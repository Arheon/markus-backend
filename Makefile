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

build-dispatcher:
	docker compose exec go_api go build -o ./build/dispatcher -tags=musl cmd/dispatcher/dispatcher.go

run:
	docker compose exec go_api ./build/app

run-dispatcher:
	docker compose exec go_api ./build/dispatcher

bash:
	docker compose exec go_api bash
