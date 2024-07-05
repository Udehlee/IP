
## IP - Stage One Task Backend-Track


## Overview

This project is a basic web server that retrieves client's IP address, location, and current weather condition( temperature).

## Features

- Obtains the client's IP address.
- Identifies the client's city using their IP address.
- Retrieves the current temperature for the  city.
- Returns a customized greeting that includes the visitor's name and the temperature details.

## API Endpoint

```sh
GET /api/hello?visitor_name=Mark

```
## Response

```json
>>>>>>> origin/main
{
    "client_ip": "127.0.0.1",
    "location": "New York",
    "greeting": "Hello, Mark!, the temperature is 11 degrees Celsius in New York"
} 
```
## Prerequisites
Before you begin, ensure you have Go installed 

## Installation

- Clone this repository to your local machine:

 ```sh
git clone https://github.com/Udehlee/IP.git

```
- Navigate to the project directory:

-  cd IP

- Run the server:

```sh
go run main.go

 ```
- The server will be listening on port 8080.

- Access the server on Browser:

```sh
http:localhost:8080

```
