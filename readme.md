# Go Spinsolf Test

Simple backend service เขียนด้วย Go + Fiber + MongoDB
ใช้สำหรับ import และ query ข้อมูลสถานีรถไฟ พร้อมฟังก์ชันค้นหาสถานีใกล้เคียง

---

```
## Project Structure
├── app/              # Application setup (Fiber, routes)
├── config/           # Config & Database connection
├── controllers/      # HTTP handlers
├── dto/              # Response DTOs
├── middleware/       # Middlewares (API Key, CORS, Logger)
├── models/           # Database models
├── parsers/          # File parsers (.csv, .json, .xlsx)
├── services/         # Business logic
├── utils/            # Helper functions (haversine, normalizer, mapper)
└── main.go           # Entry point
```

---

## Set up

1. Clone Project
  - `git clone https://github.com/Teneieiza/go-spinsolf-test.git`
  - `cd go-spinsolf-test`

2. สร้างไฟล์ .env

3. Run server
  - ใช้คำสั่ง `air` ใช้ hot reload `github.com/air-verse/`

---

## API Endpoints

- Health Check
    `GET /api/health`

- Stations
  - Import ผ่าน URL
    `POST /api/stations/import/url`
      Exam: `/api/stations/import/url?url=https://example.com/stations.json`

  - Import ผ่านไฟล์ (.csv, .json, .xlsx)
    `POST /api/stations/import/file`
      Exam: `/api/stations/import/file (form-data: file=...)`

  - Nearlest Station
    `GET /api/stations/nearby`
      Exam: `/api/stations/nearby?lat=13.75&long=100.50&limit=5`

  - Nearlest Station with pagination
    `GET /api/stations/nearby/paginated`
      Exam: `/api/stations/nearby/paginated?lat=13.75&long=100.50&page=1&page_size=10`

---

## API Key

  - ส่งใน header:
  `x-api-key: your_api_key`

  - query param
  `?api_key=your_api_key`

---

## Tech Stack

  - [go](https://go.dev/) + [Fiber](https://gofiber.io/)  `(Web framework)`
  - [MongoDb](https://www.mongodb.com/)                   `(Database)`
  - [Excelize](https://github.com/qax-os/excelize)        `(XLSX parser)`

---
