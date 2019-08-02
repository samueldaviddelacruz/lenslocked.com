package controllers

import (
	"fmt"
	"net/http"

	"github.com/samueldaviddelacruz/lenslocked.com/models"

	"github.com/samueldaviddelacruz/lenslocked.com/views"
)

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

// Users Represents a Users controller
type Users struct {
	NewView *views.View
	us      *models.UserService
}

// New is used to render the form where a user can
// create a new user account
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)

}
