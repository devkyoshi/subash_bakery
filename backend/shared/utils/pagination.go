package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPageParam extracts and validates the page parameter from the request
// Defaults to 1 if not provided or invalid
func GetPageParam(c *gin.Context) int {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// GetLimitParam extracts and validates the limit parameter from the request
// Defaults to 10 if not provided or invalid
// Maximum limit is 100
func GetLimitParam(c *gin.Context) int {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		return 10
	}
	if limit > 100 {
		return 100
	}
	return limit
}
