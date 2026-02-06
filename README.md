# Reflecta Backend

A personal reflection and mood tracking API built with Go, Fiber, and MongoDB.

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- MongoDB
- Docker (optional)

### Installation

```bash
# Clone the repository
git clone <repository-url>
cd reflecta-backend

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

### Docker

```bash
docker-compose up --build
```

---

## 📖 API Documentation

**Base URL:** `http://localhost:4000/api`

All requests require `Authorization: Bearer <token>` header (except login/register).

---

## 🔐 Authentication Endpoints

### 1. Register User

| Property | Value |
|----------|-------|
| **Method** | `POST` |
| **URL** | `/api/auth/register` |
| **Auth Required** | No |

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": "user_id",
    "name": "John Doe",
    "email": "user@example.com"
  }
}
```

**Error Responses:**
| Status | Response |
|--------|----------|
| 400 | `{"message": "Invalid input"}` |
| 400 | `{"message": "Email already exists"}` |
| 500 | `{"message": "User creation failed"}` |

---

### 2. Login

| Property | Value |
|----------|-------|
| **Method** | `POST` |
| **URL** | `/api/auth/login` |
| **Auth Required** | No |

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "token": "jwt_token_string",
  "user": {
    "id": "user_id",
    "name": "John Doe",
    "email": "user@example.com"
  }
}
```

**Error Response (400):**
```json
{
  "message": "Invalid credentials"
}
```

---

### 3. Get Profile

| Property | Value |
|----------|-------|
| **Method** | `GET` |
| **URL** | `/api/auth/profile` |
| **Auth Required** | ✅ Yes (Bearer Token) |

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "id": "user_id",
  "name": "John Doe",
  "email": "user@example.com"
}
```

**Error Responses:**
| Status | Response |
|--------|----------|
| 400 | `{"message": "Invalid user ID"}` |
| 404 | `{"message": "User not found"}` |

---

## 📝 Reflection Endpoints

> **Note:** All reflection endpoints require authentication.

### 4. Create Reflection

| Property | Value |
|----------|-------|
| **Method** | `POST` |
| **URL** | `/api/reflections` |
| **Auth Required** | ✅ Yes (Bearer Token) |

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "mood": 3,
  "note": "Had a productive day at work today..."
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `mood` | integer | Yes | Mood value (1-5) |
| `note` | string | Yes | Journal entry text (max 500 chars) |

**Mood Values:**
| Value | Emoji | Label |
|-------|-------|-------|
| 1 | 😔 | Sad |
| 2 | 😕 | Pensive |
| 3 | 😐 | Neutral |
| 4 | 🙂 | Calm |
| 5 | ✨ | Radiant |

**Success Response (200):**
```json
{
  "id": "reflection_id",
  "mood": 3,
  "note": "Had a productive day at work today...",
  "createdAt": "2026-02-06T10:30:00Z",
  "userId": "user_id"
}
```

**Error Responses:**
| Status | Response |
|--------|----------|
| 400 | `{"message": "Invalid input"}` |
| 400 | `{"message": "mood must be between 1 and 5"}` |
| 400 | `{"message": "Note must be 500 characters or less"}` |
| 500 | `{"message": "Failed to save reflection"}` |

---

### 5. Get Weekly Summary

| Property | Value |
|----------|-------|
| **Method** | `GET` |
| **URL** | `/api/reflections/weekly` |
| **Auth Required** | ✅ Yes (Bearer Token) |

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "weeklyData": [
    { "day": "MON", "mood": 2.5 },
    { "day": "TUE", "mood": 3.2 },
    { "day": "WED", "mood": 3.8 },
    { "day": "THU", "mood": 3.5 },
    { "day": "FRI", "mood": 2.8 },
    { "day": "SAT", "mood": 3.0 },
    { "day": "SUN", "mood": 4.5 }
  ],
  "dateRange": "Feb 3 — Feb 9",
  "avgMood": "Positive",
  "topEmotion": "Calm",
  "reflections": "14 Posts",
  "streak": "8 Days",
  "insight": "You felt more creative mid-week, often associated with your morning journaling habit. Notice how your calm state on Sunday correlates with your digital detox."
}
```

| Field | Type | Description |
|-------|------|-------------|
| `weeklyData` | array | Daily mood data for line/area chart |
| `weeklyData[].day` | string | Day abbreviation uppercase (MON, TUE, WED, THU, FRI, SAT, SUN) |
| `weeklyData[].mood` | number | Mood score (1.0 - 5.0 scale, decimals allowed) |
| `dateRange` | string | Display string for week range (e.g., "Feb 3 — Feb 9") |
| `avgMood` | string | Average mood label (e.g., "Positive", "Neutral", "Low") |
| `topEmotion` | string | Most frequent emotion (e.g., "Calm", "Neutral", "Sad") |
| `reflections` | string | Reflection count with label (e.g., "14 Posts") |
| `streak` | string | Current streak with label (e.g., "8 Days") |
| `insight` | string | AI-generated weekly insight/summary text |

**Error Response (500):**
```json
{
  "message": "Failed to load weekly summary"
}
```

---

### 6. Get Insights

| Property | Value |
|----------|-------|
| **Method** | `GET` |
| **URL** | `/api/reflections/insights` |
| **Auth Required** | ✅ Yes (Bearer Token) |

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response (200):**
```json
{
  "moodDistribution": [
    { "day": "Mon", "value": 65, "color": "#6D5D8B" },
    { "day": "Tue", "value": 80, "color": "#6D5D8B" },
    { "day": "Wed", "value": 50, "color": "#6D5D8B" },
    { "day": "Thu", "value": 70, "color": "#6D5D8B" },
    { "day": "Fri", "value": 45, "color": "#C9A24D" },
    { "day": "Sat", "value": 90, "color": "#6D5D8B" },
    { "day": "Sun", "value": 85, "color": "#6D5D8B" }
  ],
  "moodUplift": {
    "value": "+24%",
    "title": "Exercise correlates with higher mood",
    "description": "On days you logged physical activity, your baseline mood was significantly higher than inactive days."
  },
  "aiInsight": "Do these patterns resonate with you today?"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `moodDistribution` | array | Weekly mood values for bar chart |
| `moodDistribution[].day` | string | Day abbreviation (Mon, Tue, Wed, Thu, Fri, Sat, Sun) |
| `moodDistribution[].value` | number | Mood percentage (0-100 scale) |
| `moodDistribution[].color` | string | Hex color for the bar (e.g., "#6D5D8B") |
| `moodUplift` | object | Activity insight card data |
| `moodUplift.value` | string | Stat value (e.g., "+24%", "3.5x") |
| `moodUplift.title` | string | Insight headline |
| `moodUplift.description` | string | Detailed explanation text |
| `aiInsight` | string | AI-generated reflective question |

**Error Response (500):**
```json
{
  "message": "Failed to load insights"
}
```

---

## 📊 API Summary

| Endpoint | Method | URL | Auth | Description |
|----------|--------|-----|------|-------------|
| Register | POST | `/api/auth/register` | ❌ | User registration |
| Login | POST | `/api/auth/login` | ❌ | User authentication |
| Get Profile | GET | `/api/auth/profile` | ✅ | Get user profile |
| Create Reflection | POST | `/api/reflections` | ✅ | Save reflection entry |
| Weekly Summary | GET | `/api/reflections/weekly` | ✅ | Load weekly summary |
| Insights | GET | `/api/reflections/insights` | ✅ | Load mood insights |

---

## 🔑 Authentication

All protected endpoints require a JWT token in the `Authorization` header:

```
Authorization: Bearer <your_jwt_token>
```

The token is obtained from the `/api/auth/login` or `/api/auth/register` endpoints.

---

## ⚠️ Error Handling

All error responses follow this format:

```json
{
  "message": "Human-readable error description"
}
```

---

## 📦 Data Models

### User
```json
{
  "id": "ObjectID",
  "name": "string",
  "email": "string",
  "created_at": "int64 (timestamp)"
}
```

### Reflection
```json
{
  "id": "ObjectID",
  "user_id": "ObjectID",
  "mood": "int (1-5)",
  "note": "string (max 500 chars)",
  "date": "datetime",
  "created_at": "datetime"
}
```

---

## 📝 Notes

- Token is stored in `expo-secure-store` under key `"token"` (mobile app)
- Token is automatically attached to all requests via axios interceptor
- All timestamps are in ISO 8601 format (UTC)
- Mood values use 1-5 scale internally, converted to percentages (0-100) for charts where needed

---

## 📁 Project Structure

```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   ├── controllers/
│   ├── database/
│   ├── middleware/
│   ├── models/
│   ├── routes/
│   └── utils/
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```
