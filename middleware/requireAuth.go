package middleware

import (
	"Auth/initializers"
	"Auth/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func RequireAuth(c *gin.Context) {
	//Get the cookie of request
	tokenString, err := c.Cookie("authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	//Decode/Validate Cookie
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		//Find the user with token sub
		var user models.User
		initializers.Db.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		//Attach to request
		c.Set("user", user)

		//Continue on from this middleware
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
