package webview

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/pkg/domain"
	"github.com/gorilla/mux"
)

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	tmplSignup.Execute(w, nil)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	tmplLogin.Execute(w, nil)
}

// ShowProfile godoc
//
//	@Summary		Show an account
//	@Description	get profile by id
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200	{object}	domain.User
//	@Failure		400	{object}	domain.User
//	@Failure		404	{object}	domain.User
//	@Failure		500	{object}	domain.User
//	@Router			/profile/{id} [get]
func (s *Server) profile(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]
	requestedId, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, ok := r.Context().Value(keyUserId).(int32)
	if !ok {
		http.Error(w, "no id in request context", http.StatusInternalServerError)
		return
	}
	if id != int32(requestedId) {
		http.Error(w, "you are not authorized to visit this profile", http.StatusUnauthorized)
		return
	}

	user, err := s.service.GetUser(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tmplProfile.Execute(w, user)
}

func (s *Server) editProfile(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]
	requestedId, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, ok := r.Context().Value(keyUserId).(int32)
	if !ok {
		http.Error(w, "no id in request context", http.StatusInternalServerError)
		return
	}
	if id != int32(requestedId) {
		http.Error(w, "you are not authorized to visit this profile", http.StatusUnauthorized)
		return
	}

	r.ParseForm()
	r.PostFor
	var user domain.User
	user.ID = id
	user.FirstName = r.PostFormValue("firstName")
	user.LastName = r.PostFormValue("lastName")
	user.Email = r.PostFormValue("email")
	user.Status = r.PostFormValue("status")
	user.Description = r.PostFormValue("desription")
	err = s.service.UpdateUserDetails(context.Background(), &user)
	if err != nil {
		http.Error(w, fmt.Errorf("error updating user: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	updatedUser, err := s.service.GetUser(context.Background(), id)
	if err != nil {
		http.Error(w, fmt.Errorf("error getting updated user: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	tmplProfile.Execute(w, updatedUser)
}

func (s *Server) myProfile(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(keyUserId).(int32)
	if !ok {
		http.Error(w, "no id in request context", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/profile/%d", id), http.StatusFound)
}

func (s *Server) validateSignUp(w http.ResponseWriter, r *http.Request) {
	type Data struct {
		Username     string
		ErrorMessage string
	}

	var data Data

	r.ParseForm()
	var username string
	var password string
	var confirmPassword string
	username = r.PostFormValue("username")
	password = r.PostFormValue("password")
	confirmPassword = r.PostFormValue("confirm_password")
	fmt.Println(username, password, confirmPassword)
	data.Username = username
	if password != confirmPassword {
		data.ErrorMessage = "passwords don't match"
	} else if len(username) == 0 {
		data.ErrorMessage = "username is not specified"
	} else {
		id, err := s.service.CreateUser(context.Background(),
			&domain.User{Username: username, Password: password})
		if err != nil {
			log.Println("error creating user: ", err)
			return
		}

		token, err := s.service.GenerateToken(context.Background(), username, password)
		if err != nil {
			log.Print(err)
			http.Error(w, "error validating credentials", http.StatusUnauthorized)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "mvcAuthToken",
			Value:   token,
			Expires: time.Now().Add(domain.TokenTTL),
		})

		fmt.Println("new user id: ", id)
		http.Redirect(w, r, "/profile", http.StatusFound)
		return
	}
	tmplSignup.Execute(w, data)
}

func (s *Server) validateLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var username string
	var password string
	username = r.PostFormValue("username")
	password = r.PostFormValue("password")
	if username == "" {
		http.Error(w, "username is not specified", http.StatusBadRequest)
		return
	}
	if password == "" {
		http.Error(w, "password is not specified", http.StatusBadRequest)
		return
	}
	token, err := s.service.GenerateToken(context.Background(), username, password)
	if err != nil {
		log.Print(err)
		http.Error(w, "error validating credentials", http.StatusUnauthorized)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "mvcAuthToken",
		Value:   token,
		Expires: time.Now().Add(domain.TokenTTL),
	})
	http.Redirect(w, r, "/profile", http.StatusFound)
}
