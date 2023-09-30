package handlers

import (
	"book-manager-server/db"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type CreateRequestBody struct {
	Name            string    `json:"name"`
	Author          Author    `json:"author"`
	Category        string    `json:"category"`
	Volume          int       `json:"volume"`
	PublishedAt     time.Time `json:"published_at"`
	Summary         string    `json:"summary"`
	TableOfContents []string  `json:"table_of_contents"`
	Publisher       string    `json:"publisher"`
}

type Author struct {
	Firstname   string    `jason:"first_name"`
	Lastname    string    `json:"last_name"`
	Birthday    time.Time `json:"birthday"`
	Nationality string    `json:"nationality"`
}

func (bms *BookManagerServer) HandleCrud(w http.ResponseWriter, r *http.Request) {
	// Get authorization token from header
	token := r.Header.Get("Authentication")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the account by token
	account, err := bms.Auth.GetAccountByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve user from database
	_, err = bms.DB.GetUserByUsername(*account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodPost:
		// Parse the request body
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var crb CreateRequestBody
		if err = json.Unmarshal(reqData, &crb); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		author := db.Author{
			FirstName:   crb.Author.Firstname,
			LastName:    crb.Author.Lastname,
			Birthday:    crb.Author.Birthday,
			Nationality: crb.Author.Nationality,
		}
		err = bms.DB.CreteNewBook(&db.Book{
			Name:            crb.Name,
			Author:          author,
			Category:        crb.Category,
			Volume:          crb.Volume,
			PublishedAt:     crb.PublishedAt,
			Summary:         crb.Summary,
			TableOfContents: crb.TableOfContents,
			Publisher:       crb.Publisher,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("Could not write signup data to db")
			return
		}

		response := map[string]interface{}{
			"message": "Book created successfully",
		}

		resBody, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(resBody)

	case http.MethodGet:
		books, err := bms.DB.GetAllBooks()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := map[string][]db.Book{
			"books": books,
		}

		resBody, _ := json.Marshal(response)
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}
}

func (bms *BookManagerServer) HandleGetBook(w http.ResponseWriter, r *http.Request) {
	// Get authorization token from header
	token := r.Header.Get("Authentication")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve the account by token
	account, err := bms.Auth.GetAccountByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Retrieve user from database
	_, err = bms.DB.GetUserByUsername(*account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case http.MethodGet:
		id, err := strconv.Atoi(r.URL.Path[len("/api/v1/books/"):])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("Invalid id number")
			return
		}

		book, err := bms.DB.GetBookByID(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resBody, err := json.Marshal(book)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)

	case http.MethodPut:
		id, err := strconv.Atoi(r.URL.Path[len("/api/v1/books/"):])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("Invalid id number")
			return
		}

		// Parse the request body
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var updateForm map[string]interface{}
		if err = json.Unmarshal(reqData, &updateForm); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = bms.DB.UpdateBook(id, updateForm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("Could not update the fields")
			return
		}
		response := map[string]interface{}{
			"message": "Book updated successfully",
		}

		resBody, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)

	case http.MethodDelete:
		id, err := strconv.Atoi(r.URL.Path[len("/api/v1/books/"):])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print("Invalid id number")
			return
		}

		err = bms.DB.DeleteBook(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"message": "Book deleted successfully",
		}

		resBody, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resBody)

	}
}
