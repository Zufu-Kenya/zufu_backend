package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return device, err
	}

	device.Brand = doc.Find("div.brand").Text()
	device.Model = doc.Find("div.model").Text()
	device.ReleaseDate = doc.Find("div.release-date").Text()

	priceString := doc.Find("div.price").Text()
	price, err := strconv.ParseFloat(priceString, 64)
	if err != nil {
		return device, err
	}
	device.Price = price

	// Add scraping for brand-specific information
	if strings.Contains(url, "samsung") {
		device.Brand = "Samsung"
		device.Model = doc.Find("h1.product-details").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div > div:nth-child(1) > div > span.value").Text()
	} else if strings.Contains(url, "haier") {
		device.Brand = "Haier"
		device.Model = doc.Find("div.product-title > h1").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(2) > div > span.value").Text()
	} else if strings.Contains(url, "lg") {
		device.Brand = "LG"
		device.Model = doc.Find("h1.product-name").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(1) > div > span.value").Text()
	} else if strings.Contains(url, "mika") {
		device.Brand = "Mika"
		device.Model = doc.Find("div.product-title > h1").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(2) > div > span.value").Text()
	} else if strings.Contains(url, "hotpoint") {
		device.Brand = "Hotpoint"
		device.Model = doc.Find("div.product-title > h1").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(2) > div > span.value").Text()
	} else if strings.Contains(url, "von") {
		device.Brand = "Von"
		device.Model = doc.Find("div.product-title > h1").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(2) > div > span.value").Text()
	} else if strings.Contains(url, "bruhm") {
		device.Brand = "Bruhm"
		device.Model = doc.Find("div.product-title > h1").Text()
		device.ReleaseDate = doc.Find("div.product-essential > div:nth-child(2) > div > span.value").Text()
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
	percent := 0.05 * float64(age) // 5% reduction in price for every year
	reduction := percent * originalPrice
	currentPrice := originalPrice - reduction

	return currentPrice
}
