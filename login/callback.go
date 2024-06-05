package login

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"log"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	LogoutTimeOut = 15 // Minutes
	HttpPort      = 3000
	CallbackURL   = "/auth/google/callback"
)

type CallbackResponse struct {
	User  *User
	Error error
}

type User struct {
	Email     string
	Name      string
	UserID    string
	ExpiresAt time.Time
}

func InitProviders() error {
	if GoogleClientID == "" {
		return fmt.Errorf("GoogleClientID is not set in file client/secret.go")
	}

	url := fmt.Sprintf("http://localhost:%d%s", HttpPort, CallbackURL)
	goth.UseProviders(
		google.New(GoogleClientID, GoogleClientSecret, url, "email", "profile"),
	)

	return nil
}

func RunHTTPServer(stop chan CallbackResponse) {
	mux := http.NewServeMux()
	mux.HandleFunc(CallbackURL, func(res http.ResponseWriter, req *http.Request) {
		user, err := AuthCallback(res, req)
		stop <- CallbackResponse{
			User:  user,
			Error: err,
		}
	})

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", HttpPort), mux)
		if err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()
}

func AuthCallback(res http.ResponseWriter, req *http.Request) (*User, error) {
	user, err := CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return nil, err
	}

	tmpl, err := template.New("go").Parse(getTemplate())
	if err != nil {
		fmt.Fprintln(res, err)
		return nil, err
	}

	err = WriteUser(user)
	if err != nil {
		fmt.Fprintln(res, err)
		return nil, err
	}

	err = tmpl.Execute(res, user)
	if err != nil {
		fmt.Fprintln(res, err)
		return nil, err
	}

	return user, nil
}

func CompleteUserAuth(res http.ResponseWriter, req *http.Request) (*User, error) {
	provider, err := goth.GetProvider("google")
	if err != nil {
		return nil, err
	}

	sess := &google.Session{}
	params := req.URL.Query()
	if params.Encode() == "" && req.Method == "POST" {
		req.ParseForm()
		params = req.Form
	}

	_, err = sess.Authorize(provider, params)
	if err != nil {
		return nil, err
	}

	gu, err := provider.FetchUser(sess)
	user := &User{
		Email:     gu.Email,
		Name:      gu.Name,
		UserID:    gu.UserID,
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(LogoutTimeOut)),
	}
	return user, err
}
