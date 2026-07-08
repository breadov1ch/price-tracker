# Price Tracker

Price Tracker is a web application for monitoring product prices and receiving updates about changes.

## Features

- User registration and authentication
- Adding products to track
- Viewing the list of tracked products
- Automatic background price updates
- Docker-friendly setup

## Tech Stack

- Backend: Go, Gin, GORM, SQLite
- Frontend: React, Vite, JavaScript
- Containerization: Docker Compose

## Project Structure

- backend/price_tracker — Go-based API
- frontend/price-tracker — web interface
- data/ — application data and SQLite database
- docker-compose.yml — full stack launch with Docker Compose

## Quick Start

### 1. Clone the repository

```bash
git clone <repository-url>
cd price-tracker
```

### 2. Start with Docker Compose

```bash
docker compose up -d --build
```

### 3. Open the application

- Frontend: http://localhost
- Backend API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html

## Environment Variables

For the backend to work correctly, create a file:

```bash
backend/price_tracker/.env
```

Example content:

```env
JWT_SECRET=your_secret_key
GMAIL_APP_PASSWORD=your_gmail_app_password
```

## Useful Commands

```bash
docker compose down
docker compose logs -f backend
docker compose logs -f frontend
```

## License

This project is distributed under the MIT License.
