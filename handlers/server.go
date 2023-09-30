package handlers

import (
	"book-manager-server/authenticate"
	"book-manager-server/db"
)

type BookManagerServer struct {
	DB   *db.GormDB
	Auth *authenticate.Auth
}
