package handlers

import (
	"book-manager-server/authenticate"
	"book-manager-server/db"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type SignUpRequestBody struct {
	Username    string `json:"user_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Firstname   string `json:"first_name"`
	Lastname    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Gender      string `json:"gender"`
}

func (bms *BookManagerServer) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	// Check the API method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Print("Invalid method for signup")
		return
	}

	// Parse the request body
	reqData, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var srb SignUpRequestBody
	if err = json.Unmarshal(reqData, &srb); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = bms.DB.CreateNewUser(&db.User{
		Username:    srb.Username,
		Email:       srb.Email,
		Password:    srb.Password,
		Firstname:   srb.Firstname,
		Lastname:    srb.Lastname,
		PhoneNumber: srb.PhoneNumber,
		Gender:      srb.Gender,
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print("Could not write signup data to db")
		return
	}

	response := map[string]interface{}{
		"message": "User created successfully",
	}

	resBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(resBody)
}

type LoginRequestBody struct {
	Username string `json:"user_name"`
	Password string `json:"password"`
}

func (bms *BookManagerServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// Check the API metho
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Read request body
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the request data
	var lrb LoginRequestBody
	err = json.Unmarshal(reqBody, &lrb)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := bms.Auth.Login(authenticate.Credentials{
		Username: lrb.Username,
		Password: lrb.Password,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"access_token": token.TokenString,
	}

	resBody, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resBody)

}
