package utils

import (
	"net/mail"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetPagination(ctx *gin.Context) (int, int, error) {
	current_page := ctx.Query("current_page")
	page_size := ctx.Query("page_size")

	limit := 10
	if page_size != "" {
		page_size_int, err := strconv.Atoi(page_size)
		if err != nil {
			return 0, 0, err
		}
		limit = page_size_int
	}

	offset := 0
	if current_page != "" {
		current_page_int, err := strconv.Atoi(current_page)
		if err != nil {
			return 0, 0, err
		}
		offset = (current_page_int - 1) * limit
	}
	return limit, offset, nil
}

func IsValidUrl(str string) bool {
	_, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return true
}

func IsValidEmail(str string) bool {
	_, err := mail.ParseAddress(str)
	if err != nil {
		return false
	}
	return true
}
