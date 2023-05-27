package controllers

import (
	"Auth/initializers"
	"Auth/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

var body struct {
	Email    string
	Password string
}

func SignUp(c *gin.Context) {
	//	Get the Email/pass off the request body

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}
	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Hash password",
		})
		return
	}

	//Create the User
	user := models.User{Email: body.Email, Password: string(hash)}
	result := initializers.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to Create User",
		})
		return
	}
	//Response
	c.JSON(http.StatusOK, gin.H{})
}

func Login(c *gin.Context) {
	//Get the data in the request
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	//Look up user
	var user models.User
	initializers.Db.First(&user, "email= ?", body.Email)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No User with this email",
		})
		return
	}
	//	Compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong password",
		})
		return
	}

	//	Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	s := os.Getenv("SECRET")
	// if s is send as string without converting it will result in an error
	//because sugnedstring accepts interface but the key in HMAC must be of type []byte refer to documentation to know
	tokenString, err := token.SignedString([]byte(s))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		fmt.Println(err)
		return
	}

	// send it back
	c.SetSameSite(http.SameSiteLaxMode)
	// Normaly the secure argument is set to true but since this is our localhost we set it to false
	c.SetCookie("authorization", tokenString, 3600*2, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		//we can either send it back or use a cookie to store it
		//"token": tokenString,
	})
}

// Runs after require auth and it gets the data requireauth set in the context and can manipulate it and use
// it in any way possible(do any operation on this data)
func Validate(c *gin.Context) {
	user, _ := c.Get("user")

	//user.(models.User).Email

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
