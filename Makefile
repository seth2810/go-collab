install:
	cd frontend && yarn install --frozen-lockfile

start-backend:
	cd backend && go run main.go serve

start-frontend:
	cd frontend && yarn start

start:
	heroku local -f Procfile.dev

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

