package login

import (
	"log"
	"net/http"

	"github.com/MarNawar/carZone/middleware"
	"github.com/MarNawar/carZone/models"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var user models.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	valid := (user.UserName == "admin" && user.Password == "admin123")

	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "provide valid user name or password"})
	}

	tokenString, err := middleware.GenerateToken(user.UserName)
	if err != nil {
		log.Println("Error Generating Token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Generate Token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}
