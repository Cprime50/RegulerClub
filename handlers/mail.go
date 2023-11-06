package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Cprime50/RegulerClub/forms"
	"github.com/Cprime50/RegulerClub/models"
	"github.com/Cprime50/RegulerClub/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Sends a 6 digit code to provided email
func SendCode(email string) error {
	user, err := userModel.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no matching user found")
		}
		return err
	}
	otp, err := models.CreateOTPForEmailVerification(user)
	if err != nil {
		return errors.New("error creating OTP")
	}
	subject := "Join us at Reguler Club"
	//Body of the mail
	body := "Howdy good fellow, Copy this code " + otp + " to join us at Reguler Club"
	html := "<strong>" + body + "</strong>"

	// Initialize/Send email
	_err := util.SendMail(subject, user.Email, html, user.Username)

	if _err == nil {
		return nil
	} else {
		return _err
	}
}

// resend otp code
func ResendCode(c *gin.Context) {
	var input forms.SendEmail

	if (c.BindJSON(&input)) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Provide all fields"})
		return
	}
	user, err := userModel.GetUserByEmail(input.Email)
	if err != nil {
		log.Println("Can't Resend Code, Unable to get user by email", err)
		if user.Email == "" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User account was not found", "error": err})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something strange went wrong", "error": err})
		return
	}
	err = SendCode(user.Email)
	if err != nil {
		log.Println("Can't Resend Code, something weird happend", err)
		c.AbortWithStatusJSON(http.StatusNotImplemented, gin.H{"message": "Unable to send email", "error": err})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "A mail has been sent to you. Kindly check your inbox"})
}

func VerifyEmail(c *gin.Context) {
	var input forms.VerifyEmail

	if (c.BindJSON(&input)) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Provide all fields"})
		return
	}

	user, err := userModel.GetUserByEmail(input.Email)
	if err != nil {
		log.Println("Cant verify email, Unable to get user by email", err)
		if user.Email == "" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User account was not found", "error": err})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something strange went wrong", "error": err})
		return
	}

	// Get the OTP and check if OTP inputed matches correctly
	userOTP, err := models.GetOTPByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("OTP not found in database", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid OTP, please generate new OTP", "error": err})
			return
		}
		log.Println("something went wrong while trying to get OTP by email", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong", "error": err})
		return
	}

	// confirm OTP matches and has not expired
	__err := Verify(userOTP, input.OTP)
	if err != nil {
		log.Println("something went wrong while trying to verify OTP", __err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user account by updating the verification status in db to true
	_err := userModel.VerifyAccount(user.Email)
	if _err != nil {
		log.Println("something went wrong while trying to verify user account", _err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Error while trying to verify your account", "error": err})
		return
	}
	err_ := models.DeleteOTP(userOTP)
	if err_ != nil {
		log.Println("problem deleting OTP from database", err_)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account verified successfully, Log in"})
}

// This validates if the OTP in the database matches user OTP provided and also deletes expired OTP
func Verify(OTP *models.OTP, inputOTP string) error {
	if inputOTP != OTP.Code {
		log.Println("Verification of mail failed, Invalid verification code", "OTP is :", OTP.Code, "the user input:", inputOTP)
		return errors.New("Verification code provided is Invalid. Please look in your mail and provide the correct code")
	}

	if OTP.ExpiresAt.Before(time.Now()) {
		log.Println("Verification code is expired")
		err := models.DeleteOTP(OTP)
		if err != nil {
			log.Println("problem deleting OTP from database", err)
		}
		return errors.New("Confirmation code has expired. Please try generating a new one")
	}
	return nil
}

// Send OTP for password reset
func SendPasswordResetCode(c *gin.Context) {
	var input forms.SendEmail
	if c.BindJSON(&input) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Provide your email"})
		return
	}

	user, err := userModel.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "no matching user found"})
			return
		}
		log.Println(err, "error trying to fetch user")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something wrong happened", "error": err})
		return
	}
	otp, err := models.CreateOTPForPasswordReset(user)
	if err != nil {
		log.Println(err, "error creating OTP for passwordReset")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "problem generating OTP at the moment"})
		return
	}

	subject := "Reset your password Reguler Club"
	//Body of the mail
	body := "Howdy good fellow, Copy this code " + otp + " to reset your password"
	html := "<strong>" + body + "</strong>"

	// Initialize/Send email
	_err := util.SendMail(subject, user.Email, html, user.Username)
	if _err != nil {
		log.Println(err, "error creating Sending mail for passwordReset")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "problem sending you a mail"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "An email has been sent to you,  please check t continue with password reset"})
}

// Password Reset
func PasswordReset(c *gin.Context) {
	var input forms.ResetPassword
	// Ensure input provided doesnt go against our schema
	if c.BindJSON(&input) != nil {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"message": "Provide relevant fieds"})
		return
	}
	user, err := userModel.GetUserByEmail(input.Email)
	if err != nil {
		log.Println("Cant verify email, Unable to get user by email", err)
		if user.Email == "" || errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "User account was not found", "error": err})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something strange went wrong", "error": err})
		return
	}

	// Get the OTP and check if OTP inputed matches correctly
	userOTP, err := models.GetOTPByEmail(user.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("OTP not found in database", err)
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Invalid OTP, please generate new OTP", "error": err})
			return
		}
		log.Println("something went wrong while trying to get OTP by email", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong", "error": err})
		return
	}

	// confirm OTP matches and has not expired
	err = Verify(userOTP, input.OTP)
	if err != nil {
		log.Println("something went wrong while trying to verify OTP", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	// Ensure passsword provided and confirmedPassword macthes
	if input.Password != input.ConfirmPassword {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Passwords do not match"})
		return
	}

	// Hash new password
	newHashedPassword := userModel.SetPassword(input.Password)

	// Update user account with new password
	_err := userModel.UpdatePassword(user.Email, newHashedPassword)
	if _err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Something wrong happened while updating password", "error": err})
		return
	}

	c.AbortWithStatusJSON(http.StatusCreated, gin.H{"message": "Password has been Updated successfully, log in"})
	return
}
