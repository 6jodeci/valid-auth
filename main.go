package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/registrations", registrationsHandler)
	http.HandleFunc("/authentications", authenticationsHandler)
	http.HandleFunc("/test", testResponceHandler)
	http.ListenAndServe(":8081", nil)
}

func registrationsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		fmt.Fprintf(w, "Please enter a valid username and password.\r\n")
	} else {
		response, err := registerUser(r.FormValue("username"), r.FormValue("password"))

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, response)
		}
	}
}

func authenticationsHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()

	if ok {
		tokenDetails, err := generateToken(username, password)

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			enc := json.NewEncoder(w)
			enc.SetIndent("", "  ")
			enc.Encode(tokenDetails)
		}
	} else {
		fmt.Fprintf(w, "You require a username/password to get token.\r\n")
	}
}

func testResponceHandler(w http.ResponseWriter, r *http.Request) {
	authToken := strings.Split(r.Header.Get("Authorization"), "Bearer")[1]
	userDetails, err := validateToken(authToken)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		username := fmt.Sprint(userDetails["username"])
		fmt.Fprintf(w, "Welcome, "+username+"\r\n")
	}
}
