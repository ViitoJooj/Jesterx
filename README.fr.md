<p align="center">
  <a href="./README.md">🇧🇷 Português</a> | <a href="./README.en.md">🇺🇸 English</a> | <a href="./README.es.md">🇪🇸 Español</a> | 🇫🇷 Français
</p>

# Jesterx

Jesterx est un **SaaS de création de pages**. Il privilégie la simplicité du code et l’apprentissage pratique avec **Golang** au backend et **React + TypeScript** au frontend.

## Stack

- Backend : Golang + Gin, PostgreSQL, MongoDB
- Frontend : React + TypeScript
- Redis et RabbitMQ pour cache et files
- Projet containerisé avec Docker

## Structure

- `backend/` API et règles métier
- `frontend/` interface web
- `sql/` scripts de base de données

## Zone Admin

- Les emails définis dans `ADMIN_EMAILS` reçoivent le rôle `platform_admin`.
- Tableau de bord avec métriques d’utilisateurs, ticket moyen et revenus.
- Gestion des utilisateurs (voir, éditer les données sans exposer les mots de passe, bannir/dé-bannir, supprimer).
- Gestion des plans (prix, descriptions et limites toujours à jour et utilisés par le checkout).
- Export des utilisateurs en XLSX.
- Routes clés : `/v1/plans`, `/v1/admin/plans`, `/v1/admin/users`, `/v1/admin/users/export`, `/v1/admin/stats/overview`.
- Sur le frontend, la page `/admin` affiche les dashboards, l’éditeur de plans et la gestion des comptes pour les admins uniquement.

## Variables d’environnement

Copiez `.env.example` et ajustez les valeurs (BD, clés Stripe, ports). Ajoutez `ADMIN_EMAILS` avec les emails admin séparés par des virgules.

## Docker

```
docker compose up -d
```

Cela démarre le backend, le frontend et les services nécessaires pour un environnement local.
