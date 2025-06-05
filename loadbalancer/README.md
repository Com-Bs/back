# Nginx Load Balancer

This directory contains an nginx load balancer configuration that distributes traffic between two frontend VM instances.

## Configuration

### 1. Update VM Instance IPs

Edit `nginx.conf` and replace the placeholder IP addresses with your actual VM instance IPs:

```nginx
upstream frontend_servers {
    server YOUR_VM1_IP:80;  # Frontend VM Instance 1
    server YOUR_VM2_IP:80;  # Frontend VM Instance 2
}
```

### 2. Load Balancing Methods

The configuration uses round-robin by default. You can change the load balancing method by adding one of these lines to the `upstream` block:

- `least_conn;` - Route to server with least connections
- `ip_hash;` - Route based on client IP hash (sticky sessions)

## Deployment

### Using Docker Compose (Recommended)

1. Navigate to the loadbalancer directory:

   ```bash
   cd loadbalancer
   ```

2. Build and start the load balancer:

   ```bash
   docker-compose up -d
   ```

3. Check the status:

   ```bash
   docker-compose ps
   ```

4. View logs:
   ```bash
   docker-compose logs -f
   ```

### Using Docker directly

1. Build the image:

   ```bash
   docker build -t nginx-loadbalancer .
   ```

2. Run the container:
   ```bash
   docker run -d -p 80:80 --name nginx-loadbalancer nginx-loadbalancer
   ```

## Health Check

The load balancer includes a health check endpoint at `/health`. You can test it:

```bash
curl http://localhost/health
```

## Monitoring

Nginx access and error logs are available in the `./logs` directory when using docker-compose.

## Stopping the Load Balancer

```bash
docker-compose down
```
