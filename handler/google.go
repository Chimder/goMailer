package handler

import (
	"encoding/json"
	"goMailer/auth"
	"log"
	"net/http"
	"time"
)

type GoogleAccount struct {
	Name              string `json:"name"`
	ProviderId        string `json:"providerId"`
	ProviderAccountId string `json:"providerAccountId"`
	Email             string `json:"email"`
	Picture           string `json:"picture"`
	AccessToken       string `json:"accessToken"`
	RefreshToken      string `json:"refreshToken"`
	UserId            string `json:"userId"`
}

func RegAcc(w http.ResponseWriter, r *http.Request) {
	var newUser GoogleAccount
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("NEw Acc", newUser)

	encoded, err := auth.Encrypt(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("encodeEEE", encoded)
	// Устанавливаем куки
	cookie := &http.Cookie{
		Name:     "googleMailer_" + newUser.ProviderId,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	// w.Write([]byte("sing google new account"))

}
