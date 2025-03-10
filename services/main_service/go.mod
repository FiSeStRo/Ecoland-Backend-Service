module github.com/FiSeStRo/Ecoland-Backend-Service/services/main-service

go 1.23.4

require (
	github.com/FiSeStRo/Ecoland-Backend-Service/pkg v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.2.1
	golang.org/x/crypto v0.36.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.9.0 // indirect
)

replace github.com/FiSeStRo/Ecoland-Backend-Service/pkg => ../../go_pkg
