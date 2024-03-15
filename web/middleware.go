package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

// RequireAuth makes sure the user is authenticated
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.Session.Exists(r.Context(), "user_id") {
			app.Session.Put(r.Context(), "error", "Login first")
			http.Redirect(w, r, "/join/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireVerified(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Session.Exists(r.Context(), "verified") && app.Session.GetBool(r.Context(), "verified") == false {
			app.Session.Put(r.Context(), "error", "You must verify your email at personal cabinet!")
			http.Redirect(w, r, "/verify-mail", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func User(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Session.Exists(r.Context(), "user_role") && app.Session.GetInt(r.Context(), "user_role") != 3 {
			app.Session.Put(r.Context(), "error", "You don't have permission to access this page!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Head(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Session.Exists(r.Context(), "user_role") && app.Session.GetInt(r.Context(), "user_role") != 2 {
			app.Session.Put(r.Context(), "error", "You don't have permission to access this page!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func Admin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Session.Exists(r.Context(), "user_role") && app.Session.GetInt(r.Context(), "user_role") != 1 {
			app.Session.Put(r.Context(), "error", "You don't have permission to access this page!")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
