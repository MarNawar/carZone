// package middleware

// import (
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	jwt "github.com/golang-jwt/jwt/v5"
// )

// type Claims struct {
// 	UserName string `json:username`
// 	jwt.RegisteredClaims
// }
// var SECRET_KEY string = os.Getenv("SECRET_KEY")

// func AuthMiddleware()gin.HandlerFunc{
// 	return func(c *gin.Context){
// 		clientToken := c.Request.Header.Get("token")
// 		if clientToken == "" {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization header provided"})

// 			c.Abort()
// 			return
// 		}

// 		claims := &Claims{}

// 		token, err := jwt.ParseWithClaims(
// 			clientToken,
// 			claims,
// 			func(token *jwt.Token) (interface{}, error) {
// 				return []byte(SECRET_KEY), nil
// 			},
// 		)

// 		if err!= nil || !token.Valid{
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}

// 		context.WithValue()
// 	}
// }

package middleware

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

var jwtKey = []byte(SECRET_KEY)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func GenerateToken(userName string)(string, error){
	expiration := time.Now().Local().Add(time.Hour * 24)

	// claims := &jwt.RegisteredClaims{
	// 	ExpiresAt: jwt.NewNumericDate(expiration),
	// 	Subject: userName,
	// 	IssuedAt: jwt.NewNumericDate(time.Now()),
	// }

	claims := &Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtKey)
	if err != nil{
		return "", err
	}
	return signedToken, nil
}