package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

const (
	OAuthTimeLimit   = 5 * time.Minute
	AdminAccessLimit = 24 * time.Hour
)

type GoogleOAuthData struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

// OAuthResponse is the response body of oauth login requests
type OAuthResponse struct {
	Token string `json:"token"`
}

// Render preprocess the response before it's sent to the wire
func (rsp *OAuthResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// do nothing
	return nil
}

// NewAlbumResponse is the constructor method for AlbumResponse type
func NewOAuthResponse(token string) *OAuthResponse {
	resp := OAuthResponse{Token: token}

	return &resp
}

// AuthRouter handles all routes for authentication
func (handler Handler) AuthRouter() chi.Router {
	r := chi.NewRouter()

	r.HandleFunc("/google/login", handler.OAuthGoogleLogin)
	r.HandleFunc("/google/callback", handler.oauthGoogleCallback)

	return r
}

// OAuthGoogleLogin is the login endpoint for Google OAuth API
func (handler Handler) OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauthState := generateStateOauthCookie(w)
	u := handler.googleOAuthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {

	var expiration = time.Now().Add(OAuthTimeLimit)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (handler Handler) oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		err := fmt.Errorf("invalid Google OAuth state: %s", r.FormValue("state"))
		log.Println(err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	data, err := handler.getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userData := GoogleOAuthData{}
	if err := json.Unmarshal(data, &userData); err != nil {
		log.Println(err.Error())
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// check if the user is admin
	if !(userData.VerifiedEmail && userData.Email == handler.adminID) {
		log.Printf("Unauthorized access. UserData[%s]", data)
		render.Render(w, r, ErrUnauthorized(errors.New("unauthorized access")))
		return
	}

	// respond with an admin JWT on success
	token, _ := handler.generateAdminJWT(time.Now().Add(AdminAccessLimit))
	render.Render(w, r, NewOAuthResponse(token))
}

func (handler Handler) getUserDataFromGoogle(code string) ([]byte, error) {
	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

	// Use code to get token and get user info from Google.
	token, err := handler.googleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange error: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}

	return contents, nil
}
