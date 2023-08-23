package user

import (
	"database/sql"
	"html"
	"log"
	"net/http"
	"strings"

	httputils "github.com/ecuyle/gomine/internal/http"
	"github.com/ecuyle/gomine/internal/passwords"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserOptions struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       string
	Username string
	Hash     string
}

func insertUser(user *User) error {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return err
	}

	defer db.Close()
	transaction, err := db.Begin()

	if err != nil {
		return err
	}

	statement, err := transaction.Prepare("insert into users(id, username, hash) values(?, ?, ?)")

	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(user.ID, user.Username, user.Hash)

	if err != nil {
		return err
	}

	err = transaction.Commit()

	if err != nil {
		return err
	}

	return nil
}

func makeUser(username, password string) (*User, error) {
	// hash, err := password.GenerateHashFromPassword(password)
	hash, err := passwords.GenerateHashFromPassword(password)

	if err != nil {
		return nil, err
	}

	id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	escapedUsername := html.EscapeString(strings.TrimSpace(username))

	return &User{ID: id.String(), Username: escapedUsername, Hash: hash}, nil
}

func PostUser(context *gin.Context) {
	var options UserOptions

	if err := context.BindJSON(&options); err != nil {
		log.Println(err)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := makeUser(options.Username, options.Password)

	if err != nil {
		log.Println("Error making user")
		log.Println(err)
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = insertUser(user)

	if err != nil {
		log.Println("Error inserting user")
		log.Println(err)
		context.String(http.StatusInternalServerError, err.Error())
		return
	}

	httputils.RespondWithStatusCreated(context, map[string]string{
		"id":       user.ID,
		"username": user.Username,
	})
}
