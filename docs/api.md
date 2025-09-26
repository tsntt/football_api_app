# API Endpoints
Here are all the API endpoints that are available for the Football Web App and their corresponding responses and examples.
---
## PUBLIC
### Health Check
GET http://localhost:8080/health

### Register User
POST http://localhost:8080/auth/register
Content-Type: application/json

{
  "name": "usuario1",
  "password": "senha123"
}

### Login User
POST http://localhost:8080/auth/login
Content-Type: application/json

{
  "name": "usuario1",
  "password": "senha123"
}

### Login Admin (default admin user)
POST http://localhost:8080/auth/login
Content-Type: application/json

{
  "name": "admin",
  "password": "admin123"
}
---
# PROTECTED
### Get Championships
GET http://localhost:8080/championship
Authorization: Bearer YOUR_JWT_TOKEN_HERE

### Get Matches by Championship ID
GET http://localhost:8080/championship/2021/matches
Authorization: Bearer YOUR_JWT_TOKEN_HERE

### Get Matches with Filters
GET http://localhost:8080/championship/2021/matches?team=Manchester%20United&stage=ROUND_3
Authorization: Bearer YOUR_JWT_TOKEN_HERE

### Subscribe to Team as Fan
POST http://localhost:8080/fans
Authorization: Bearer YOUR_JWT_TOKEN_HERE
Content-Type: application/json

{
  "user_id": 1,
  "team": "Manchester United"
}
---
## ADMIN ONLY
### Admin - Get All Matches (Admin only)
GET http://localhost:8080/admin
Authorization: Bearer YOUR_ADMIN_JWT_TOKEN_HERE

### Admin - Broadcast Match Notification (Admin only)
POST http://localhost:8080/admin/broadcast/327568
Authorization: Bearer YOUR_ADMIN_JWT_TOKEN_HERE

###
### Example responses:
###

### Register Response (200)
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "User successfully created!"
}

### Login Response (200)
HTTP/1.1 200 OK
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX..."
}

### Championship Response (200)
HTTP/1.1 200 OK
Content-Type: application/json

[
  {
    "id": 2021,
    "name": "Premier League",
    "code": "PL",
    "type": "LEAGUE",
    "emblem": "https://crests.football-data.org/PL.png",
    "area": {
      "id": 2072,
      "name": "England",
      "code": "ENG",
      "flag": "https://crests.football-data.org/770.svg"
    },
    "currentSeason": {
      "id": 1564,
      "startDate": "2023-08-11",
      "endDate": "2024-05-19",
      "currentMatchday": 15
    }
  }
]

### Fan Subscribe Response (200)
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "Subscribed to Manchester United"
}

### Broadcast Response (200)
HTTP/1.1 200 OK
Content-Type: application/json

{
  "message": "Broadcast sent to 5 fans",
  "data": {
    "match_id": 327568,
    "fans_count": 5,
    "message": "🏆 Manchester United vs Liverpool - Status: SCHEDULED"
  }
}

### Error Response Example (401)
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "message": "Authorization header required"
}

### Error Response Example (400)
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "message": "validation error: Key: 'UserRequest.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}