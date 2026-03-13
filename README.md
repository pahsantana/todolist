# 📋 Todo List API

## 🏗️ Architecture

The project follows **Clean Architecture** principles, organizing code into layers with well-defined responsibilities:

```
main.go          → application entry point/depency injection
config/           → environment variables
internal/
├── domain/
│   ├── entities/ → entities and business rules
│   └── repository/ → persistence interface
├── services/     → use cases and application rules
├── repositories/ → MongoDB implementation
├── handlers/     → HTTP layer
└── middleware/   → request logger
└── dto/   
```

> The main rule: dependencies always point inward. The `domain` layer knows no other layer.

---

## 🚀 Getting Started

### Prerequisites
- [Go 1.22+](https://go.dev/dl/)
- [Docker](https://www.docker.com/)

### 1. Clone the repository
```bash
git clone https://github.com/pahsantana/todolist.git
cd todolist
```

### 2. Set up environment and edit `.env` with your settings:
```bash
cp .env.example .env
```

### 3. Start the database
```bash
docker compose up -d mongo mongo-express
```

### 4. Install dependencies and run
```bash
go mod tidy
go run main.go
```

API will be available at `http://localhost:8080`

---

## 🐳 Run everything with Docker

```bash
docker compose up -d
```

| Service | URL |
|---|---|
| API | http://localhost:8080 |
| Mongo Express | http://localhost:8081 |

---

## 📡 Endpoints

| Method | Route | Description |
|---|---|---|
| GET | /health | Health check |
| POST | /tasks | Create task |
| GET | /tasks | List tasks |
| GET | /tasks?status=pending | Filter by status |
| GET | /tasks?priority=high | Filter by priority |
| GET | /tasks/:id | Get task by ID |
| PUT | /tasks/:id | Update task |
| DELETE | /tasks/:id | Delete task |

---

## 📦 Data Model

```json
{
  "id": "uuid",
  "title": "Study Golang",
  "description": "Review goroutines",
  "status": "pending",
  "priority": "high",
  "due_date": "2026-12-01",
  "created_at": "2026-03-10T21:00:00Z",
  "updated_at": "2026-03-10T21:00:00Z"
}
```

**Allowed status:** `pending` | `in_progress` | `completed` | `cancelled`

**Priority levels:** `low` | `medium` | `high`

---

## ✅ Business Rules

- Title is required (min. 3 / max. 100 characters)
- Priority is required
- `due_date` cannot be in the past (format `YYYY-MM-DD`)
- `completed` tasks cannot be edited — only deleted
- Status and priority are validated against allowed values

---

## 🔁 Request Examples

# API Requests

## Health Check

**curl:**
```bash
curl http://localhost:8080/health
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET
```

---

## Create Task

**curl:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Study Golang","priority":"high","due_date":"2026-12-01"}'
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"title":"Study Golang","priority":"high","due_date":"2026-12-01"}'
```

---

## List Tasks

**curl:**
```bash
curl http://localhost:8080/tasks
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method GET
```

---

## Filter by Status

**curl:**
```bash
curl "http://localhost:8080/tasks?status=pending"
curl "http://localhost:8080/tasks?status=in_progress"
curl "http://localhost:8080/tasks?status=completed"
curl "http://localhost:8080/tasks?status=cancelled"
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks?status=pending" -Method GET
Invoke-WebRequest -Uri "http://localhost:8080/tasks?status=in_progress" -Method GET
Invoke-WebRequest -Uri "http://localhost:8080/tasks?status=completed" -Method GET
Invoke-WebRequest -Uri "http://localhost:8080/tasks?status=cancelled" -Method GET
```

---

## Filter by Priority

**curl:**
```bash
curl "http://localhost:8080/tasks?priority=low"
curl "http://localhost:8080/tasks?priority=medium"
curl "http://localhost:8080/tasks?priority=high"
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks?priority=low" -Method GET
Invoke-WebRequest -Uri "http://localhost:8080/tasks?priority=medium" -Method GET
Invoke-WebRequest -Uri "http://localhost:8080/tasks?priority=high" -Method GET
```

---

## Filter by Status and Priority

**curl:**
```bash
curl "http://localhost:8080/tasks?status=pending&priority=high"
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks?status=pending&priority=high" -Method GET
```

---

## Get Task by ID

**curl:**
```bash
curl http://localhost:8080/tasks/{id}
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method GET
```

---

## Update Task

**curl:**
```bash
# Update status
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'

# Update title
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title"}'

# Update multiple fields
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated title","status":"in_progress","priority":"low"}'
```

**PowerShell:**
```powershell
# Update status
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"status":"in_progress"}'

# Update title
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"title":"Updated title"}'

# Update multiple fields
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"title":"Updated title","status":"in_progress","priority":"low"}'
```

---

## Delete Task

**curl:**
```bash
curl -X DELETE http://localhost:8080/tasks/{id}
```

**PowerShell:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method DELETE
```

---

## Error Cases

**Invalid priority:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"title":"Test","priority":"urgent"}'
# Expected: 400 Bad Request
```

**Due date in the past:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"title":"Test","priority":"high","due_date":"2020-01-01"}'
# Expected: 400 Bad Request
```

**Task not found:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks/invalid-id" -Method GET
# Expected: 404 Not Found
```

**Edit completed task:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks/{id}" -Method PUT -Headers @{"Content-Type"="application/json"} -Body '{"status":"completed"}'
# First complete the task, then try to edit it
# Expected: 422 Unprocessable Entity
```

**Title too short:**
```powershell
Invoke-WebRequest -Uri "http://localhost:8080/tasks" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"title":"ab","priority":"high"}'
# Expected: 400 Bad Request
```

## 🛠️ Tech Stack

| Technology | Purpose |
|---|---|
| Go 1.22 | Main language |
| Gin | HTTP framework |
| MongoDB 7 | Database |
| Zap | Structured logging |
| Docker | Containerization |
