# Geospatial Country Information API

## Overview

The **Geospatial Country Information API** is a RESTful web service that provides detailed geospatial and demographic data about countries. Users can query information such as population, area size, languages, timezone, currency, and the national flag by providing a country's name.

This project demonstrates modern API development using Golang, leveraging the Echo framework for fast and flexible routing, caching for optimized performance, and Swagger for interactive API documentation.

### Features

- **Country Query**: Retrieve geospatial and demographic data for any country.
- **Caching**: Frequently requested country data is cached for faster responses and reduced API calls.
- **Rate Limiting**: Prevents abuse by limiting the number of requests per client.
- **Error Handling**: Consistent and structured error messages.
- **Interactive Documentation**: Built-in Swagger UI for exploring API endpoints.

### Technologies

- **Golang**: The backend programming language.
- **Echo**: A high-performance web framework for routing and middleware.
- **Swagger**: For API documentation.
- **Ristretto**: In-memory caching for optimized performance.

### API Endpoints

`GET /info`
Retrieve detailed information about a country.

**Query Parameters**:

- `country` (string, required): The name of the country.

**Response Example**:

```json
{
  "data": {
    "country": "Germany",
    "capital": "Berlin",
    "population": 83240525,
    "area_size_km2": 357022,
    "languages": {"deu": "German"},
    "timezone": ["UTC+1"],
    "currency": {"EUR": {"name": "Euro", "symbol": "€"}},
    "flag": "https://flagcdn.com/de.svg"
  }
}
```

**Error Responses**:

- `400 Bad Request`: Missing or invalid parameters.
- `404 Not Found`: Country not found.
- `500 Internal Server Error`: Unexpected error.

## Setup Instructions

### Prerequisites

- Install **Go** (version 1.20 or higher).
- Install **Git** for cloning the repository.

### Installation

1. **Clone the repository**:

    ```bash
    git clone <https://github.com/khamisilawrence/golang-geo-info-api.git>
    cd golang-geo-info-api
    ```

2. **Install dependencies**:

    ```bash
    go mod tidy
    ```

3. **Generate Swagger documentation**:

    ```bash
    swag init --parseDependency --parseInternal
    ```

4. **Start the server**:

    ```bash
    go run main.go
    ```

5. Visit the API:

    - Swagger UI: <http://localhost:8080/swagger/index.html>
    - Health Check: <http://localhost:8080/>

## Testing

### Unit Tests

Run unit tests for individual functions:

```bash
go test ./...
```

## File Structure

```python
geo-info-api/
├── docs/                 # Swagger documentation
├── handlers/             # Request handlers
├── utils/                # Utility functions (caching, error handling, API clients)
├── main.go               # Entry point of the application
├── go.mod                # Module file for dependencies
├── go.sum                # Dependency versions
└── README.md             # Project documentation
```

## License

This project is licensed under the [MIT License](https://opensource.org/license/mit). See the [`LICENSE`](LICENSE) file for details.
