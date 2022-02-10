package http

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	pa "github.com/Lambels/patrickarvatu.com"
	"github.com/go-chi/chi/v5"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// registerAuthRoutes registers the auth routes under r.
func (s *Server) registerAuthRoutes(r chi.Router) {
	r.Delete("/logout", s.handleLogout)
	r.Get("/github", s.handleGithubOAuth)
	r.Post("/github/callback", s.handleGithubCallback)

	r.Route("/user", func(r chi.Router) {
		r.Use(s.jsonResponseTypeMiddleware)
		r.Use(s.requireAuthMiddleware)
		r.Get("/me", s.handleMe)
		r.Get("/check-auth", s.handleCheckAuth)
	})
}

// handleLogout handels DELETE '/oauth/logout'.
// sets an empty session with no user id indicating an unauth request.
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	// clear session.
	if err := s.setSession(w, pa.Session{}); err != nil {
		SendError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleGithubOAuth handles GET '/oauth/github'.
// prepare state and redirect to github oauth endpoint for transaction.
func (s *Server) handleGithubOAuth(w http.ResponseWriter, r *http.Request) {
	ses, err := s.getSession(r)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// generate random bytes for state.
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		SendError(w, r, err)
		return
	}

	// set state under session.
	ses.State = base64.StdEncoding.EncodeToString(buf)

	// set state to check on callback.
	if err := s.setSession(w, ses); err != nil {
		SendError(w, r, err)
		return
	}

	// redirect to provider.
	http.Redirect(w, r, s.NewOAuthConfig(pa.AuthSourceGitHub).AuthCodeURL(ses.State), http.StatusFound)
}

// handleGithubCallback handles POST '/oauth/github/callback'.
// validates state for possible csrf attack.
// exchanges resource owners grant token for access token.
// fetches at least user id from github.
// creates auth object and user.
func (s *Server) handleGithubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state") // state sent back from github.
	code := r.FormValue("code")   // temp grant code.

	// read session from request.
	ses, err := s.getSession(r)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// validate that the state comming from the request matches the state from the request to auth.
	if ses.State != state {
		SendError(w, r, fmt.Errorf("request state and response state mismatch"))
		return
	}

	// exchange grant code for access token and refresh token, but github doesent provide
	// refresh tokens.
	token, err := s.NewOAuthConfig(pa.AuthSourceGitHub).Exchange(r.Context(), code)
	if err != nil {
		SendError(w, r, err)
		return
	}

	// github api client
	client := github.NewClient(oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token.AccessToken,
	})))

	// pass empty string to get the auth user under the token.
	user, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		SendError(w, r, err)
		return
	} else if user.ID == nil {
		SendError(w, r, fmt.Errorf("user id not returned by github, abording auth process"))
		return
	}

	var name string
	if user.Name != nil {
		name = *user.Name
	} else if user.Login != nil {
		name = *user.Login
	}

	// get optional email.
	var email string
	var isAdmin bool
	if user.Email != nil {
		email = *user.Email
		if email == s.conf.User.AdminUserEmail {
			isAdmin = true
		}
	}

	auth := &pa.Auth{
		User: &pa.User{
			Name:    name,
			Email:   email,
			IsAdmin: isAdmin,
		},
		Source:      pa.AuthSourceGitHub,
		SourceID:    strconv.FormatInt(*user.ID, 10),
		AccessToken: token.AccessToken,
	}
	if !token.Expiry.IsZero() {
		auth.Expiry = &token.Expiry
	}

	// create the auth.
	if err := s.AuthService.CreateAuth(r.Context(), auth); err != nil {
		SendError(w, r, err)
		return
	}

	// clear state.
	ses.State = ""
	// set userID.
	ses.UserID = auth.UserID // populated on creation

	if err := s.setSession(w, ses); err != nil {
		SendError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// handleMe handels GET '/oauth/user/me'.
// returns an userProfileResponse.
func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	v := pa.UserFromContext(r.Context())
	sendData := struct {
		user   *pa.User `json:"user"`
		pfpURL string   `json:"pfpUrl"`
	}{
		user:   v,
		pfpURL: v.AvatarURL(10),
	}
	SendJSON(w, sendData)
}

// handleCheckAuth handles GET '/oauth/user/check-auth'.
// if request reaches here user is auth.
func (s *Server) handleCheckAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
