# API Documentation

This document provides a detailed description of the API endpoints.

[Postman Collection](https://.postman.co/workspace/Personal-Workspace~54934cc3-4386-4d24-ad9c-76441e3e236d/collection/1936338-ae95e92a-beb7-4222-818a-f5a1a6edce12?action=share&creator=1936338&active-environment=1936338-a44ab973-0ad0-49e9-b85f-aad4fa919495)

## Authentication

### `POST api/v1/auth/register`

Registers a new user.

**Request Body:**

```json
{
  "name": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "message": "User created successfully"
}
```

### `POST api/v1/auth/login`

Logs in a user.

**Request Body:**

```json
{
  "name": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "token": "string"
}
```

### `POST api/v1/auth/logout`

Logs out a user.

**Headers:**
Authorization: Bearer YOUR_JWT_TOKEN_HERE

**Response:**

```json
{
  "message": "Logout successful"
}
```

---

## Championships

### `GET api/v1/championship`

Retrieves a list of all championships.

**Headers:**
Authorization: Bearer YOUR_JWT_TOKEN_HERE

**Response:**

```json
[
  {
    "id": 2013,
    "name": "Campeonato Brasileiro Série A",
    "code": "BSA",
    "type": "LEAGUE",
    "emblem": "https://crests.football-data.org/bsa.png",
    "currentSeason": {
      "id": 2371,
      "startDate": "2025-03-29",
      "endDate": "2025-12-21",
      "currentMatchday": 25,
      "winner": null
    },
    "seasons": null
  }
]
```

### `GET api/v1/championship/:id/matches`

Retrieves a list of matches for a specific championship.

**Headers:**
Authorization: Bearer YOUR_JWT_TOKEN_HERE

**Query Parameters:**

- `team` (optional): Filters matches by team name.
- `stage` (optional): Filters matches by stage (e.g., "group-stage", "quarter-finals").

**Response:**

```json
[
  {
        "id": 534938,
        "utcDate": "2025-03-29T21:30:00Z",
        "status": "FINISHED",
        "matchday": 1,
        "stage": "REGULAR_SEASON",
        "group": "",
        "lastUpdated": "2025-09-26T00:20:44Z",
        "homeTeam": {
            "id": 1776,
            "name": "São Paulo FC",
            "shortName": "São Paulo",
            "tla": "PAU",
            "crest": "https://crests.football-data.org/1776.png"
        },
        "awayTeam": {
            "id": 1778,
            "name": "SC Recife",
            "shortName": "Recife",
            "tla": "REC",
            "crest": "https://crests.football-data.org/1778.png"
        },
        "score": {
            "winner": "DRAW",
            "duration": "REGULAR",
            "fullTime": {
                "home": 0,
                "away": 0
            },
            "halfTime": {
                "home": 0,
                "away": 0
            }
        },
        "competition": {
            "id": 2013,
            "name": "Campeonato Brasileiro Série A",
            "code": "BSA",
            "type": "LEAGUE",
            "emblem": "https://crests.football-data.org/bsa.png",
            "currentSeason": null,
            "seasons": null
        },
        "season": {
            "id": 2371,
            "startDate": "2025-03-29",
            "endDate": "2025-12-21",
            "currentMatchday": 24,
            "winner": null
        }
    },
    ...
]
```

---

## Fans

### `POST api/v1/fans`

Subscribes a fan to a championship. This endpoint requires authentication.

**Headers:**
Authorization: Bearer YOUR_JWT_TOKEN_HERE

**Request Body:**

```json
{
  "user_id": "integer",
  "team_id": "integer",
  "team_name": "São Paulo",
}
```

**Response:**

```json
{
  "message": "Successfully subscribed to São Paulo"
}
```

---

## Admin

These endpoints require administrator privileges.

### `GET api/v1/admin/`

Retrieves a list of all matches.

**Headers:**
Authorization: Bearer YOUR_ADMIN_JWT_TOKEN_HERE

**Response:**

```json
[
  {
        "id": 534938,
        "utcDate": "2025-03-29T21:30:00Z",
        "status": "FINISHED",
        "matchday": 1,
        "stage": "REGULAR_SEASON",
        "group": "",
        "lastUpdated": "2025-09-26T00:20:44Z",
        "homeTeam": {
            "id": 1776,
            "name": "São Paulo FC",
            "shortName": "São Paulo",
            "tla": "PAU",
            "crest": "https://crests.football-data.org/1776.png"
        },
        "awayTeam": {
            "id": 1778,
            "name": "SC Recife",
            "shortName": "Recife",
            "tla": "REC",
            "crest": "https://crests.football-data.org/1778.png"
        },
        "score": {
            "winner": "DRAW",
            "duration": "REGULAR",
            "fullTime": {
                "home": 0,
                "away": 0
            },
            "halfTime": {
                "home": 0,
                "away": 0
            }
        },
        "competition": {
            "id": 2013,
            "name": "Campeonato Brasileiro Série A",
            "code": "BSA",
            "type": "LEAGUE",
            "emblem": "https://crests.football-data.org/bsa.png",
            "currentSeason": null,
            "seasons": null
        },
        "season": {
            "id": 2371,
            "startDate": "2025-03-29",
            "endDate": "2025-12-21",
            "currentMatchday": 24,
            "winner": null
        }
    },
    ...
]
```

```json
// ws will update this 
{
    "channel_id": 1,
    "total_sent": 10,
    "sent_count": 5,
    "failed_count": 5,
    "is_completed": true,
    "error_details": []
}
```

### `POST api/v1/admin/broadcast/:match_id`

Broadcasts a match.

**Headers:**
Authorization: Bearer YOUR_ADMIN_JWT_TOKEN_HERE

**Response:**

```json
{
    "message": "Broadcast sent to 1 fans",
    "data": {
        "match_id": 534938,
        "message": "🏆 São Paulo FC vs SC Recife - Status: FINISHED",
        "notification_id": "match_534938",
        "targets_count": 1
    }
}
```
