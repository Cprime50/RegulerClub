package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Cprime50/RegulerClub/forms"
	"github.com/Cprime50/RegulerClub/models"
	"github.com/Cprime50/RegulerClub/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

var (
	userModel models.User
)

func Register(c *gin.Context) {
	var input forms.Register
	var errs []string

	// Bind form input as JSON
	if err := c.BindJSON(&input); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationError := range validationErrors {
				if validationError.Tag() == "required" {
					errs = append(errs, fmt.Sprintf("%s not provided", validationError.Field()))
				}
			}
		}
	}

	// Check if email exists
	result, _ := userModel.GetUserByEmail(input.Email)
	if result.Email != "" {
		errs = append(errs, "This email already exists")
	}

	if result.Username != "" {
		errs = append(errs, "This username already exists")
	}

	// Ensure password provided and confirmedPassword match
	if input.Password != input.ConfirmPassword {
		errs = append(errs, "Passwords do not match")
	}

	// If there are any errors, return them
	if len(errs) > 0 {
		for _, err := range errs {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
		}
		return
	}

	// update the user table with new data
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: userModel.SetPassword(input.Password),
		Verified: false,
		RoleID:   4,
	}
	// check if email already in use
	// hashpassword

	// create user
	err := models.CreateUser(&user)
	if err != nil {
		log.Fatal("Unable to create accoount", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Problem creating an account", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Your account has been created succesfully.", "User": input.Email})

	// send code after registeration
	_err := SendCode(input.Email)
	if _err != nil {
		c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"message": "Problem sending you a confirmation code to your mail", "error": _err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Kindly check your mail for a confirmation code", "User": input.Email})
}

// Login
func Login(c *gin.Context) {
	var input forms.Login

	// Bind input
	if err := c.BindJSON(&input); err != nil {
		// check for error in required field
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}

		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage, "message": "please correctly provide the relevant fields"})
		return
	}

	// Check to make sure user provides atleast username and password
	if input.Email == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Provide either an email or username"})
		return
	}

	// check fir user by their email
	user, err := userModel.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "email not found, please create an account", "error": err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong", "error": err.Error()})
		return
	}

	// check if password matches
	_, err = userModel.CheckPassword(input.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "incorrect password provided", "error": err.Error()})
		return
	}

	if user.Verified == false {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Please verify your email to login"})
		return
	}

	// generate jwt
	jwt, err := util.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid token", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"jwt": jwt, "Email": input.Email, "message": "Logged in successfully"})

}

//TODO
// GET all users
// Get user
// Update user
// delete user
// assign role to user
// ban user
