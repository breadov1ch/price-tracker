package handlers

import (
	"net/http"
	"price-tracker/internal/database"
	"price-tracker/internal/models"
	"price-tracker/internal/worker"
	"time"

	"github.com/gin-gonic/gin"
)

type AddProductInput struct {
	Name string `json:"name" binding:"required"`
	URL  string `json:"url" binding:"required,url"`
}

// AddProduct godoc
// @Summary      Додати товар для відстеження
// @Description  Додає посилання на товар до акаунта поточного користувача. Якщо товар новий — створює його.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        input body      AddProductInput  true  "URL товару"
// @Success      201  {object}   models.Product   "Товар успішно додано"
// @Failure      400  {object}   map[string]string "Некоректний URL"
// @Failure      401  {object}   map[string]string "Неавторизовано"
// @Security     BearerAuth
// @Router       /products [post]
func AddUserProduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	var input AddProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	var product models.Product
	result := database.DB.Where("url = ?", input.URL).First(&product)

	if result.Error != nil {
		product = models.Product{
			URL:       input.URL,
			LastPrice: 0,
			ImageURL:  "",
			Name:      input.Name,
			UpdatedAt: time.Now(),
		}

		if err := database.DB.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with create product"})
			return
		}
		go worker.StartPriceTicker(&product)
	}
	err := database.DB.Model(&user).Association("Products").Append(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with save product to database"})
	}

	c.JSON(http.StatusCreated, product)
}

// DeleteUserProduct godoc
// @Summary      Видалити товар із відстеження
// @Description  Видаляє зв'язок між поточним користувачем та товаром за його ID. Якщо товар більше ніхто не відстежує, він повністю видаляється з бази даних.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID товару, який потрібно видалити"
// @Success      200  {object}   map[string]string "Товар успішно видалено з відстеження"
// @Failure      400  {object}   map[string]string "Некоректний ID"
// @Failure      401  {object}   map[string]string "Неавторизовано"
// @Failure      404  {object}   map[string]string "Товар або користувача не знайдено"
// @Security     BearerAuth
// @Router       /products/{id} [delete]
func DeleteUserProduct(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID is required"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	err := database.DB.Model(&user).Association("Products").Delete(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with delete product from database"})
		return
	}

	var count int64
	database.DB.Table("user_products").Where("product_id = ?", product.ID).Count(&count)

	if count == 0 {
		if err := database.DB.Delete(&product).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"message": "Product successfully untracked"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product successfully untracked"})
}

// GetUserProducts godoc
// @Summary      Отримати товари користувача
// @Description  Повертає список усіх товарів, які відстежує поточний авторизований користувач
// @Tags         products
// @Accept       json
// @Produce      json
// @Success      200  {array}    models.Product   "Список товарів"
// @Failure      401  {object}   map[string]string "Неавторизовано"
// @Security     BearerAuth
// @Router       /products [get]
func GetUserProducts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User is not authorized"})
	}
	var user models.User
	if err := database.DB.Preload("Products").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error with gettng user's products"})
		return
	}
	c.JSON(http.StatusOK, user.Products)
}
