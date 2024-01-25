package helper

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
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

func AddDecimal(saldo float64, requestGram string) float64 {
	decimalValue, err := decimal.NewFromString(requestGram)
	if err != nil {
		fmt.Println("Error parsing string to decimal:", err)
		return 0
	}

	result := decimalValue.Add(decimal.NewFromFloat(saldo))

	resultFloat, _ := result.Float64()
	return resultFloat
}

func DecimalFromString(gram string) (float64, error) {
	decimalValue, err := decimal.NewFromString(gram)
	if err != nil {
		fmt.Println("Error parsing string to decimal:", err)
		return 0, err
	}

	result := decimalValue.Add(decimal.NewFromFloat(0))

	resultFloat, _ := result.Float64()
	return resultFloat, nil
}

func ValidateGram(gram string) bool {

	resultFloat, _ := DecimalFromString(gram)
	if hasMoreThan3DecimalDigits(resultFloat) {
		fmt.Println("Error: Value should have at most 3 digits after the decimal point")
		return false
	}

	return true
}

func removeTrailingZeros(str string) string {
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return str
	}

	return strconv.FormatFloat(num, 'f', -1, 64)
}

func hasMoreThan3DecimalDigits(num float64) bool {
	numStr := removeTrailingZeros(fmt.Sprintf("%f", num))

	dotIndex := strings.Index(numStr, ".")

	if len(numStr)-dotIndex-1 > 3 {
		return true
	}

	return false
}

func ValidationFormatError(err error) interface{} {
	var validateErr validator.ValidationErrors
	if errors.As(err, &validateErr) {
		for _, e := range err.(validator.ValidationErrors) {
			return e.Error()
		}
	}
	return "Invalid Input"
}

func TimeNow() int {
	currentTime := time.Now()

	epochTimeSeconds := currentTime.Unix()
	return int(epochTimeSeconds)
}
