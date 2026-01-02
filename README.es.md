<p align="center">
  <a href="./README.md">🇧🇷 Português</a> | <a href="./README.en.md">🇺🇸 English</a> | 🇪🇸 Español | <a href="./README.fr.md">🇫🇷 Français</a>
</p>

# Jesterx

Jesterx es un **SaaS para crear páginas**. El proyecto prioriza la simplicidad del código y el aprendizaje práctico con **Golang** en el backend y **React + TypeScript** en el frontend.

## Pila

- Backend: Golang + Gin, PostgreSQL, MongoDB
- Frontend: React + TypeScript
- Redis y RabbitMQ para cache y colas
- Todo containerizado con Docker

## Estructura

- `backend/` API y lógica de negocio
- `frontend/` interfaz web
- `sql/` scripts de base de datos

## Área de administración

- Los correos definidos en `ADMIN_EMAILS` reciben la función `platform_admin`.
- Dashboard con métricas de usuarios, ticket medio y facturación.
- Gestión de usuarios (ver, editar datos sin mostrar contraseñas, banear/desbanear, eliminar).
- Gestión de planes (precios, descripciones y límites actualizados, usados por el checkout).
- Exportación de usuarios en XLSX.
- Rutas clave: `/v1/plans`, `/v1/admin/plans`, `/v1/admin/users`, `/v1/admin/users/export`, `/v1/admin/stats/overview`.
- En el frontend, la ruta `/admin` muestra dashboards, edición de planes y gestión de cuentas solo para admins.

## Variables de entorno

Copie `.env.example` y configure según su entorno (DB, claves de Stripe, puertos). Incluya `ADMIN_EMAILS` con los correos de administradores separados por comas.

## Docker

```
docker compose up -d
```

Esto levanta backend, frontend y servicios de soporte para un entorno local.
