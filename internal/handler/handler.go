package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/Isotton1/web-authenticatior/internal/common"
	"github.com/Isotton1/web-authenticatior/models"
	"github.com/Isotton1/web-authenticatior/pkg/crypto"
	"github.com/Isotton1/web-authenticatior/pkg/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := loadTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := loadTemplate(w, "signup.html", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		comfirmation := r.FormValue("confirm password")
		if password == comfirmation {
			salt, err := crypto.GenerateRandomString(128)
			if err != nil {
				log.Fatal(err)
			}
			pepper, err := crypto.GenerateRandomString(128)
			if err != nil {
				log.Fatal(err)
			}
			passwordHash := crypto.NewHash(salt + password + pepper)
			encryptedSalt, err := crypto.Encrypt(salt, password)
			if err != nil {
				log.Fatal(err)
			}
			encryptedPepper, err := crypto.Encrypt(pepper, password)
			if err != nil {
				log.Fatal(err)
			}
			user := models.User{
				Username:     username,
				PasswordHash: passwordHash,
				Salt:         encryptedSalt,
				Pepper:       encryptedPepper,
			}
			err = database.InsertUser(&user)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "<div>Login successful!</div>")
			w.WriteHeader(http.StatusAccepted)
		} else {
			fmt.Fprintf(w, "<div>Invalid credentials</div>")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		err := loadTemplate(w, "login.html", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := database.GetUser(username)
		if err != nil {
			if err == common.ErrNoUserFound {
				fmt.Fprintf(w, "<div>Do not exist a account with this username</div>")
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
			}
		}
		salt, err := crypto.Decrypt(user.Salt, password)
		pepper, err := crypto.Decrypt(user.Pepper, password)
		passwordHash := crypto.NewHash(salt + password + pepper)
		if passwordHash == user.PasswordHash {
			fmt.Fprintf(w, "<div>Login successful!</div>")
			w.WriteHeader(http.StatusAccepted)
		} else {
			fmt.Fprintf(w, "<div>Invalid credentials</div>")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}

}

func loadTemplate(w http.ResponseWriter, fileName string, data interface{}) error {
	tmpl, err := template.ParseFiles(fileName)
	if err != nil {
		return errors.New("Error loading the template: " + err.Error())
	}
	err = tmpl.ExecuteTemplate(w, fileName, data)
	if err != nil {
		return err
	}
	return nil

}
