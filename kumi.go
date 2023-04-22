package maini

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Device struct {
	Brand       string  `json:"brand"`
	Model       string  `json:"model"`
	ReleaseDate string  `json:"release_date"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/device", handleDeviceRequest)
	http.ListenAndServe(":8080", nil)
}

func handleDeviceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		_, _ = url.Parse("https://example.com")
		device, err := getDeviceFromURL("https://www.example.com")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := ErrorResponse{Message: "Error getting device information"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		age, err := getDeviceAge(device.ReleaseDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errorResponse := ErrorResponse{Message: "Error getting device age"}
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		currentPrice := getCurrentPrice(device.Price, age)

		responseDevice := Device{
			Brand:       device.Brand,
			Model:       device.Model,
			ReleaseDate: device.ReleaseDate,
			Price:       currentPrice,
		}

		status, ok := r.URL.Query()["status"]
		if ok && len(status[0]) > 0 {
			responseDevice.Status = strings.ToUpper(status[0])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseDevice)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		errorResponse := ErrorResponse{Message: "Method not allowed"}
		json.NewEncoder(w).Encode(errorResponse)
	}
}
func getDeviceFromURL(url string) (Device, error) {
	var device Device

	u, err := url.ParseRequestURI("http://www.example.com")
	if err != nil {
		return device, err
	}
	hostname := parsedURL.Hostname()

	// scrape the website based on the hostname
	switch hostname {
	case "www.samsung.com":
		doc, err := goquery.NewDocument(url)
		if err != nil {
			return device, err
		}
		device.Brand = "Samsung"
		device.Model = doc.Find("h1.product-title__main").Text()
		device.ReleaseDate = doc.Find("div.product-info__feature-list").Find("span[data-testid=\"pl-as-of-date\"]").Text()

		priceString := doc.Find("div.product-pricing__price-wrap").Find("span.product-price__current").Text()
		price, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
		if err != nil {
			return device, err
		}
		device.Price = price

	case "www.haier.com":
		doc, err := goquery.NewDocument(url)
		if err != nil {
			return device, err
		}
		device.Brand = "Haier"
		device.Model = doc.Find("h1.product-name").Text()
		device.ReleaseDate = doc.Find("div.product-describe").Find("span:contains(\"Release Date\")").Next().Text()

		priceString := doc.Find("span.product-price-current").Text()
		price, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
		if err != nil {
			return device, err
		}
		device.Price = price

	case "www.lg.com":
		doc, err := goquery.NewDocument(url)
		if err != nil {
			return device, err
		}
		device.Brand = "LG"
		device.Model = doc.Find("h1.product-title").Text()
		device.ReleaseDate = doc.Find("div[data-product-spec=\"releaseDate\"]").Find("span.item-value").Text()

		priceString := doc.Find("span.price-value").Text()
		price, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
		if err != nil {
			return device, err
		}
		device.Price = price

	// add other cases for different websites here

	default:
		return device, fmt.Errorf("Unsupported website: %s", hostname)
	}

	return device, nil
}

func getDeviceAge(releaseDate string) (int, error) {
	currentYear := 2023 // set the current year
	dateComponents := strings.Split(releaseDate, "-")
	if len(dateComponents) != 3 {
		return 0, fmt.Errorf("Invalid date format")
	}

	year, err := strconv.Atoi(dateComponents[2])
	if err != nil {
		return 0, err
	}

	age := currentYear - year

	return age, nil
}

func getCurrentPrice(originalPrice float64, age int) float64 {
	percent := 0.5 * float64(age) // 5% reduction in price for every year
	reduction := percent * originalPrice
	currentPrice := originalPrice - reduction

	return currentPrice
}
