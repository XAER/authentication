package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/lib"
	"authentication/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandler(c *gin.Context) {
	// Getting the request body as models.RegisterRequest
	var registerRequest models.RegisterRequest

	// Binding the request body to the models.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Hash password
	hashedPassword, err := lib.CreateHashedPassword(registerRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	// Creating the user and inserting it into the database
	// Creating a new user
	user := models.Users{
		Username:  registerRequest.Username,
		Password:  hashedPassword,
		Email:     registerRequest.Email,
		FirstName: registerRequest.FirstName,
		LastName:  registerRequest.LastName,
	}
	// Inserting the user into the database
	result := config.App.D.Create(&user)

	errors.Is(result.Error, gorm.ErrInvalidValue)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User created!",
	})
}

func LoginHandler(c *gin.Context) {
	var loginRequest models.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	// check if user wrote email or username
	idt := helpers.IsEmailOrUsername(loginRequest.Identifier)

	if idt == "username" {
		userName := models.Users{
			Username: string(loginRequest.Identifier),
		}
		resultUserName := config.App.D.Where(&userName).First(&userName)
		errors.Is(resultUserName.Error, gorm.ErrInvalidValue)

		if resultUserName.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found!",
			})
			return
		}

		// Check if password is correct
		fmt.Println(userName.Password)
		correct, err := lib.CheckPassword(loginRequest.Password, userName.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if correct == true {
			c.JSON(http.StatusOK, gin.H{
				"message": "Password correct!",
			})
			return
		}

	}
	if idt == "email" {
		userEmail := models.Users{
			Email: idt,
		}

		resultUserEmail := config.App.D.Where(&userEmail).First(&userEmail)

		errors.Is(resultUserEmail.Error, gorm.ErrInvalidValue)
		if resultUserEmail.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found!",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "User found!",
		})
		return
	}

	// Check if the user exists

}
