package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

type WeatherData struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`

	Current struct {
		Tempc float64 `json:"temp_c"`
	} `json:"current"`
}

// IPAddr gets the ip address of the client
func IPAddr(r *http.Request) (string, error) {

	xForwarded := r.Header.Get("x-forwarded-for")
	if xForwarded != "" {
		ips := strings.Split(xForwarded, ",")
		return strings.TrimSpace(ips[0]), nil // Return the first IP without error

	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("error getting IP Addr: %v", err)
	}

	return ip, nil
}

// WeatherInfo gets weather information based on the IP address
func WeatherInfo(location string) (WeatherData, error) {

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, location)

	var weatherData WeatherData

	resp, err := http.Get(url)
	if err != nil {
		return weatherData, fmt.Errorf("unable to fetch weather details: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return weatherData, fmt.Errorf("unable to read response body: %v", err)
	}

	if err := json.Unmarshal(body, &weatherData); err != nil {
		return weatherData, fmt.Errorf("unable to unmarshal JSON: %v\nResponse body: %s", err, body)
	}

	return weatherData, nil
}

func hello(w http.ResponseWriter, r *http.Request) {

	visitor := r.URL.Query().Get("visitor_name")

	visitor = strings.TrimSpace(visitor)
	if visitor == "" {
		visitor = "Guest"
	}

	IP, err := IPAddr(r)
	if err != nil {
		log.Fatal(err)
	}

	weatherData, err := WeatherInfo(IP)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to fetch weather information: %v", err), http.StatusInternalServerError)
		return
	}

	locationName := weatherData.Location.Name
	tempCelsius := weatherData.Current.Tempc

	greeting := fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", visitor, tempCelsius, locationName)
	response := Response{
		ClientIP: IP,
		Location: locationName,
		Greeting: greeting,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", index)
	mux.HandleFunc("/api/hello", hello)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)

}
