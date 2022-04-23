package controllers

import (
	"encoding/json"
	"github.com/jgsheppa/search_engine/context"
	"github.com/jgsheppa/search_engine/models"
	"github.com/jgsheppa/search_engine/rand"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type User struct {
	us models.UserService
}

func Auth(us models.UserService) *User {
	return &User{
		us: us,
	}
}

type RegistrationForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the registration form
//
// POST /register
func (u *User) Create(w http.ResponseWriter, r *http.Request) {
	var form RegistrationForm

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	user.Role = "user"
	if err := u.us.Create(&user); err != nil {
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusNotFound)
		return
	}

}

type LoginForm struct {
	Email    string `schema:"email" json:"email"`
	Password string `schema:"password" json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (u *User) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Path:     "/",
		Value:    user.Remember,
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	return nil
}

// Login godoc
// @Summary Login to the Redisearch API
// @Tags Auth
// @Param credentials body LoginForm true "Credentials"
// @Success 201 {string} string "Ok"
// @Failure 422 {object} models.ApiError
// @Router /api/auth/login [post]
func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		json.NewEncoder(w).Encode(models.LargePayloadError)
		w.WriteHeader(models.LargePayloadError.HttpStatus)
		return
	}
	if err := r.Body.Close(); err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}
	if err := json.Unmarshal(body, &form); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(models.ValidationError)
		w.WriteHeader(models.ValidationError.HttpStatus)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		panic(err)
	}

	err = u.signIn(w, user)
	if err != nil {
		panic(err)
	}

	token := user.Remember
	response := LoginResponse{
		Message: "Login Successful",
		Token:   token,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout godoc
// @Summary Logout of the Redisearch API
// @Tags Auth
// @Success 201 {string} string "Ok"
// @Failure 401 {object} models.ApiError
// @Router /api/auth/logout [post]
func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	if user == nil {
		json.NewEncoder(w).Encode(models.AuthError)
		return
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)

	token, _ := rand.RememberToken()
	user.Remember = token
	err := u.us.Update(user)
	if err != nil {
		json.NewEncoder(w).Encode(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
	}

	user = context.User(r.Context())

	response := Response{
		Message: "Logout Successful",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
