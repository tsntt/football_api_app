# API Documentation

This document provides a detailed description of the API endpoints.

## Authentication

### `POST /auth/register`

Registers a new user.

**Request Body:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "message": "User created successfully"
}
```

### `POST /auth/login`

Logs in a user.

**Request Body:**

```json
{
  "username": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "token": "string"
}
```

### `POST /auth/logout`

Logs out a user.

**Response:**

```json
{
  "message": "Logout successful"
}
```

---

## Championships

### `GET /championship`

Retrieves a list of all championships.

**Response:**

```json
[
  {
    "id": "integer",
    "name": "string",
    "season": "string"
  }
]
```

### `GET /championship/:id/matches`

Retrieves a list of matches for a specific championship.

**Query Parameters:**

- `team` (optional): Filters matches by team name.
- `stage` (optional): Filters matches by stage (e.g., "group-stage", "quarter-finals").

**Response:**

```json
[
  {
    "id": "integer",
    "championship_id": "integer",
    "team_a": "string",
    "team_b": "string",
    "match_date": "string",
    "stage": "string"
  }
]
```

---

## Fans

### `POST /fans`

Subscribes a fan to a championship. This endpoint requires authentication.

**Request Body:**

```json
{
  "championship_id": "integer"
}
```

**Response:**

```json
{
  "message": "Successfully subscribed to championship"
}
```

---

## Admin

These endpoints require administrator privileges.

### `GET /admin`

Retrieves a list of all matches.

**Response:**

```json
[
  {
    "id": "integer",
    "championship_id": "integer",
    "team_a": "string",
    "team_b": "string",
    "match_date": "string",
    "stage": "string"
  }
]
```

### `POST /admin/broadcast/:match_id`

Broadcasts a match.

**Response:**

```json
{
  "message": "Match broadcast initiated"
}
```
