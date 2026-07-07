package worker

import (
	"log"
	"price-tracker/internal/database"
	"price-tracker/internal/models"
	"time"
)

func StartPriceTicker(product *models.Product) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	log.Printf("[ВОРКЕР] Запущен бесконечный трекинг для товара ID %d (%s)\n", product.ID, product.Name)

	shouldContinue := checkAndUpdatePrice(product)
	if !shouldContinue {
		return
	}

	for range ticker.C {
		shouldContinue := checkAndUpdatePrice(product)
		if !shouldContinue {
			return
		}
	}
}

func checkAndUpdatePrice(product *models.Product) bool {
	var dbProduct models.Product
	err := database.DB.First(&dbProduct, product.ID).Error
	if err != nil {
		log.Printf("[ВОРКЕР СТОП] Товар ID %d больше не существует в БД. Останавливаю горутину.\n", product.ID)
		return false
	}
	oldPrice := dbProduct.LastPrice

	productData, err := parseProduct(product.URL)
	if err != nil {
		log.Printf("[ВОРКЕР ОШИБКА] Не удалось обновить данные для ID %d: %v\n", product.ID, err)
		return true
	}

	price := productData.Price

	product.LastPrice = price
	product.UpdatedAt = time.Now()

	product.ImageURL = productData.ImageURL

	database.DB.Save(&product)

	if oldPrice != price {
		sendMessageToGmail(product, oldPrice, price)
		log.Printf("[ВОРКЕР ИЗМЕНЕНИЕ ЦЕНЫ] Товар ID %d: старая цена %.2f грн, новая цена %.2f грн\n", product.ID, oldPrice, price)
	}

	log.Printf("[ВОРКЕР УСПЕХ] Товар ID %d обновлен. Цена: %.2f грн. Картинка: %s\n", product.ID, price, productData.ImageURL)
	return true
}

func InitPriceTrackers() {
	var products []models.Product

	if err := database.DB.Find(&products).Error; err != nil {
		log.Fatalf("[КРИТИЧЕСКАЯ ОШИБКА] Не удалось загрузить товары для трекинга: %v\n", err)
	}

	log.Printf("[ИНИЦИАЛИЗАЦИЯ] Найдено %d товаров в БД. Запускаю фоновые воркеры...\n", len(products))

	for i := range products {
		go StartPriceTicker(&products[i])
	}

	log.Println("[ИНИЦИАЛИЗАЦИЯ] Все фоновые воркеры успешно восстановлены!")
}
