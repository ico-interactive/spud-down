package main

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	. "spud-down/types"

	dl_api "github.com/deadlock-api/openapi-clients/go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var dl_client *dl_api.APIClient

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

func parseInt32(s string) (int32, error) {
	val, err := strconv.ParseInt(s, 10, 32)
	return int32(val), err
}

func getMatchHistory(c *gin.Context) {
	accountIdParam := c.Param("accountId")
	accountId, err := parseInt32(accountIdParam)
	if err != nil {
		c.Error(errors.New("could not find account id"))
	}
	matches, _, err := dl_client.PlayersAPI.MatchHistory(context.Background(), accountId).ForceRefetch(true).Execute()
	if err != nil {
		c.Error(errors.New("could not fetch from deadlock-api"))
	}
	c.JSON(http.StatusOK, matches)
}

func main() {

	// load env
	env, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal("error: could not load .env file")
	}
	sl_key := env["SL_API_KEY"]
	sl_url := env["SL_API_URL"]
	sl_mock_url := env["SL_MOCK_API_URL"]
	dl_api_url := env["DL_API_URL"]
	gin_port := cmp.Or(env["GIN_PORT"], "8080")
	_, _, _ = sl_key, sl_url, sl_mock_url

	// deadlock-api
	dl_cfg := dl_api.NewConfiguration()
	dl_cfg.Host = dl_api_url
	dl_client = dl_api.NewAPIClient(dl_cfg)

	// gin http server
	router := gin.Default()
	router.Use(ErrorHandler())
	router.GET("/profile/:accountId", getProfileLocal)
	router.GET("/player/history/:accountId", getMatchHistory)
	http.ListenAndServe(fmt.Sprintf(":%s", gin_port), router)

}
