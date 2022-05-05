package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/lib"
	"authentication/models"
	"errors"
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
		correct, err := lib.CheckPassword(loginRequest.Password, userName.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if correct == true {
			token, err := lib.CreateToken(userName.Username, userName.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while generating token!",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   token,
			})
			return
		}

	}
	if idt == "email" {
		userEmail := models.Users{
			Email: string(loginRequest.Identifier),
		}

		resultUserEmail := config.App.D.Where(&userEmail).First(&userEmail)
		errors.Is(resultUserEmail.Error, gorm.ErrInvalidValue)

		if resultUserEmail.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "User not found!",
			})
			return
		}

		// Check if password is correct
		correct, err := lib.CheckPassword(loginRequest.Password, userEmail.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
			return
		}

		if correct == true {
			token, err := lib.CreateToken(userEmail.Username, userEmail.Password)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error while generating token!",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successful",
				"token":   token,
			})
			return
		}
	}

}

func CheckToken(c *gin.Context) {
	// Getting the token from the request body
	var tokenRequest models.TokenRequest
	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	// Checking if the token is valid
	valid, err := lib.ValidateToken(tokenRequest.Token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if valid == true {
		c.JSON(http.StatusOK, gin.H{
			"message": "Token is valid",
			"status":  "OK",
		})
	}

}
