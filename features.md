Berdasarkan project requirements yang ada di prompting.md, berikut adalah fitur-fitur yang belum diimplementasikan:

1. **Auth Feature**:
- User registration dengan email/password validation
- Login dengan JWT token menggunakan JWK
- JWK rotation capability
- Refresh token mechanism
- Password hashing dengan bcrypt

2. **Cart Feature** (belum diimplementasikan sama sekali):
- Add/Remove/Update cart items
- Calculate cart total
- Expire cart after certain time
- Store cart data in MongoDB
- Publish events saat cart changes

3. **Order Feature** (belum diimplementasikan sama sekali):
- Create order dari cart
- Order status management (PENDING → CONFIRMED → PROCESSING → SHIPPED → DELIVERED)
- Store order data di PostgreSQL
- Publish order events

4. **Payment Feature** (belum diimplementasikan sama sekali):
- Process payment simulation
- Payment status management
- Store payment records
- Publish payment events

5. **Additional Message Brokers** (baru RabbitMQ yang diimplementasikan):
- Kafka integration
- NSQ integration
- NATS integration

6. **Database Implementations**:
- MongoDB integration
- Database migrations untuk semua tables

7. **Monitoring & Observability**:
- Structured logging dengan logrus/zap
- Circuit breaker pattern
- Health check endpoints

8. **Testing**:
- Unit tests untuk setiap usecase
- Integration tests
- Mock repositories
- Test fixtures

9. **API Server**:
- HTTP handlers dengan Gin
- Input validation
- Error handling middleware
- API documentation

10. **Configuration**:
- Viper configuration management
- Environment-based configs
- Multiple broker configs

Saat ini yang sudah diimplementasikan dengan baik adalah:
1. Saga Pattern dengan orchestration
2. RabbitMQ message broker
3. Basic PostgreSQL integration
4. Worker service untuk saga steps
5. Docker setup dengan docker-compose
6. Basic deployment configuration

Apakah Anda ingin saya membantu mengimplementasikan salah satu fitur yang belum tersedia?
