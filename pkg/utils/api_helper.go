package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type Error struct {
	StatusCode int
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

func SuccessMessage(message string, title string, code int) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Title:   title,
			Message: message,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_, err = w.Write(jsonResponse)
		if err != nil {
			return
		}
	}
}

func ErrorMessage(title string, message string, code int) func(c *gin.Context, err error) {

	return func(c *gin.Context, err error) {
		response := struct {
			Title   string `json:"title"`
			Message string `json:"message"`
			Code    int    `json:"code"`
		}{
			Title:   title,
			Message: message,
			Code:    code,
		}
		if err != nil {
			PrintError(err)
		}
		c.JSON(code, response)
	}
}

func PrintError(err error) {
	fmt.Println("\033[31m" + err.Error() + "\033[0m")
}
