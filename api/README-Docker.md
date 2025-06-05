# Docker Setup for Learning Go API

## 🚀 Quick Start

### 1. Start Everything with Docker Compose

```bash
docker-compose up -d
```

### 2. View Logs

```bash
# View all logs
docker-compose logs -f

# View only API logs
docker-compose logs -f api

# View only MongoDB logs
docker-compose logs -f mongodb
```

### 3. Stop Everything

```bash
docker-compose down
```

### 4. Stop and Remove Volumes (Clean Reset)

```bash
docker-compose down -v
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the api directory with:

```env
# MongoDB Configuration (already configured in docker-compose.yml)
MONGO_URI=mongodb://admin:password123@mongodb:27017/learning_go_db?authSource=admin
MONGO_DATABASE=learning_go_db

# For external MongoDB, use:
# MONGO_URI=mongodb://username:password@your-mongo-host:27017/your-database?authSource=admin
# MONGO_DATABASE=your_database_name
```

### Default MongoDB Credentials

- **Username**: `admin`
- **Password**: `password123`
- **Database**: `learning_go_db`
- **Port**: `27018` (external), `27017` (internal)

## 📡 Access Points

- **API**: https://localhost:8443
- **MongoDB**: localhost:27018

## 🔧 Development Commands

### Rebuild API after code changes

```bash
docker-compose up -d --build api
```

### Access MongoDB shell

```bash
docker-compose exec mongodb mongosh -u admin -p password123
```

### Access API container shell

```bash
docker-compose exec api sh
```

## 📊 Services

### MongoDB

- **Image**: mongo:7.0
- **Container**: learning_go_mongodb
- **Data Persistence**: Volume `mongodb_data`
- **Health Check**: Included

### Go API

- **Build**: From local Dockerfile
- **Container**: learning_go_api
- **Port**: 8443 (HTTPS)
- **Dependencies**: Waits for MongoDB to be healthy

## 🛠 Troubleshooting

### Check service status

```bash
docker-compose ps
```

### Restart a specific service

```bash
docker-compose restart api
# or
docker-compose restart mongodb
```

### Clean rebuild everything

```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```
