package helper

import (
	"log"
	"time"

	"github.com/teris-io/shortid"
)

type ResponseError struct {
	Error   bool        `json:"error"`
	ReffId  string      `json:"reff_id"`
	Message interface{} `json:"message"`
}

func APIResponseError(status bool, reff_id string, message interface{}) ResponseError {

	jsonResponse := ResponseError{
		Error:   status,
		ReffId:  reff_id,
		Message: message,
	}

	return jsonResponse
}

func GenShortId() string {
	sid, err := shortid.New(1, shortid.DefaultABC, 2342)
	shortid.SetDefault(sid)
	shortId, err := shortid.Generate()

	if err != nil {
		log.Fatalf("Failed to generate id: %v", err)
	}

	return shortId
}

func TimeNow() int {
	currentTime := time.Now()

	epochTimeSeconds := currentTime.Unix()
	return int(epochTimeSeconds)
}
