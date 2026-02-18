# TextToJSON

**TextToJSON** is a high-performance Go API designed to transform unstructured text into structured JSON data using the power of **Google Gemini AI**. Built with **Hexagonal Architecture (Ports and Adapters)**, this project demonstrates clean code practices, separation of concerns, and robust error handling.

> **Goal**: Turn chaotic text (tweets, emails, recipes, logs) into clean, predictable JSON with zero regex.

---

## Features

- **AI-Powered Structuring**: Leverages Google Gemini (LLM) to intelligently infer schemas and extract fields like dates, names, locations, and more.
- **Hexagonal Architecture**: Strict separation between core logic, driving adapters (HTTP), and driven adapters (AI providers).
- **Auto-Retry Mechanism**: Intelligent resiliency that retries failed LLM calls up to 3 times to ensure valid JSON output.
- **Swagger Documentation**: Built-in interactive API documentation.
- **12-Factor App**: Configuration via environment variables.

---

## Tech Stack

- **Language**: [Go (Golang)](https://go.dev/)
- **AI Engine**: [Google Gemini](https://ai.google.dev/)
- **Documentation**: [Swagger (Swag)](https://github.com/swaggo/swag)
- **Configuration**: [Godotenv](https://github.com/joho/godotenv)
- **Architecture**: Hexagonal (Ports & Adapters)

---

## Getting Started

### Prerequisites

- Go 1.22+ installed
- A [Google Gemini API Key](https://aistudio.google.com/app/apikey)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/jpp0ca/TextToJSON.git
   cd TextToJSON
   ```

2. **Configure Environment**
   Create a `.env` file in the root directory:
   ```bash
   cp .env.example .env  # if .env.example exists, otherwise just create .env
   ```
   
   Add your API key:
   ```env
   GEMINI_API_KEY=your_google_gemini_api_key_here
   ```

3. **Install Dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the Application**
   ```bash
   go run cmd/api/main.go
   ```

   The server will start at `http://localhost:8080`.

---

## API Documentation (Swagger)

Once the application is running, access the interactive documentation at:

 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

---

## ⚡ Usage Example

Send a `POST` request to `/structure` with raw text.

### Request
**Endpoint**: `POST http://localhost:8080/structure`

**Body**:
```json
{
  "raw_text": "Flight AZ204 from NYC to Rome departs on 2023-10-12 at 18:00 and arrives at 08:00 the next day."
}
```

### Response
The API will inspect the text and return a structured JSON object:

```json
{
  "data": {
    "flight_number": "AZ204",
    "origin": "NYC",
    "destination": "Rome",
    "departure_date": "2023-10-12",
    "departure_time": "18:00",
    "arrival_time": "08:00",
    "arrival_day_offset": 1
  }
}
```

---

## Running Tests

The project includes comprehensive unit tests across all architecture layers (domain, application, and HTTP adapters), using hand-written mocks for interface dependencies.

```bash
go test ./internal/... -v
```

---

## Project Structure

```
.
├── cmd/
│   └── api/            # Application entry point
├── internal/
│   ├── adapters/       # Implementation of interfaces (HTTP, Gemini Client)
│   ├── application/    # Application business logic (Usecases)
│   ├── domain/         # Core business entities (Rules)
│   └── ports/          # Interfaces (Input/Output ports)
├── docs/               # Swagger generated files
├── .env                # Environment variables (gitignored)
└── go.mod              # Go module definition
```

---

## License

This project is licensed under the terms of the LICENSE file included in the repository.
