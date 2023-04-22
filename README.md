# zufu_backend

# Device Information API
This is an ongoing project to build an API that retrieves device information from a website and calculates the current price of the device based on its release date and age. The API is built using the Go programming language.

# Usage
To use the API, make a GET request to the /device endpoint. The API will retrieve device information from a hardcoded website and calculate the current price based on the device's release date and age. The response will be in JSON format and include the device brand, model, release date, and current price.

# Query Parameters
You can optionally include a status query parameter to specify the status of the device. If present, the status value will be included in the response.

### Example Request

GET /device?status=available HTTP/1.1
Host: localhost:8080

### Example Response

HTTP/1.1 200 OK
Content-Type: application/json

{
    "brand": "Example Brand",
    "model": "Example Model",
    "release_date": "2022-03-15",
    "price": 950.0,
    "status": "AVAILABLE"
}

# Contributing
This is an open project and contributions are welcome. However, please note that Zufu Kenya retains the rights of ownership and admission. If you would like to contribute, please fork the project and create a new branch for your changes. When your changes are ready, submit a pull request and they will be reviewed and merged into the main branch if approved.
