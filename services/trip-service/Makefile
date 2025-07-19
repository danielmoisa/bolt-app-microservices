
# Temporary access for pgAdmin
forward-postgres:
	kubectl port-forward service/postgres 5432:5432

migrate:
	migrate -path ./migrations -database "postgresql://postgres:password@localhost:5432/bolt_app?sslmode=disable" up
