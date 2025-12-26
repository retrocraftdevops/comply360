module github.com/comply360/commission-service

go 1.21

require (
	github.com/comply360/shared v0.0.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.5.0
	github.com/lib/pq v1.10.9
	github.com/rabbitmq/amqp091-go v1.9.0
	github.com/shopspring/decimal v1.3.1
)

replace github.com/comply360/shared => ../../packages/shared
