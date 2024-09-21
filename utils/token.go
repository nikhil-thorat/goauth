package utils

import "net/http"

func GetTokenFromRequest(r *http.Request) string {

	authToken := r.Header.Get("Authorization")

	if authToken != "" {
		return authToken
	}

	return ""
}
