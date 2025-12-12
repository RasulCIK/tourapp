package controllers

import (
	"fmt"
	"net/http"
	"time"
	"tourapp/internal/logger"
	"tourapp/internal/models"
	"tourapp/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("my_secret_key")
var validate = validator.New()

type UserController struct {
	UserRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{UserRepo: userRepo}
}

func (u *UserController) Login(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.UserRepo.GetByEmail(input.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

	logger.Logger.Info().
		Str("email", user.Email).
		Time("login_time", time.Now()).
		Msg("User logged in")
}

func (u *UserController) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	fmt.Sscan(idParam, &id)

	user, err := u.UserRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)

	logger.Logger.Info().
		Uint("user_id", id).
		Msg("User data retrieved")
}

func (u *UserController) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	fmt.Sscan(idParam, &id)

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.UserRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Username = input.Username
	user.Email = input.Email
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}

	u.UserRepo.Update(user)
	c.JSON(http.StatusOK, user)

	logger.Logger.Info().
		Uint("user_id", user.ID).
		Str("email", user.Email).
		Msg("User updated")
}

func (u *UserController) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	var id uint
	fmt.Sscan(idParam, &id)

	if err := u.UserRepo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})

	logger.Logger.Info().
		Uint("user_id", id).
		Msg("User deleted")
}

func (u *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	if err := validate.Struct(user); err != nil {
		
		errors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, e := range errors {
			errorMessages = append(errorMessages, fmt.Sprintf("%s %s", e.Field(), e.Tag()))
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMessages})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	if err := u.UserRepo.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot create user"})
		return
	}

	logger.Logger.Info().
		Str("email", user.Email).
		Msg("User registered")

	c.JSON(http.StatusCreated, user)
}
