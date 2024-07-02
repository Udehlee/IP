package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Response struct {
	ClientIP string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

// IPAddr fetch the IPAddress of the requester
func IPAddr() (string, error) {

	ipifyURL := "https://api.ipify.org?format=json"

	resp, err := http.Get(ipifyURL)
	if err != nil {
		return "", fmt.Errorf("unable to fetch IP address: %w", err)
	}
	defer resp.Body.Close()

	ipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read IP data: %w", err)
	}

	ip := struct {
		IP string `json:"ip"`
	}{}

	if err := json.Unmarshal(ipData, &ip); err != nil {
		return "", fmt.Errorf("failed to unmarshal IP data: %w", err)
	}

	return ip.IP, nil
}

// Location gets the current city
func Location(ip string) (string, error) {

	apiKey := os.Getenv("IPGEOLOCATION_API_KEY")

	url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", apiKey, ip)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get location: %s", resp.Status)
	}

	city, ok := result["city"].(string)
	if !ok {
		return "", fmt.Errorf("failed to get city")
	}

	return city, nil
}

func Temp(city string) (string, error) {

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	weatherData := struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f°C", weatherData.Main.Temp), nil
}

func hello(w http.ResponseWriter, r *http.Request) {

	visitorName := r.URL.Query().Get("visitor_name")

	clientIP, err := IPAddr()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	location, err := Location(clientIP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	temp, err := Temp(location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		ClientIP: clientIP,
		Location: location,
		Greeting: fmt.Sprintf("Hello, %s! The temperature is %s in %s", visitorName, temp, location),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	http.HandleFunc("/api/hello", hello)
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
