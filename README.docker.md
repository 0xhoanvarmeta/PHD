# Docker Setup Guide

This guide helps you run the PHD Admin system using Docker and Docker Compose.

## Prerequisites

- Docker installed ([Get Docker](https://docs.docker.com/get-docker/))
- Docker Compose installed (included with Docker Desktop)

## Quick Start (Development)

### 1. Setup Environment Variables

```bash
# Copy the docker environment template
cp .env.docker .env

# Edit .env and add your blockchain configuration
# IMPORTANT: Change JWT_SECRET in production!
```

### 2. Start Services

```bash
# Start PostgreSQL and Backend
docker-compose up -d

# View logs
docker-compose logs -f

# Or start with logs
docker-compose up
```

The services will:
- ✅ Start PostgreSQL on port `5432`
- ✅ Automatically run database migrations
- ✅ Seed the default admin account
- ✅ Start backend server on port `3000`

### 3. Access Services

- **API**: http://localhost:3000/api
- **Swagger Docs**: http://localhost:3000/api/docs
- **PostgreSQL**: localhost:5432
  - Database: `phd_admin`
  - User: `postgres`
  - Password: `postgres`

### 4. Default Admin Account

```
Username: admin
Password: Admin@123
```

**⚠️ IMPORTANT**: Change this password after first login!

## Development Commands

```bash
# Stop services
docker-compose down

# Stop and remove volumes (clears database)
docker-compose down -v

# Restart backend only
docker-compose restart backend

# View backend logs
docker-compose logs -f backend

# View postgres logs
docker-compose logs -f postgres

# Execute command in backend container
docker-compose exec backend sh
docker-compose exec backend npm run migration:run

# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d phd_admin
```

## Production Deployment

### 1. Setup Production Environment

```bash
# Create production .env file
cp .env.docker .env

# Edit .env with production values
nano .env
```

**Required production settings:**
```env
# Strong JWT secret (use: openssl rand -base64 32)
JWT_SECRET=your-super-secure-random-secret-key

# Database credentials (change from defaults!)
DB_USERNAME=phd_admin
DB_PASSWORD=strong-password-here
DB_DATABASE=phd_admin

# Blockchain configuration
BLOCKCHAIN_NETWORK=mainnet
CONTRACT_ADDRESS=0x...
PRIVATE_KEY=0x...

# CORS (restrict to your domain)
CORS_ORIGIN=https://yourdomain.com

# Logging
LOG_LEVEL=info
NODE_ENV=production
```

### 2. Start Production Stack

```bash
# Build and start production containers
docker-compose -f docker-compose.prod.yml up -d

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### 3. Production Commands

```bash
# Stop production services
docker-compose -f docker-compose.prod.yml down

# View production logs
docker-compose -f docker-compose.prod.yml logs -f backend

# Run migrations in production
docker-compose -f docker-compose.prod.yml exec backend npm run migration:run

# Backup database
docker-compose -f docker-compose.prod.yml exec postgres pg_dump -U postgres phd_admin > backup.sql

# Restore database
docker-compose -f docker-compose.prod.yml exec -T postgres psql -U postgres phd_admin < backup.sql
```

## Database Management

### Manual Migrations

```bash
# Run migrations
docker-compose exec backend npm run migration:run

# Revert last migration
docker-compose exec backend npm run migration:revert

# Generate new migration
docker-compose exec backend npm run migration:generate -- apps/backend/src/app/database/migrations/YourMigrationName
```

### Database Backup & Restore

```bash
# Backup
docker-compose exec postgres pg_dump -U postgres phd_admin > backup_$(date +%Y%m%d_%H%M%S).sql

# Restore
docker-compose exec -T postgres psql -U postgres phd_admin < backup.sql

# Create new database
docker-compose exec postgres createdb -U postgres new_database_name
```

## Troubleshooting

### Backend won't start

```bash
# Check backend logs
docker-compose logs backend

# Restart backend
docker-compose restart backend

# Rebuild backend
docker-compose up -d --build backend
```

### Database connection issues

```bash
# Check if postgres is healthy
docker-compose ps

# Check postgres logs
docker-compose logs postgres

# Verify database exists
docker-compose exec postgres psql -U postgres -l
```

### Reset everything

```bash
# Stop and remove all containers, networks, and volumes
docker-compose down -v

# Remove images
docker-compose down --rmi all -v

# Start fresh
docker-compose up -d
```

## Docker Compose Files

- `docker-compose.yml` - Development setup with hot-reload
- `docker-compose.prod.yml` - Production setup with optimized build
- `Dockerfile.backend` - Development Dockerfile
- `Dockerfile.backend.prod` - Production Dockerfile (multi-stage build)

## Performance Tips

### Development
- Hot reload is enabled by default
- Source code is mounted as volume for instant changes
- Database data persists in named volume

### Production
- Use multi-stage build for smaller image size
- Only production dependencies installed
- Optimized for performance and security
- Health checks enabled

## Security Checklist for Production

- [ ] Change default PostgreSQL password
- [ ] Generate strong JWT secret
- [ ] Update admin password after first login
- [ ] Configure proper CORS origins
- [ ] Use environment-specific `.env` file
- [ ] Enable SSL/TLS in production
- [ ] Set up firewall rules
- [ ] Regular database backups
- [ ] Monitor logs for security issues
- [ ] Keep Docker images updated

## Environment Variables Reference

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | postgres | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USERNAME` | postgres | Database username |
| `DB_PASSWORD` | postgres | Database password |
| `DB_DATABASE` | phd_admin | Database name |
| `JWT_SECRET` | dev-secret | JWT signing secret |
| `JWT_EXPIRES_IN` | 24h | JWT expiration time |
| `PORT` | 3000 | Backend server port |
| `NODE_ENV` | development | Environment mode |
| `BLOCKCHAIN_NETWORK` | testnet | Blockchain network |
| `CONTRACT_ADDRESS` | - | Smart contract address |
| `PRIVATE_KEY` | - | Blockchain private key |
| `CORS_ORIGIN` | * | CORS allowed origins |
| `LOG_LEVEL` | debug | Logging level |

## Next Steps

After starting the services:

1. **Test the API**: Visit http://localhost:3000/api/docs
2. **Login**: Use default admin credentials
3. **Change Password**: Update admin password immediately
4. **Create Scripts**: Start creating and managing scripts
5. **Monitor Events**: Check event logs for blockchain triggers

## Support

For issues or questions:
- Check logs: `docker-compose logs -f`
- View container status: `docker-compose ps`
- Restart services: `docker-compose restart`
- Report issues: [GitHub Issues](https://github.com/your-repo/issues)
