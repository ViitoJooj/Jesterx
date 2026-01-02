<p align="center">
  <a href="./README.md">🇧🇷 Português</a> | 🇺🇸 English | <a href="./README.es.md">🇪🇸 Español</a> | <a href="./README.fr.md">🇫🇷 Français</a>
</p>

# Jesterx

Jesterx is a **SaaS page builder**. The initial version (v1) will be open and public on GitHub, serving both as a functional product and as a technical experiment. The project focuses on **code simplicity**, hands-on learning, and idea validation. It exists to test my skills with **Golang**, study scalable architecture, and, if possible, generate revenue in the future.

### Stack and Technical Decisions

- **Golang + Gin** on the backend
- **React + TypeScript** on the frontend
- **PostgreSQL** for users
- **MongoDB** for website persistence
- **Redis** for caching and scalability
- **RabbitMQ** for notifications and asynchronous tasks
- Authentication via **OAuth2**
- Public APIs
- Everything **containerized with Docker**
- Infrastructure on **Azure**, with **CI/CD**
- **Automated tests in Python**, chosen for their simplicity and speed

---

## Overview

This repository brings together everything you need to start an e-commerce site:

- API for business rules
- Web interface
- Database structure

The idea is to allow any developer to clone the project, set up the environment, and start working without much initial configuration.

---

## Project structure

The project follows a **modular** architecture, divided into three main parts:

```
jesterx/
       ├─ backend/           # API and server logic
       ├─ frontend/          # User interface
       ├─ sql/               # Database scripts
       ├─ .env.example       # Example of environment variables
       ├─ docker-compose.yml
       ├─ LICENSE.md
```

---

## Backend

The backend concentrates all the application logic, such as:

- User authentication
- Products
- Orders
- Communication with the database

The API follows the REST standard, with the possibility of future adaptation if necessary.

---

## Frontend

The frontend is responsible for the store's interface, including:

- Product listing
- Shopping cart
- Login and registration
- Checkout
- Admin area (dashboard for platform admins)

It directly consumes the backend API.

---

## Admin area

- Only emails listed in `ADMIN_EMAILS` receive the `platform_admin` role.
- Dashboard with user growth, average ticket, revenue, and plan usage metrics.
- User management (view, edit profile data and plan without exposing passwords, ban/unban, delete).
- Plan management (update prices, descriptions, limits, and benefits). Checkout always uses the latest values.
- XLSX export of users directly from the interface or API.

### Key new routes

- `GET /v1/plans` – public list of plans and limits.
- `GET /v1/admin/plans` and `PUT /v1/admin/plans/:plan_id` – query and update prices/descriptions/limits.
- `GET /v1/admin/users`, `PUT /v1/admin/users/:user_id`, `PUT /v1/admin/users/:user_id/ban`, `DELETE /v1/admin/users/:user_id` – user management.
- `GET /v1/admin/users/export` – XLSX export.
- `GET /v1/admin/stats/overview` – data for admin dashboards.

On the frontend, the `/admin` route is protected by role and shows dashboards, plan editor, and account management.

---

## Database

The database scripts are located in the `sql/` folder, including:

- Table creation
- Relationships
- Initial data (when applicable)

---

## Environment Configuration

### Environment Variables

Copy the example file:

```bash
cp .env.example .env
```

Then adjust the variables according to your environment, such as database, ports, and access keys.
Add the list of admin emails in `ADMIN_EMAILS` (comma-separated) to unlock the `/admin` dashboard and `/v1/admin` routes.

## Docker

The project has a `docker-compose.yml` file to facilitate local setup:

```
docker compose up -d
```

This will bring up the backend, frontend, Redis, RabbitMQ, and database.

## Features

- User registration and authentication
- Product CRUD (Create, Read, Delete)
- Shopping cart
- Order system
- Checkout
- Administrative panel
- Payment processing integrations (future)

## Local Development

For local development:

```
# Backend cd backend
# Install dependencies and run the server

# Frontend cd frontend
# Install dependencies and run the app
```

## Contribution

#### Feel free to contribute:

1. Fork the repository
2. Create a branch (feature/my-feature)
3. Commit your changes
4. Open a Pull Request

If in doubt, check the <a href="./CONTRIBUTING.md">CONTRIBUTING.md</a>

## License

This project is licensed under the terms of the file <a href="LICENSE.md">LICENSE.md</a>

## Author

Developed by ViitoJooj (819SauCe)
