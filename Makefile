up:
	docker compose up -d

down:
	docker compose down

build: 
	docker compose build

create-env:
	cp .env.example .env

download-dependencies:
	docker compose exec go_api go mod download

build-app:
	docker compose exec go_api go build -o ./build/app -tags=musl cmd/app/app.go

run:
	docker compose exec go_api ./build/app
bash:
	docker compose exec go_api bash
