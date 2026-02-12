<p align="center">
  <a href="../README.md">🇧🇷 Português</a> | 🇺🇸 English | <a href="./README_cn.md">🇨🇳 CN</a>
</p>

# Jesterx

Jesterx is a **page creation SaaS**.  
The initial version (v1) will be open and public on GitHub, serving both as a functional product and as a technical experiment. The project focuses on **code simplicity**, hands-on learning, and idea validation.  
It exists to test my skills with **Golang**, study scalable architecture, and potentially generate revenue in the future.

### Stack and technical decisions
- **Golang + Gin** on the backend
- **React + TypeScript** on the frontend
- **PostgreSQL** for users
- **MongoDB** for site persistence
- **Redis** for caching and scalability
- **RabbitMQ** for notifications and asynchronous tasks
- Authentication via **OAuth2**
- Public APIs
- Everything **containerized with Docker**
- Infrastructure on **Azure**, with **CI/CD**
- **Automated tests in Python**, chosen for simplicity and speed

---

## Overview
This repository brings together everything needed to start an e-commerce:
- API for business rules
- Web interface
- Database structure

The idea is to allow any developer to clone the project, spin up the environment, and start working without much initial configuration.

---

## Backend
The backend concentrates all application logic, such as:
- User authentication
- Products
- Orders
- Database communication

The API follows the REST pattern, with possibility for future adaptation if needed.

---

## Frontend
The frontend is responsible for the store interface, including:
- Product listing
- Shopping cart
- Login and registration
- Checkout
- Administrative area (dashboard for admins)

It directly consumes the backend API.

---

## Administrative Area
- Only users with email listed in `ADMIN_EMAILS` receive the `platform_admin` role.
- Dashboard with metrics on users created, average ticket, revenue, and most used plans.
- User management (view, edit profile and plan data without exposing password, ban/unban, delete).
- Plan management (change prices, descriptions, limits, and benefits). Checkout always uses updated values.
- User export in XLSX directly from the interface and via API.

---

## Database
In the `migrations/` folder are the database scripts, including:
- Table creation
- Relationships
- Initial data (when applicable)

---

## Environment setup

### Environment variables
Copy the example file:
```bash
cp .env.example .env
```
Then adjust the variables according to your environment, such as database, ports, and access keys.

Include the list of administrative emails in `ADMIN_EMAILS` (separated by commas) to enable the `/admin` dashboard and `/v1/admin` routes.

## Docker
The project has a _docker-compose.yml_ to facilitate local setup:
```bash
docker compose up -d
```
This will spin up the backend, frontend, redis, rabbitMQ, and database.

## Features
- User registration and authentication
- (post/get) REST communication in software services
- Product CRUD
- Shopping cart
- Order system
- Checkout
- Administrative panel
- Payment gateway integrations

## Local development
For local development:
install dependencies
```bash
go mod tidy
```

## Contributing
#### Want to contribute?
Just:
1. Fork it
2. Create a branch (author/your-feature)
3. Commit your changes
4. Open a Pull Request

- For any questions, check out <a href="./CONTRIBUTING.md">CONTRIBUTING.md</a>

## License
This project is licensed according to the <a href="LICENSE.md">LICENSE.md</a> file