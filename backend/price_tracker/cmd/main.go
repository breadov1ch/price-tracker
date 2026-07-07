package main

import (
	"log"
	"price-tracker/internal/database"
	"price-tracker/internal/handlers"
	"price-tracker/internal/middleware"
	"price-tracker/internal/worker"
	"strings"

	_ "price-tracker/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Price Tracker API
// @version         1.0
// @description     API сервер для відстеження цін на товари.
// @host            localhost:8080
// @BasePath        /api

// @securityDefinitions.apiKey BearerAuth
// @in                         header
// @name                       Authorization
// @description                Введіть токен у форматі: Bearer <ваш_токен>
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Попередження: Не вдалося завантажити файл .env")
	}
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowOriginFunc = func(origin string) bool {
		return strings.HasPrefix(origin, "http://localhost") || strings.HasPrefix(origin, "http://127.0.0.1")
	}

	r.Use(cors.New(corsConfig))
	database.InitDB()
	worker.InitPriceTrackers()

	// testUser := models.User{
	// 	Email:    "gmail.com",
	// 	Name:     "Danilo Test",
	// 	Password: "supersecretpassword",
	// }

	// // Зберігаємо в SQLite через GORM
	// // FirstOrCreate перевірить, чи є вже такий юзер, якщо немає — створить
	// result := database.DB.FirstOrCreate(&testUser, models.User{Email: "gmail.com"})
	// if result.Error != nil {
	// 	log.Println("Помилка запису в БД:", result.Error)
	// } else {
	// 	log.Println("База даних працює! Тестовий юзер ID:", testUser.ID)
	// }

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)
	r.POST("/api/auth/refresh", handlers.Refresh)
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/products", handlers.AddUserProduct)
		protected.GET("/products", handlers.GetUserProducts)
		protected.DELETE("/products/:id", handlers.DeleteUserProduct)
	}

	r.Run(":8080")
}
