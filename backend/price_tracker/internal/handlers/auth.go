package handlers

import (
	"fmt"
	"net/http"
	"price-tracker/internal/database"
	"price-tracker/internal/models"
	"time"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary      Реєстрація нового користувача
// @Description  Створює новий акаунт користувача в базі даних
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body      RegisterInput  true  "Дані для реєстрації"
// @Success      21 stage_success  {object}       map[string]interface{} "Успішна реєстрація"
// @Failure      400             {object}       map[string]string      "Помилка у запиті або дублікат Email"
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	newUser := models.User{
		Email:     input.Email,
		Name:      input.Name,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}
	result := database.DB.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email alreade registred"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": newUser.ID})
}

// Login godoc
// @Summary      Авторизація користувача
// @Description  Перевіряє пошту/пароль та повертає JWT-токен
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body      LoginInput  true  "Дані для входу"
// @Success      200             {object}    map[string]interface{} "Успішний вхід, повертає токен"
// @Failure      400             {object}    map[string]string      "Неправильний пароль або логін"
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password must be required"})
		return
	}

	var user models.User
	result := database.DB.Where("email = ?", input.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie("refresh_token", refreshTokenString, 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message":      "Login was successful",
		"access_token": accessTokenString,
	})
}

func Refresh(c *gin.Context) {
	cookieToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token не знайдено, увійдіть знову"})
		return
	}

	token, err := jwt.Parse(cookieToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Недійсний або протухлий refresh токен"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не вдалося прочитати дані"})
		return
	}

	userID := claims["user_id"].(float64)

	newAccessClaims := jwt.MapClaims{
		"user_id": uint(userID),
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	}
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	newAccessTokenString, err := newAccessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessTokenString,
	})
}
