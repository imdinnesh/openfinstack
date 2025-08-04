# KYC Service - OpenFinStack

The **KYC (Know Your Customer)** service is a dedicated microservice within the **OpenFinStack** fintech platform. It manages user verification workflows, including document submission, status tracking, and integration with different verification strategies (e.g., mock/manual/external).

---

## 🚀 Features

- User KYC submission and retrieval  
- Admin-based KYC verification flow  
- Support for multiple verification types (mock, manual, pluggable)  
- Kafka-based event publishing on status updates  
- Modular service, repository, and handler layers  
- Clean architecture with clear separation of concerns

---

## 🧱 Tech Stack

- Go (Gin + GORM)  
- PostgreSQL  
- Redis (for session/state management, optional)  
- Kafka (for event-driven communication)  
- gRPC (optional, for inter-service communication)

---

## 📁 Folder Structure

```
services/
└── kyc/
    ├── cmd/                 # Application entry point
    ├── config/              # Environment configs and constants
    ├── handler/             # HTTP request handlers
    ├── models/              # GORM models
    ├── repository/          # Database logic
    ├── service/             # Business logic
    ├── verifier/            # Verifier plugins (manual, mock)
    ├── routes/              # Route registration
    ├── proto/               # gRPC definitions (optional)
    └── migrations/          # GORM migrations
```

---

## ⚙️ Setup Instructions

### 1. Environment Variables

Create a `.env` file in `services/kyc/`:

```
PORT=8081
DATABASE_URL=postgres://user:password@localhost:5432/openfinstack
REDIS_URL=localhost:6379
KYC_VERIFIER=mock
```

### 2. Run Locally

```bash
cd services/kyc
go run cmd/main.go
```

---

## 📌 API Endpoints

### Submit KYC

```
POST /api/v1/kyc
```

**Request Body:**

```json
{
  "user_id": 1,
  "document_type": "aadhaar",
  "document_number": "1234-5678-9012"
}
```

---

### Get KYC by User ID

```
GET /api/v1/kyc/user/:user_id
```

---

### Admin Verify KYC

```
POST /api/v1/kyc-admin/verify/:kyc_id
```

---

## 🔁 Verifier Plug-in System

The KYC module supports multiple verification backends:

- `mock`: Always returns successful verification (used for dev/testing)  
- `manual`: Admin manually verifies each KYC request  

You can add new verification backends by implementing:

```go
type Verifier interface {
    Verify(kyc *models.KYC) (*models.VerificationResult, error)
}
```

And registering it in `verifier/factory.go`.

---

## 🪄 Migrations

We use `gormigrate` for safe database schema evolution. Migration files live in `services/kyc/migrations`.

### Example Usage:

```go
m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
    migrations.AddKYCFields20250803(),
    migrations.AddKYCStatusLogs20250804(),
})
_ = m.Migrate()
```

---

## 📣 Kafka Integration

When a KYC is verified, an event is published to Kafka:

```go
kafka.Publish("kyc.status.updated", KYCEventPayload)
```

This allows other services to listen for updates (e.g., notification service).

---

## ✅ TODO

- Add support for document image uploads  
- Retry mechanism for failed KYCs  
- Add audit trail with timestamped logs  
- Implement gRPC interface for inter-service use  
- Role-based access controls for admin endpoints

---
