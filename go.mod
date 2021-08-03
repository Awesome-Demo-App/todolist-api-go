module github.com/awesome-demo-app/todolist-api

go 1.16

require (
	github.com/prometheus/client_golang v1.11.0
	github.com/slok/go-http-metrics v0.9.0
	go.opentelemetry.io/otel v1.0.0-RC2
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.0.0-RC2
	go.opentelemetry.io/otel/sdk v1.0.0-RC2
	go.opentelemetry.io/otel/trace v1.0.0-RC2
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12
)
