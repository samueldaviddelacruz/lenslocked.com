package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	llctx "github.com/samueldaviddelacruz/lenslocked.com/context"

	"github.com/gorilla/csrf"
	"github.com/samueldaviddelacruz/lenslocked.com/models"
	"golang.org/x/oauth2"
)

func NewAuths(os models.OAuthService, dbxOAuth *oauth2.Config) *Oauths {
	return &Oauths{
		os:       os,
		dbxOAuth: dbxOAuth,
	}
}

// Users Represents a Users controller
type Oauths struct {
	os       models.OAuthService
	dbxOAuth *oauth2.Config
}

func (o *Oauths) DropboxConnect(w http.ResponseWriter, r *http.Request) {
	state := csrf.Token(r)
	cookie := http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	url := o.dbxOAuth.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (o *Oauths) DropboxCallback(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	state := r.FormValue("state")
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if cookie == nil || cookie.Value != state {
		http.Error(w, "Invalid state provided", http.StatusBadRequest)
		return
	}

	cookie.Value = ""
	cookie.Expires = time.Now()
	http.SetCookie(w, cookie)
	code := r.FormValue("code")
	token, err := o.dbxOAuth.Exchange(context.TODO(), code)

	user := llctx.User(r.Context())
	existing, err := o.os.Find(user.ID, models.OauthDropbox)
	if err == models.ErrNotFound {

	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		o.os.Delete(existing.ID)
	}

	userOAuth := models.OAuth{
		UserID:  user.ID,
		Token:   *token,
		Service: models.OauthDropbox,
	}
	err = o.os.Create(&userOAuth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%+v", token)
}

func (o *Oauths) DropboxTest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := r.FormValue("path")

	user := llctx.User(r.Context())
	userOath, err := o.os.Find(user.ID, models.OauthDropbox)
	if err != nil {
		panic(err)
	}
	token := userOath.Token
	dropBoxQuery := struct {
		Path string `json:"path"`
	}{
		Path: path,
	}
	dropBoxQueryBytes, err := json.Marshal(dropBoxQuery)
	if err != nil {
		panic(err)
	}

	client := o.dbxOAuth.Client(context.TODO(), &token)
	req, err := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder",
		bytes.NewReader(dropBoxQueryBytes),
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)

}
