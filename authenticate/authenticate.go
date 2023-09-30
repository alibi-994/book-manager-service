package authenticate

import (
	"book-manager-server/db"
	"crypto/rand"
	"errors"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	db        *db.GormDB
	secretKey []byte
}

type Credentials struct {
	Username string
	Password string
}

type Claims struct {
	jwt.MapClaims
	Username string `json:"user_name"`
}

type Token struct {
	TokenString string
}

func GenerateSecretKey() ([]byte, error) {
	secretKey := make([]byte, 31)
	_, err := rand.Read(secretKey)
	if err != nil {
		return nil, err
	}
	return secretKey, nil
}

func InitAuth(gdb *db.GormDB) (*Auth, error) {
	secretKey, err := GenerateSecretKey()
	if err != nil {
		return nil, err
	}

	if gdb == nil {
		return nil, errors.New("DB can not be nil")
	}

	return &Auth{
		db:        gdb,
		secretKey: secretKey,
	}, nil
}

func (a *Auth) Login(cred Credentials) (Token, error) {
	user, err := a.db.GetUserByUsername(cred.Username)
	if err != nil {
		log.Print("User does not exist")
		return Token{}, err
	}

	// Check if the password match
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(cred.Password))
	if err != nil {
		log.Print("Password does not match")
		return Token{}, err
	}

	// Generate jwtToken
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Username:  cred.Username,
		MapClaims: jwt.MapClaims{},
	})

	tokenString, err := tokenJWT.SignedString(a.secretKey)
	if err != nil {
		return Token{}, nil
	}

	return Token{TokenString: tokenString}, nil

}

func (a *Auth) CheckToken(tokenstr string) (*Claims, error) {
	c := &Claims{}
	tkn, _ := jwt.ParseWithClaims(tokenstr, c, func(t *jwt.Token) (interface{}, error) {
		return a.secretKey, nil
	})

	if !tkn.Valid {
		return nil, errors.New("Unauthorized")
	}

	return c, nil
}

func (a *Auth) GetAccountByToken(token string) (*string, error) {
	if token == "" {
		log.Print("Access denied")
		return nil, errors.New("Token not Valid")
	}

	claim, err := a.CheckToken(token)
	if err != nil {
		return nil, err
	}

	return &claim.Username, nil
}
