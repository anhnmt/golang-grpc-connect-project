package utils

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	// HeaderAuthorize header authorize
	HeaderAuthorize = "authorization"
)

// TotalPage returns the total number of pages.
func TotalPage(total int64, pageSize int64) (totalPage int64) {
	if total%pageSize == 0 {
		totalPage = total / pageSize
	} else {
		totalPage = total/pageSize + 1
	}

	if totalPage == 0 {
		totalPage = 1
	}

	return
}

// CurrentPage returns the current page.
func CurrentPage(page int64, totalPages int64) int64 {
	if page <= 0 || totalPages < page {
		return 1
	} else if page > totalPages {
		return totalPages
	}

	return page
}

// AuthFromHeader is a helper function for extracting the :authorization header from the http header of the request.
//
// It expects the `:authorization` header to be of a certain scheme (e.g. `basic`, `bearer`), in a
// case-insensitive format (see rfc2617, sec 1.2). If no such authorization is found, or the token
// is of wrong scheme, an error with status `Unauthenticated` is returned.
func AuthFromHeader(header http.Header, expectedScheme string) (string, error) {
	val := header.Get(HeaderAuthorize)
	if val == "" {
		return "", fmt.Errorf("Request unauthenticated with " + expectedScheme)
	}

	splits := strings.SplitN(val, " ", 2)
	if len(splits) < 2 {
		return "", fmt.Errorf("Bad authorization string")
	}

	if !strings.EqualFold(splits[0], expectedScheme) {
		return "", fmt.Errorf("Request unauthenticated with " + expectedScheme)
	}

	return splits[1], nil
}
