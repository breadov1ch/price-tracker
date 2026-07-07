package worker

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"price-tracker/internal/database"
	"price-tracker/internal/models"

	"gopkg.in/gomail.v2"
)

func sendMessageToGmail(product *models.Product, oldPrice, newPrice float64) {
	var users []models.User
	err := database.DB.Model(product).Association("Users").Find(&users)
	if err != nil {
		log.Printf("Error fetching users for product %s: %v", product.URL, err)
		return
	}

	tmpl, err := template.ParseFiles("templates/price_change.html")
	if err != nil {
		log.Printf("Error parsing email template: %v", err)
		return // КРИТИЧЕСКОЕ ИСПРАВЛЕНИЕ: выходим, чтобы не поймать панику ниже
	}

	// Выносим инициализацию диалера за пределы цикла.
	// Зачем подключаться к SMTP-серверу заново для каждого пользователя, верно?
	// Один диалер отправит всё гораздо быстрее.
	d := gomail.NewDialer("smtp.gmail.com", 587, "server.breadmeet@gmail.com", os.Getenv("GMAIL_APP_PASSWORD"))

	for _, user := range users {
		data := map[string]any{
			"UserName":    user.Name,
			"ProductName": product.Name,
			"OldPrice":    oldPrice,
			"NewPrice":    newPrice,
			"URL":         product.URL,
			"ImageURL":    product.ImageURL, // <-- ДОБАВИЛИ КАРТИНКУ ДЛЯ ШАБЛОНА
		}

		var body bytes.Buffer
		if err := tmpl.Execute(&body, data); err != nil {
			log.Printf("Ошибка заполнения шаблона для %s: %v", user.Email, err)
			continue // Используем continue вместо return, чтобы ошибка одного юзера не ломала отправку остальным
		}

		m := gomail.NewMessage()
		m.SetHeader("From", "server.breadmeet@gmail.com")
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "💰 Цена на товар изменилась!")
		m.SetBody("text/html", body.String())

		if err := d.DialAndSend(m); err != nil {
			log.Printf("Ошибка отправки письма для %s: %v", user.Email, err)
		} else {
			log.Printf("Письмо успешно отправлено для %s", user.Email)
		}
	}
}
