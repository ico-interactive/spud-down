package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Profile struct {
	AccountID                 int64   `json:"accountId"`
	Name                      string  `json:"name"`
	AvatarURL                 string  `json:"avatarUrl"`
	PerformanceRankMessage    *string `json:"performanceRankMessage"`
	LastUpdated               string  `json:"lastUpdated"`
	CalibrationMatches        int     `json:"calibrationMatches"`
	TwitchUsername            *string `json:"twitchUsername"`
	YoutubeChannelURL         *string `json:"youtubeChannelUrl"`
	Region                    string  `json:"region"`
	CalibrationMatchID        int64   `json:"calibrationMatchId"`
	CalibrationResetMatchID   *int64  `json:"calibrationResetMatchId"`
	FontID                    string  `json:"fontId"`
	GlowStyleID               string  `json:"glowStyleId"`
	GlowColorID               string  `json:"glowColorId"`
	OutlineID                 string  `json:"outlineId"`
	AnimationID               string  `json:"animationId"`
	AvatarShapeID             string  `json:"avatarShapeId"`
	AvatarEffectID            string  `json:"avatarEffectId"`
	AvatarEffectColorID       string  `json:"avatarEffectColorId"`
	TitleID                   string  `json:"titleId"`
	UpdatedWithinLast1Minutes bool    `json:"updatedWithinLast1Minutes"`
	LocalisedLastUpdated      string  `json:"localisedLastUpdated"`
	PPScore                   int     `json:"ppScore"`
	SelectedWidgets           *string `json:"selectedWidgets"`
	EstimatedRankNumber       int     `json:"estimatedRankNumber"`
}

func check(e error) {
	if e != nil {
		panic(e)
		// log.Fatal(e)
	}
}

func getProfile(c *gin.Context) {
	dat, err := os.ReadFile("./test/data/profiles/151453298")
	fmt.Printf("Read example profile JSON file: %d bytes\n", len(dat))
	check(err)
	decoder := json.NewDecoder(bytes.NewReader(dat))
	fmt.Println("Decoding example profile JSON stream...")
	for {
		var profile Profile
		err := decoder.Decode(&profile)
		check(err)
		if err == io.EOF {
			break
		}
		c.IndentedJSON(http.StatusOK, profile)
	}
}

func main() {

	// load env
	env, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal("error: could not load .env file")
	}
	apiKey := env["API_KEY"]
	apiURL := env["API_URL"]
	ginPort := cmp.Or(env["GIN_PORT"], "8080")
	_, _ = apiKey, apiURL

	// gin http server
	router := gin.Default()
	router.GET("/profile", getProfile)
	http.ListenAndServe(fmt.Sprintf(":%s", ginPort), router)

}
