# Reflecta Backend

Reflecta is a personal reflection and mood-tracking API. Users register an account, log short daily reflections with a mood score, and get back weekly summaries and insights about their emotional patterns (average mood, top emotion, journaling streak, and generated insight text).

Built with **Go**, **Fiber**, and **MongoDB**.

## Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.25 |
| Web framework | [Fiber v2](https://gofiber.io/) |
| Database | MongoDB (official Go driver) |
| Auth | JWT (`golang-jwt/jwt/v5`) + bcrypt password hashing |
| Config | `.env` file loaded via `godotenv` |
| Deployment | Docker + Docker Compose |

## Features

- **User accounts** — register and log in with email/password; passwords are bcrypt-hashed and a JWT is returned on success.
- **Reflections** — authenticated users create journal entries with a mood score (1–5) and a note (max 500 characters).
- **Weekly summary** — per-day average mood for the current week (Mon–Sun), average mood label, most frequent emotion, post count, daily journaling streak, and a generated insight message.
- **Personal insights** — mood distribution by day of week (as 0–100 chart values with colors), a "mood uplift" insight card, and a reflective prompt question.
- **Rate limiting** — global limit of 100 requests per minute per client.
- **Request logging** — Fiber's logger middleware on all routes.

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go              # Entry point: env, DB, middleware, routes
├── internal/
│   ├── config/                  # Environment variable loading
│   ├── controllers/             # Auth and reflection handlers
│   ├── database/                # MongoDB connection
│   ├── middleware/              # JWT auth middleware
│   ├── models/                  # User and Reflection models
│   ├── routes/                  # Route definitions
│   └── utils/                   # JWT helpers, validation, mutex lock
├── docker-compose.yml           # API + MongoDB services
├── Dockerfile
└── go.mod
```

## Getting Started

### Prerequisites

- Go 1.25+
- MongoDB (local instance, or use Docker Compose below)

### Configuration

Copy `.env.example` to `.env` and fill in your values:

```env
PORT=4000
MONGO_URI=mongodb://localhost:27017
JWT_SECRET=your-secret-key
```

### Run locally

```bash
go mod download
go run cmd/server/main.go
```

The server starts on `http://localhost:4000`.

### Run with Docker

```bash
docker-compose up --build
```

This starts the API on port `4000` and a MongoDB 7 container on port `27017` with a persistent volume.

## API Reference

Base URL: `http://localhost:4000/api`

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/auth/register` | No | Create an account, returns JWT + user |
| POST | `/api/auth/login` | No | Log in, returns JWT + user |
| GET | `/api/auth/profile` | Yes | Get the current user's profile |
| GET | `/api/user/me` | Yes | Get the current user's ID from the token |
| POST | `/api/reflections` | Yes | Create a reflection entry |
| GET | `/api/reflections/weekly` | Yes | Weekly mood summary |
| GET | `/api/reflections/insights` | Yes | Mood distribution and insights |

Protected endpoints require the header:

```
Authorization: Bearer <token>
```

Errors are returned as `{"message": "..."}` with an appropriate HTTP status code.

### Register / Login

```json
POST /api/auth/register
{ "name": "Jane Doe", "email": "jane@example.com", "password": "secret123" }

POST /api/auth/login
{ "email": "jane@example.com", "password": "secret123" }
```

Both respond with:

```json
{
  "token": "<jwt>",
  "user": { "id": "...", "name": "Jane Doe", "email": "jane@example.com" }
}
```

### Create Reflection

```json
POST /api/reflections
{ "mood": 4, "note": "Had a productive day." }
```

- `mood` — integer 1–5 (1 Sad, 2 Pensive, 3 Neutral, 4 Calm, 5 Radiant)
- `note` — string, max 500 characters

### Weekly Summary

`GET /api/reflections/weekly` returns per-day averages for the current Monday–Sunday week:

```json
{
  "weeklyData": [{ "day": "MON", "mood": 3.5 }, ...],
  "dateRange": "Jul 14 — Jul 20",
  "avgMood": "Positive",
  "topEmotion": "Calm",
  "reflections": "5 Posts",
  "streak": "3 Days",
  "insight": "You had a great week! ..."
}
```

### Insights

`GET /api/reflections/insights` returns chart-ready data for the current week:

```json
{
  "moodDistribution": [{ "day": "Mon", "value": 70, "color": "#6D5D8B" }, ...],
  "moodUplift": {
    "value": "+24%",
    "title": "Positive trend detected",
    "description": "..."
  },
  "aiInsight": "What patterns do you see in your reflections?"
}
```

`value` is the day's average mood converted to a 0–100 scale; low-mood days (≤ 2.5 average) are colored amber (`#C9A24D`) instead of the default purple.

## Data Models

**User** (`users` collection)

| Field | Type | Notes |
|-------|------|-------|
| `_id` | ObjectID | |
| `name` | string | |
| `email` | string | Must be unique (checked at registration) |
| `password` | string | bcrypt hash, never returned in JSON |
| `created_at` | int64 | Unix timestamp |

**Reflection** (`reflections` collection)

| Field | Type | Notes |
|-------|------|-------|
| `_id` | ObjectID | |
| `user_id` | ObjectID | Owner of the entry |
| `mood` | int | 1–5 |
| `note` | string | Max 500 characters |
| `date` | datetime | Set at creation |
| `created_at` | datetime | Used for weekly grouping and streaks |
