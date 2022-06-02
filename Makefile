storage-up:
	docker-compose -f docker-compose-pg-storage.yml up -d --remove-orphans

storage-down:
	docker-compose -f docker-compose-pg-storage.yml down

server:
	go run main.go