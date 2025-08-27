# Open Finstack 🚀  

Welcome to **Open Finstack**, a **monorepo** built for powering robust financial technology services.  
This repository is designed with **reusability**, **maintainability**, and **scalability** in mind, making it easier to build, extend, and deploy fintech microservices.  

---

## 📂 Repository Structure  

Open Finstack follows a **modular architecture** for clarity and efficiency:  

- **`deploy/`** 🐳  
  Contains all **Dockerfiles** and deployment configurations for containerizing and orchestrating services.  

- **`pkg/`** 📦  
  Shared libraries and resources used across services (e.g., logging, Kafka, Redis).  
  Helps maintain **consistency** and **avoids duplication**.  

- **`services/`** 🧩  
  Houses all microservices, each within its own subfolder for **clear separation of concerns**.  

---

## ⚡ Core Services  

### 1. Auth Service 🔑  
- **Purpose:** Manages user **authentication**, **authorization**, and **session management**.  
- **Features:**  
  - User registration & login  
  - JWT token generation & validation  
  - Password management  
  - Role-Based Access Control (RBAC)  
- **Location:** `services/auth/`  
- **More Info:** See [services/auth/readme.md](services/auth/readme.md)  

---

### 2. KYC / Profile Service 👤  
- **Purpose:** Handles **Know Your Customer (KYC)** workflows and user profile management.  
- **Features:**  
  - Profile creation & updates  
  - Document uploads & verification  
  - KYC status tracking  
  - Compliance data retrieval  
- **Location:** `services/kyc/`  
- **More Info:** See [services/kyc/readme.md](services/kyc/readme.md)  

---

### 3. Wallet Service 💰  
- **Purpose:** Manages **user wallets**, deposits, withdrawals, and peer-to-peer transfers.  
- **Features:**  
  - Wallet creation  
  - Balance tracking  
  - Atomic transfers  
  - Transaction history  
  - Publishes wallet events to Kafka for downstream services  
- **Location:** `services/wallet/`  
- **More Info:** See [services/wallet/readme.md](services/wallet/readme.md)  

---

### 4. Ledger Service 📚  
- **Purpose:** Provides an **immutable double-entry accounting system** ensuring financial integrity.  
- **Features:**  
  - Transaction recording & retrieval  
  - Reversals for error handling  
  - Consumes wallet events to post ledger entries automatically  
  - Serves as the **ultimate audit log** for the platform  
- **Location:** `services/ledger/`  
- **More Info:** See [services/ledger/readme.md](services/ledger/readme.md)  

---

### 5. Gateway Service 🚪  
- **Purpose:** Acts as the **front door** to Open Finstack, routing requests to the right microservice.  
- **Features:**  
  - Centralized authentication  
  - API rate limiting  
  - Dynamic routing via YAML config  
  - Built-in observability & metrics  
- **Location:** `gateway/`  
- **More Info:** See [gateway/readme.md](gateway/readme.md)  

---

## 🏦 Why Open Finstack?  

- ✅ **Event-driven & modular** design  
- ✅ **Auditable & reliable** financial tracking via Ledger  
- ✅ **Scalable microservice** architecture  
- ✅ **Secure & compliant** by design  

---

