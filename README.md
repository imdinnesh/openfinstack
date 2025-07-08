# Open Finstack 🚀

Welcome to Open Finstack, a monorepo designed for building robust financial technology services. This repository is structured to promote reusability, maintainability, and efficient deployment of various microservices.

---

## Repository Structure 🏗️

Open Finstack adopts a clear and modular structure to streamline development and deployment:

* **`deploy/`**: This folder contains all the **Dockerfiles** 🐳 and related deployment configurations. You'll find everything needed to containerize and deploy your services here.

* **`pkg/`**: This directory houses **shared resources** 🤝 and packages. Any common utilities, libraries, or interfaces used across multiple services are centralized here to avoid duplication and ensure consistency.

* **`services/`**: This is where all the individual **microservices** 🧩 reside. Each service is self-contained within its own subfolder, promoting clear separation of concerns.

---

## Services 💡

Open Finstack currently includes the following core services:

### 1. Auth Management 🔑

* **Purpose**: Handles user authentication, authorization, and session management.
* **Key Features**: User registration, login, token generation/validation, password management, and role-based access control.
* **Location**: `services/auth/`
* **Further Details**: Refer to `services/auth/readme.md` for specific API documentation, setup instructions, and usage examples.

### 2. KYC/Profile Management 👤

* **Purpose**: Manages Know Your Customer (KYC) processes and user profile information.
* **Key Features**: User profile creation and updates, document uploads for verification, status tracking of KYC applications, and data retrieval for compliance.
* **Location**: `services/kyc/`
* **Further Details**: Refer to `services/kyc/readme.md` for specific API documentation, setup instructions, and usage examples.