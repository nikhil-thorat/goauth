package auth

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nikhil-thorat/goauth/configs"
	"github.com/nikhil-thorat/goauth/types"
	"github.com/nikhil-thorat/goauth/utils"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := utils.GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			err := fmt.Errorf("FAILED TO VALIDATE TOKEN : %v", err)
			permissionDenied(w, err)
			return
		}

		if !token.Valid {
			err := fmt.Errorf("INVALID TOKEN")
			permissionDenied(w, err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		str := claims["userID"].(string)

		userID, err := strconv.Atoi(str)
		if err != nil {
			err := fmt.Errorf("FAILED TO CONVERT userID to INT : %v", err)
			permissionDenied(w, err)
			return
		}

		user, err := store.GetUserByID(userID)
		if err != nil {
			err := fmt.Errorf("FAILED TO GET USER BY ID : %v", err)
			permissionDenied(w, err)
			return
		}

		if !user.IsVerified {
			err := fmt.Errorf("USER EMAIL NOT VERIFIED")
			permissionDenied(w, err)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func CreateJWT(secret []byte, userID int) (string, error) {

	expiration := time.Second * time.Duration(configs.Envs.JWTExpiration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(int(userID)),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("UNEXPECTED STRING METHOD : %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil

	})
}

func permissionDenied(w http.ResponseWriter, err error) {
	utils.WriteErrorJSON(w, http.StatusForbidden, err)
}

func GetUserIDFromContext(ctx context.Context) int {

	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
