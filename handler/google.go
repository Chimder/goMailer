package handler

import (
	"encoding/json"
	"goMailer/auth"
	"log"
	"net/http"
	"time"
)

type Empty struct {
}

type GoogleAccount struct {
	Name              string `json:"name,omitempty"`
	ProviderId        string `json:"provider_id,omitempty"`
	ProviderAccountId string `json:"provider_account_id,omitempty"`
	Email             string `json:"email,omitempty"`
	Picture           string `json:"picture,omitempty"`
	AccessToken       string `json:"access_token,omitempty"`
	RefreshToken      string `json:"refresh_token,omitempty"`
	UserId            string `json:"user_id,omitempty"`
}
// @Success 200 {array} Empty

// @Summary RegGoogleAcc
// @Description
// @Tags Google
// @ID get-user-list-manga
// @Accept  json
// @Produce  json
// @Param  body body string true "Reg Body"
// @Router /google/reg [get]
func RegGoogleAcc(w http.ResponseWriter, r *http.Request) {
	var newUser GoogleAccount
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoded, err := auth.Encrypt(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "googleMailer_" + newUser.ProviderId,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("sing google new account"))
}

// @Summary Delete
// @Description delete google session
// @Tags Google
// @ID delete google session
// @Accept  json
// @Produce  json
// @Param  email query string true "id"
// @Success 200 {array} Empty
// @Router /google/delete/ [get]


func DeleteGoogleCookie(w http.ResponseWriter, r *http.Request) {
	ProviderId := r.URL.Query().Get("id")
	cookieName := "googleMailer_" + ProviderId

	log.Println("Emaillll", cookieName)
	_, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, "Cookie not found", http.StatusNotFound)
		return
	}

	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("deleted"))
}

// @Summary Get Google Session
// @Description Get Google Session
// @Tags Google
// @ID get-google-session
// @Accept  json
// @Produce  json
// @Success 200 {array} GoogleAccount
// @Router /google/session/ [get]


func GetGoogleSession(w http.ResponseWriter, r *http.Request) {
	allCookie := r.Cookies()

	log.Println("ALLCCC", allCookie)

}
