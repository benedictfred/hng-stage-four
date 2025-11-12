# Push Notification Service

A microservice for sending push notifications via Firebase Cloud Messaging (FCM) with RabbitMQ queue support.

## Features

- **Device Management**: Register and manage device tokens for push notifications
- **Queue-Based Processing**: Asynchronous push notification processing using RabbitMQ
- **Token Validation**: Automatic token validation during registration and before sending
- **Rich Notifications**: Support for title, body, image, and link in notifications
- **Retry Mechanism**: Automatic retry with exponential backoff (max 5 retries)
- **Dead Letter Queue**: Failed messages after max retries are moved to DLQ
- **Queue Statistics**: Monitor queue lengths and processing status
- **API Documentation**: Interactive Swagger/OpenAPI documentation

## Prerequisites

- Go 1.24+
- PostgreSQL 15+
- RabbitMQ 3.x
- Firebase Cloud Messaging credentials (service account JSON file)

## Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd push-service
   ```

2. **Set up FCM credentials**:
   - Place your Firebase service account JSON file as `service-account.json` in the project root
   - Or set `FCM_CREDENTIALS_JSON` environment variable with the JSON content

3. **Start all services**:
   ```bash
   make docker-compose-up
   ```
   Or manually:
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**:
   ```bash
   make migrate-up
   ```

5. **Access the services**:
   - API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html
   - RabbitMQ Management: http://localhost:15672 (guest/guest)

### Local Development

1. **Install dependencies**:
   ```bash
   go mod download
   ```

2. **Set up environment variables**:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=push_service
   export DB_PASSWORD=push_service_password
   export DB_NAME=push_service
   export RABBITMQ_HOST=localhost
   export RABBITMQ_PORT=5672
   export FCM_USE_FILE=true
   # Set FCM_CREDENTIALS_JSON or place service-account.json in project root
   ```

3. **Start PostgreSQL and RabbitMQ** (if not using Docker):
   ```bash
   # Using Docker Compose for dependencies only
   docker-compose up -d postgres rabbitmq
   ```

4. **Run database migrations**:
   ```bash
   make migrate-up
   ```

5. **Generate Swagger documentation**:
   ```bash
   make swagger
   ```

6. **Run the application**:
   ```bash
   make run
   ```

## API Documentation

### Swagger UI

Once the service is running, access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

The Swagger UI provides:
- Complete API endpoint documentation
- Request/response schemas
- Example requests
- Try-it-out functionality

### API Endpoints

#### Health Checks
- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check (includes database connectivity)

#### Device Management
- `POST /v1/devices` - Register a new device
- `GET /v1/devices?user_id={user_id}` - Get user's devices
- `DELETE /v1/devices/{token}` - Unregister a device

#### Push Notifications
- `POST /v1/push/send` - Send push notification to a user (queued)
- `POST /v1/push/send-bulk` - Send push notifications to multiple users (queued)
- `POST /v1/push/test-direct` - Test direct FCM send (bypasses queue)

#### Queue Management
- `GET /v1/queue/stats` - Get queue statistics

### Example API Calls

#### Register a Device
```bash
curl -X POST http://localhost:8080/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "token": "fcm_device_token_here",
    "platform": "android"
  }'
```

#### Send Push Notification
```bash
curl -X POST http://localhost:8080/v1/push/send \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "title": "Hello",
    "body": "This is a test notification",
    "image": "https://example.com/image.jpg",
    "link": "https://example.com"
  }'
```

#### Get Queue Statistics
```bash
curl http://localhost:8080/v1/queue/stats
```

## Docker

### Building the Image

```bash
make docker-build
# or
docker build -t push-service .
```

### Running with Docker

```bash
make docker-run
# or
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e RABBITMQ_HOST=rabbitmq \
  -v $(pwd)/service-account.json:/app/service-account.json:ro \
  push-service
```

### Docker Compose

The `docker-compose.yml` file includes:
- **push-service**: Main application
- **postgres**: PostgreSQL database
- **rabbitmq**: RabbitMQ message broker

**Commands**:
```bash
# Start all services
make docker-compose-up

# Stop all services
make docker-compose-down

# View logs
make docker-compose-logs
# or
docker-compose logs -f

# Rebuild and start
make docker-compose-build
make docker-compose-up
```

## Configuration

Configuration is managed via `config.yaml` and environment variables. Key settings:

### Server
- `SERVER_PORT`: HTTP server port (default: 8080)
- `SERVER_MODE`: Gin mode (debug/release)

### Database
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name

### RabbitMQ
- `RABBITMQ_HOST`: RabbitMQ host
- `RABBITMQ_PORT`: RabbitMQ port (default: 5672)
- `RABBITMQ_USERNAME`: RabbitMQ username
- `RABBITMQ_PASSWORD`: RabbitMQ password
- `RABBITMQ_VHOST`: Virtual host (default: /)

### Queue
- `QUEUE_WORKER_PREFETCH_COUNT`: Number of messages to prefetch (default: 10)
- `QUEUE_RETRY_MAX_RETRIES`: Maximum retry attempts (default: 5)
- `QUEUE_RETRY_BACKOFF`: Retry backoff duration (default: 5s)
- `QUEUE_VALIDATION_ENABLED`: Enable token validation (default: true)

### FCM
- `FCM_USE_FILE`: Use service account file (true/false)
- `FCM_CREDENTIALS_JSON`: FCM credentials as JSON string (alternative to file)
- `FCM_PROJECT_ID`: Firebase project ID

## Development

### Generate Swagger Documentation

```bash
make swagger
```

This generates Swagger documentation in `docs/swagger/` directory.

### Running Tests

```bash
make test
```

### Database Migrations

```bash
# Create a new migration
make migrate-create

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## Architecture

### Queue Processing Flow

1. **Enqueue**: Push notifications are enqueued to RabbitMQ
2. **Worker**: Background worker consumes messages from the queue
3. **Validation**: Device tokens are validated (if enabled)
4. **Send**: Notifications are sent via FCM
5. **Retry**: Failed messages are retried with exponential backoff
6. **DLQ**: Messages exceeding max retries are moved to dead letter queue

### Queue Structure

- **Main Queue**: `push_notifications_queue` - Primary queue for new notifications
- **Retry Queue**: `push_retries_queue` - Messages waiting for retry
- **Dead Letter Queue**: `push_dead_letters_queue` - Failed messages after max retries

## License

MIT
