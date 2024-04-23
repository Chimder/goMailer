package handler

import (
	"encoding/json"
	"goMailer/auth"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/felixstrobel/mailtm"
)

type TempAccount struct {
	ID         string    `json:"id"`
	Address    string    `json:"address"`
	Quota      int       `json:"quota"`
	Used       int       `json:"used"`
	IsDisabled bool      `json:"isDisabled"`
	IsDeleted  bool      `json:"isDeleted"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`

	Password string
	Token    string
}

type TempHandler struct {
}

type Addressee struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type DetailedMessage struct {
	ID             string      `json:"id"`
	AccountID      string      `json:"accountId"`
	MessageID      string      `json:"msgid"`
	From           Addressee   `json:"from"`
	To             []Addressee `json:"to"`
	CC             []Addressee `json:"cc"`
	BCC            []Addressee `json:"bcc"`
	Subject        string      `json:"subject"`
	Seen           bool        `json:"seen"`
	Flagged        bool        `json:"flagged"`
	IsDeleted      bool        `json:"isDeleted"`
	Verifications  []string    `json:"verifications"`
	Retention      bool        `json:"retention"`
	RetentionDate  time.Time   `json:"retentionDate"`
	Text           string      `json:"text"`
	Html           []string    `json:"html"`
	HasAttachments bool        `json:"hasAttachments"`
	Attachments    []string    `json:"attachments"`
	Size           int         `json:"size"`
	DownloadUrl    string      `json:"downloadUrl"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

type Message struct {
	ID             string      `json:"id"`
	AccountID      string      `json:"accountId"`
	MessageID      string      `json:"msgid"`
	From           Addressee   `json:"from"`
	To             []Addressee `json:"to"`
	Subject        string      `json:"subject"`
	Intro          string      `json:"intro"`
	Seen           bool        `json:"seen"`
	IsDeleted      bool        `json:"isDeleted"`
	HasAttachments bool        `json:"hasAttachments"`
	Size           int         `json:"size"`
	DownloadUrl    string      `json:"downloadUrl"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

// @Summary RegTempAcc
// @Description get randome tempAcc
// @Tags Temp
// @ID get-temp-mail
// @Accept  json
// @Produce  json
// @Success 200 {array} TempAccount
// @Router /temp/reg [get]
func (t *TempHandler) RegTempEmail(w http.ResponseWriter, r *http.Request) {

	client, err := mailtm.New()

	account, err := client.NewAccount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("ACCTEMP", account)

	encoded, err := auth.Encrypt(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "googleMailer_" + account.ID,
		Value:    encoded,
		Path:     "/",
		Expires:  time.Now().Add(1 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
	w.Write([]byte("create temp new account"))
}

// @Summary Temp
// @Description get all Temp messages
// @Tags Temp
// @ID get-temp-messages
// @Accept  json
// @Produce  json
// @Param  id query string true "id"
// @Param  pageToken query int false "pageToken"
// @Success 200 {array} Message
// @Router /temp/messages [get]
func (t *TempHandler) GetTempMessages(w http.ResponseWriter, r *http.Request) {
	pageToken := r.URL.Query().Get("pageToken")
	pageTokenInt, err := strconv.Atoi(pageToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ProviderId := r.URL.Query().Get("id")
	cookieName := "googleMailer_" + ProviderId

	var account *mailtm.Account
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = auth.Decrypt(cookie.Value, &account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	client, err := mailtm.New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	messages, err := client.GetMessages(account, pageTokenInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Temp
// @Description get one Temp messages
// @Tags Temp
// @ID get-one-temp-message
// @Accept  json
// @Produce  json
// @Param  messageId query string true "messageId"
// @Param  id query string true "id"
// @Success 200 {object} DetailedMessage
// @Router /temp/message [get]
func (t *TempHandler) GetTempMessage(w http.ResponseWriter, r *http.Request) {
	messageId := r.URL.Query().Get("messageId")
	accountId := r.URL.Query().Get("id")
	cookieName := "googleMailer_" + accountId

	var account *mailtm.Account
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = auth.Decrypt(cookie.Value, &account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("OO", account)

	client, err := mailtm.New()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	message, err := client.GetMessageByID(account, messageId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	log.Println("OONE", message)
}
