package authentication

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ecuyle/gomine/internal/passwords"
	"github.com/ecuyle/gomine/internal/token"
	"github.com/ecuyle/gomine/internal/user"
	"github.com/gin-gonic/gin"
)

type AuthenticationOptions struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func retrieveAccessTokenIfCredentialsValid(username, password string) (string, error) {
	db, err := sql.Open("sqlite3", "./gomine.db")

	if err != nil {
		return "", err
	}

	defer db.Close()
	transaction, err := db.Begin()

	if err != nil {
		return "", err
	}

	fmt.Println("3")
	statement, err := transaction.Prepare("select * from users where username=?")

	if err != nil {
		return "", err
	}

	defer statement.Close()

	user := user.User{}
	err = statement.QueryRow(username).Scan(&user.ID, &user.Username, &user.Hash)

	if err != nil {
		return "", err
	}

	err = passwords.ComparePasswordWithHash(password, user.Hash)

	if err != nil {
		return "", err
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		return "", err
	}

	return token, nil
}

func AuthenticateUser(c *gin.Context) {
	var options AuthenticationOptions

	if err := c.ShouldBindJSON(&options); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := retrieveAccessTokenIfCredentialsValid(options.Username, options.Password)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": token})

}
