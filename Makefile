# config data
user := root
database := test

dev:
	docker-compose up -d
	@echo "docker development setup started."
	
dev-build:
	docker-compose up --build -d
	@echo "docker compose image rebuilded."

stop:
	docker-compose down
	@echo "Docker compose Stopped"

restart-db: delete-data-db create-db insert-data-db