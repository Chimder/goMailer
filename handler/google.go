package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goMailer/auth"
	"goMailer/config"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Empty struct {
}

type Response struct {
	MessagesData  []MessageData `json:"messagesData"`
	NextPageToken string        `json:"nextPageToken"`
}

type GoogleAccount struct {
	Name              string `json:"name,omitempty"`
	ProviderId        string `json:"providerId,omitempty"`
	ProviderAccountId string `json:"providerAccountId,omitempty"`
	Email             string `json:"email,omitempty"`
	Picture           string `json:"picture,omitempty"`
	AccessToken       string `json:"accessToken,omitempty"`
	RefreshToken      string `json:"refreshToken,omitempty"`
	UserId            string `json:"userId,omitempty"`
}

type GoogleHandler struct {
}

type MessageData struct {
	MessageId       string
	Subject         string
	From            string
	To              string
	Date            string
	Snippet         string
	IsUnread        bool
	IsBodyWithParts bool
	BodyData        string
}


// @Summary RegGoogleAcc
// @Description
// @Tags Google
// @ID get-user-list-manga
// @Accept  json
// @Produce  json
// @Param  body body string true "Reg Body"
// @Success 200 {array} Empty
// @Router /google/reg [post]
func (h *GoogleHandler) RegGoogleAcc(w http.ResponseWriter, r *http.Request) {
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
		Name:     "googleMailer_" + newUser.ProviderAccountId,
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
// @Param  id query string true "id"
// @Success 200 {array} Empty
// @Router /google/delete [get]
func (h *GoogleHandler) DeleteGoogleCookie(w http.ResponseWriter, r *http.Request) {
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
// @Router /google/session [get]
func (h *GoogleHandler) GetGoogleSession(w http.ResponseWriter, r *http.Request) {
	allCookies := r.Cookies()
	var accounts []GoogleAccount

	for _, cookie := range allCookies {
		if strings.HasPrefix(cookie.Name, "googleMailer_") {
			var account GoogleAccount
			err := auth.Decrypt(cookie.Value, &account)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accounts = append(accounts, account)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(accounts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Messages
// @Description  Get Messages and content
// @Tags Google
// @ID get-google-mess
// @Accept  json
// @Produce  json
// @Param  id query string true "id"
// @Param  pageToken query string false "pageToken"
// @Success 200 {array} Empty
// @Router /google/messages [get]
func (h *GoogleHandler) MessagesAndContent(w http.ResponseWriter, r *http.Request) {
	pageToken := r.URL.Query().Get("pageToken")
	ProviderId := r.URL.Query().Get("id")
	cookieName := "googleMailer_" + ProviderId

	var account GoogleAccount
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = auth.Decrypt(cookie.Value, &account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     config.LoadEnv().GOOGLE_CLIENT_ID,
		ClientSecret: config.LoadEnv().GOOGLE_CLIENT_SECRET,
		Endpoint:     google.Endpoint,
	}
	token := &oauth2.Token{
		AccessToken:  account.AccessToken,
		RefreshToken: account.RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	client := config.Client(ctx, token)
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	////////////////////////////////////////////////////
	user := "me"
	call := srv.Users.Messages.List(user).MaxResults(80)
	if pageToken != "" {
		call.PageToken(pageToken)
	}
	rr, err := call.Do()
	if err != nil {
		http.Error(w, "Unable to retrieve messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []*gmail.Message
	task := func(m *gmail.Message) error {
		msg, err := srv.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			return fmt.Errorf("error fetching message %s: %v", m.Id, err)
		}
		messages = append(messages, msg)
		return nil
	}
	////////////////////////////////////////////////////////////////
	var group sync.WaitGroup
	for _, m := range rr.Messages {
		group.Add(1)
		go func(m *gmail.Message) {
			defer group.Done()
			err := task(m)
			if err != nil {
				log.Printf("Error retrieving message %s: %v", m.Id, err)
			}
		}(m)
	}
	group.Wait() //all tasks to finish
	/////////////////////////////////////////////////////////////////////////
	var messagesData []MessageData
	for _, msg := range messages {
		headers := msg.Payload.Headers
		subject := getHeaderValue(headers, "Subject")
		from := extractName(getHeaderValue(headers, "From"))
		to := getHeaderValue(headers, "To")
		date := getHeaderValue(headers, "Date")

		isUnread := false
		for _, label := range msg.LabelIds {
			if label == "UNREAD" {
				isUnread = true
				break
			}
		}

		var bodyData string
		isBodyWithParts := false
		if len(msg.Payload.Parts) > 0 {
			bodyData = decodeBase64(msg.Payload.Parts[0].Body.Data)
		} else {
			isBodyWithParts = true
			bodyData = decodeBase64(msg.Payload.Body.Data)
		}

		messageData := MessageData{
			MessageId:       msg.Id,
			Subject:         subject,
			From:            from,
			To:              to,
			Date:            date,
			Snippet:         msg.Snippet,
			IsUnread:        isUnread,
			IsBodyWithParts: isBodyWithParts,
			BodyData:        bodyData,
		}
		messagesData = append(messagesData, messageData)
	}

	log.Println("LEEEG", len(messagesData))
	response := Response{
		MessagesData:  messagesData,
		NextPageToken: rr.NextPageToken,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getHeaderValue(headers []*gmail.MessagePartHeader, name string) string {
	for _, header := range headers {
		if header.Name == name {
			return header.Value
		}
	}
	return ""
}

func extractName(from string) string {
	re := regexp.MustCompile("(.*)<.*>")
	match := re.FindStringSubmatch(from)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return from
}

func decodeBase64(s string) string {
	data, _ := base64.URLEncoding.DecodeString(s)
	return string(data)
}
