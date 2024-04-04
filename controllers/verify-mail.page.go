package controllers

import (
	"bytes"
	"github.com/seemsod1/goober/helpers"
	"github.com/seemsod1/goober/helpers/forms"
	"github.com/seemsod1/goober/helpers/render"
	models "github.com/seemsod1/goober/models/app_models"
	"github.com/seemsod1/goober/models/entities"
	"html/template"
	"net/http"
	"time"
)

func (m *Repository) VerifyMail(w http.ResponseWriter, r *http.Request) {
	if !m.App.Session.Exists(r.Context(), "verified") {
		m.App.Session.Put(r.Context(), "error", "You are already verified")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	render.RenderTemplate(w, r, "verify-mail.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) VerifyMailPost(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Can't parse form", http.StatusBadRequest)
		return
	}
	code := r.Form.Get("code")

	form := forms.New(r.PostForm)
	digit1 := r.Form.Get("digit-1")
	digit2 := r.Form.Get("digit-2")
	digit3 := r.Form.Get("digit-3")
	digit4 := r.Form.Get("digit-4")
	digit5 := r.Form.Get("digit-5")
	digit6 := r.Form.Get("digit-6")
	form.Required("digit-1", "digit-2", "digit-3", "digit-4", "digit-5", "digit-6")
	form.MinLength("digit-1", 1)
	form.Maxlength("digit-1", 1)
	form.MinLength("digit-2", 1)
	form.Maxlength("digit-2", 1)
	form.MinLength("digit-3", 1)
	form.Maxlength("digit-3", 1)
	form.MinLength("digit-4", 1)
	form.Maxlength("digit-4", 1)
	form.MinLength("digit-5", 1)
	form.Maxlength("digit-5", 1)
	form.MinLength("digit-6", 1)
	form.Maxlength("digit-6", 1)

	if !form.Valid() {
		render.RenderTemplate(w, r, "verify-mail.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	code = digit1 + digit2 + digit3 + digit4 + digit5 + digit6

	userId := m.App.Session.GetInt(r.Context(), "user_id")

	var user entities.User
	res := m.App.DB.Where("id = ?", userId).First(&user)
	if res.Error != nil {
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	var verifdata entities.VerificationData
	res = m.App.DB.Where("email = ?", user.Email).Where("code = ?", code).Where("expires_at > ?", time.Now()).First(&verifdata)
	if res.Error != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid verification code")
		http.Redirect(w, r, "/verify-mail", http.StatusSeeOther)
		return
	}

	m.App.DB.Delete(&verifdata)
	user.IsVerified = true
	m.App.DB.Save(&user)

	m.App.Session.Pop(r.Context(), "verified")
	m.App.Session.Put(r.Context(), "flash", "Your email is verified")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) GetVerificationCode(w http.ResponseWriter, r *http.Request) {
	userId := m.App.Session.GetInt(r.Context(), "user_id")

	var user entities.User
	res := m.App.DB.Where("id = ?", userId).First(&user)
	if res.Error != nil {
		w.WriteHeader(http.StatusSeeOther)
		return
	}
	needVerif := m.App.Session.Exists(r.Context(), "verified")
	if !needVerif {
		m.App.Session.Put(r.Context(), "error", "You are already verified")
	}

	var verifdata entities.VerificationData
	res = m.App.DB.Where("email = ?", user.Email).Where("expires_at > ?", time.Now()).First(&verifdata)
	if res.Error == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	var verData []entities.VerificationData
	m.App.DB.Where("email = ?", user.Email).Find(&verData)
	for _, v := range verData {
		m.App.DB.Delete(&v)
	}

	code := helpers.GenerateRandomString(6)
	verifdata = entities.VerificationData{
		Email:     user.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	m.App.DB.Create(&verifdata)

	emailTemplate := `
        <p>Dear {{.Name}},</p>
        <p>Thank you for subscribing to our mailing list. To verify your email address and complete the subscription process, please use the following verification code:</p>
        <p><strong>Verification Code:</strong> {{.VerificationCode}}</p>
        <p>Please enter this code on the website to confirm your subscription. If you didn't request this email or have any questions, please contact us immediately.</p>
        <p>Thank you,<br>Goober Team</p>
    `

	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		// handle error
		panic(err)
	}

	var emailContentBuf bytes.Buffer
	err = tmpl.Execute(&emailContentBuf, struct {
		Name             string
		VerificationCode string
	}{
		Name:             user.Name,
		VerificationCode: code,
	})
	if err != nil {
		// handle error
		panic(err)
	}

	msg := models.MailData{
		To:      user.Email,
		From:    "",
		Subject: "Verification code",
		Content: template.HTML(emailContentBuf.String()),
	}

	m.App.MailChan <- msg

	w.WriteHeader(http.StatusOK)

}
