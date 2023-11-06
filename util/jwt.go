package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cprime50/RegulerClub/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// generate token
func GenerateJWT(user *models.User) (string, error) {
	tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.RoleID,
		// TODO make token expiry loaction based instead of unix
		"issued_at":  time.Now().Unix(),
		"expires_at": time.Now().Add(time.Minute * time.Duration(tokenTTL)).Unix(),
	})
	// sign token with secret key encoding
	return token.SignedString(privateKey)
}

// validate token
func ValidateJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// Validate admin role
func ValidateAdminRoleJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 1 {
		return nil
	}
	return errors.New("invalid admin token provided")
}

// Validate schooluser role
func ValidateModRoleJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 2 || userRole == 1 {
		return nil
	}
	return errors.New("invalid MOD token provided")
}

// Validate paidUser role
func ValidatePrimeUserRoleJWT(c *gin.Context) error {
	token, err := getToken(c)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 3 || userRole == 1 {
		return nil
	}
	return errors.New("invalid paidUser token provided")
}

// fetch user details from token
func CurrentUser(c *gin.Context) *models.User {
	var userModel models.User
	err := ValidateJWT(c)
	if err != nil {
		return &models.User{}
	}
	token, _ := getToken(c)
	claims, _ := token.Claims.(jwt.MapClaims)
	userId := uint(claims["id"].(float64))

	user, err := userModel.GetUserById(userId)
	if err != nil {
		return &models.User{}
	}
	return user
}

// TODO  improve if necessary
func getToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

// extract token from request authorization header
func getTokenFromRequest(c *gin.Context) string {
	bearerToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
