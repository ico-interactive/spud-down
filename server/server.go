package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	. "spud-down/types"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getProfileLocal(c *gin.Context) {
	// load test data
	id := c.Param("accountId")
	dat, readErr := os.ReadFile("./test/data/profiles/" + id)
	if readErr != nil {
		c.Error(errors.New("could not read profile data"))
	}

	// decode json
	var profile Profile
	decoder := json.NewDecoder(bytes.NewReader(dat))
	if err := decoder.Decode(&profile); err == io.EOF {
		return
	} else if err != nil {
		c.Error(errors.New("could not parse profile data"))
	}
	c.JSON(http.StatusOK, profile)
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}

func main() {

	// load env
	env, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("error: could not load .env file")
	}
	apiKey := env["API_KEY"]
	apiURL := env["API_URL"]
	ginPort := cmp.Or(env["GIN_PORT"], "8080")
	_, _ = apiKey, apiURL

	// gin http server
	router := gin.Default()
	router.Use(ErrorHandler())
	router.GET("/profile/:accountId", getProfileLocal)
	http.ListenAndServe(fmt.Sprintf(":%s", ginPort), router)

}
