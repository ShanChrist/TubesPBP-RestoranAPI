package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("Bebasapasaja123!")
var tokenName = "token"

type Claims struct {
	ID       int    `json:"json:id"`
	Name     string `json:"json:name"`
	UserType int    `json:"json:type"`
	jwt.StandardClaims
}

func generateToken(w http.ResponseWriter, id int, name string, userType int) {
	tokenExpiryTime := time.Now().Add(5 * time.Minute)
	//create claims with user data
	claims := &Claims{
		ID:       id,
		Name:     name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}
	//encrypt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//jwtKey := os.Getenv("JWT_TOKEN")
	jwtToken, err := token.SignedString(jwtKey)
	if err != nil {
		log.Println(err)
		return
	}

	//set token to cookies
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    jwtToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetUserToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func Authenticate(next http.HandlerFunc, accesType int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateUserToken(r, accesType)
		if !isValidToken {
			// sendErrorResponse(w)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accesType int) bool {
	isAccesTokenValid, id, name, userType := validateTokenFromCookies(r)
	fmt.Print(id, name, userType, accesType, isAccesTokenValid)

	if isAccesTokenValid {
		isUserValid := userType == accesType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateTokenFromCookies(r *http.Request) (bool, int, string, int) {
	if cookie, err := r.Cookie(tokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &Claims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken, accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", -1
}
