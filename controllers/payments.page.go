package controllers

import (
	"encoding/json"
	"net/http"
)

type walletData struct {
	Address string `json:"walletAddress"`
}

func (m *Repository) PaymentsConnectMetamask(w http.ResponseWriter, r *http.Request) {
	var wallet walletData
	if err := json.NewDecoder(r.Body).Decode(&wallet); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	responseData := map[string]interface{}{
		"message": "Wallet connected successfully",
		"address": wallet.Address,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)

}
