package worker

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

type ExtractMethod string

const (
	FromText      ExtractMethod = "text"      // Из текста внутри тега (например, <span>1500</span>)
	FromAttribute ExtractMethod = "attribute" // Из значения атрибута (например, data-qaprice="1500")
)

type PriceInstruction struct {
	Selector      string
	Method        ExtractMethod
	AttrKey       string
	ImageSelector string
	ImageMethod   ExtractMethod
	ImageAttrKey  string
}

type ProductData struct {
	Price    float64
	ImageURL string
}

var ShopRegistry = map[string]PriceInstruction{
	"prom.ua": {
		Selector:      "div[data-qaid='product_price']",
		Method:        FromAttribute,
		AttrKey:       "data-qaprice",
		ImageSelector: "img.XUh3F.o91O9",
		ImageMethod:   FromAttribute,
		ImageAttrKey:  "src",
	},
	"comfy.ua": {
		Selector:      "meta[name='product:sale_price:amount']",
		Method:        FromAttribute,
		AttrKey:       "content",
		ImageSelector: "img.contain",
		ImageMethod:   FromAttribute,
		ImageAttrKey:  "src",
	},
	"grokholsky.com": {
		Selector:      "div.page-product",
		Method:        FromAttribute,
		AttrKey:       "data-discount",
		ImageSelector: "picture img",
		ImageMethod:   FromAttribute,
		ImageAttrKey:  "src",
	},
}

func extractDomain(productURL string) (string, error) {
	parsedURL, err := url.Parse(productURL)
	if err != nil {
		return "", err
	}
	host := strings.TrimPrefix(parsedURL.Hostname(), "www.")
	return host, nil
}

func sanitizePrice(raw string) string {
	res := strings.Join(strings.Fields(raw), "")
	res = strings.ReplaceAll(res, "\u00a0", "")
	res = strings.ReplaceAll(res, "грн", "")
	res = strings.ReplaceAll(res, "₴", "")
	res = strings.ReplaceAll(res, ",", ".")
	return strings.TrimSpace(res)
}

func parseProduct(productURL string) (ProductData, error) {
	var data ProductData

	domain, err := extractDomain(productURL)
	if err != nil {
		return data, fmt.Errorf("ошибка парсинга URL: %w", err)
	}

	instruction, exists := ShopRegistry[domain]
	if !exists {
		return data, fmt.Errorf("магазин %s еще не добавлен в список инструкций", domain)
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)
	c.SetRequestTimeout(15 * time.Second)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Language", "uk-UA,uk;q=0.9,en-US;q=0.8,en;q=0.7")
		r.Headers.Set("Cache-Control", "no-cache")
	})

	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			log.Printf("[Внимание] Сайт вернул статус: %d\n", r.StatusCode)
		}
		if strings.Contains(string(r.Body), "cloudflare") || strings.Contains(string(r.Body), "Just a moment") {
			log.Println("[БЛОКИРОВКА] Сайт включил защиту от роботов (Cloudflare)!")
		}
	})

	var parseErr error
	var tagFound bool

	c.OnHTML(instruction.Selector, func(e *colly.HTMLElement) {
		tagFound = true
		var rawPrice string

		if instruction.Method == FromAttribute {
			rawPrice = e.Attr(instruction.AttrKey)
		} else {
			rawPrice = e.Text
		}

		cleanPrice := sanitizePrice(rawPrice)
		parsedPrice, err := strconv.ParseFloat(cleanPrice, 64)
		if err != nil {
			parseErr = fmt.Errorf("не удалось перевести '%s' в число: %w", cleanPrice, err)
			return
		}
		data.Price = parsedPrice
	})

	c.OnHTML(instruction.ImageSelector, func(e *colly.HTMLElement) {
		var rawImg string
		if instruction.ImageMethod == FromAttribute {
			rawImg = e.Attr(instruction.ImageAttrKey)
		} else {
			rawImg = e.Text
		}

		// AbsoluteURL преобразует относительные ссылки /uploads/... в абсолютные
		data.ImageURL = e.Request.AbsoluteURL(rawImg)
	})

	err = c.Visit(productURL)
	if err != nil {
		return data, fmt.Errorf("ошибка при запросе к сайту: %w", err)
	}

	if parseErr != nil {
		return data, parseErr
	}

	if !tagFound {
		return data, fmt.Errorf("селектор цены '%s' не найден на странице", instruction.Selector)
	}

	return data, nil
}

func main() {
	urlsFromDB := []string{
		"https://prom.ua/ua/p3010567509-setevoe-zaryadnoe-ustrojstvo.html",
		"https://comfy.ua/ua/holodil-nik-bosch-kgn39ul316.html",
		"https://grokholsky.com/ua/product/apple/mac/apple-macbook-air/macbook-air-13-6-m5-2026/noutbuk-apple-macbook-air-13-6-2026-apple-m5-10-core-gpu-16gb-1tb-midnight-mdhf4/",
	}

	fmt.Println("--- Запуск парсера товаров ---")

	for _, u := range urlsFromDB {
		fmt.Printf("Сканирую ссылку: %s\n", u)

		product, err := parseProduct(u)
		if err != nil {
			log.Printf("[ОШИБКА] %v\n\n", err)
			continue
		}

		fmt.Printf("[УСПЕХ] Актуальная цена товара: %.2f грн\n", product.Price)
		fmt.Printf("[УСПЕХ] URL картинки товара: %s\n\n", product.ImageURL)
	}
}
