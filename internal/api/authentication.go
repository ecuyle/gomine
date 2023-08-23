package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/ecuyle/gomine/internal/authentication"
	"github.com/gin-gonic/gin"
)

func TokenValid(c *gin.Context) bool {
	rawToken := ExtractToken(c)

	return authentication.CheckTokenValidity(rawToken)
}

func ExtractToken(c *gin.Context) string {
	return authentication.ExtractTokenFromBearerToken(c.Request.Header.Get("Authorization"))
}

func ExtractTokenID(c *gin.Context) (string, error) {
	rawToken := ExtractToken(c)

	return authentication.ExtractUserIdFromToken(rawToken)
}

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

	user := User{}
	err = statement.QueryRow(username).Scan(&user.ID, &user.Username, &user.Hash)

	if err != nil {
		return "", err
	}

	err = authentication.ComparePasswordwithHash(password, user.Hash)

	if err != nil {
		return "", err
	}

	token, err := authentication.GenerateToken(user.ID)

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
