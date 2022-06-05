# Storage up if you don't have native database in your host
storage-up:
	docker-compose -f docker-compose-pg-storage.yml up -d --remove-orphans

# Storage down to deactivate it
storage-down:
	docker-compose -f docker-compose-pg-storage.yml down

# Storage test up (Storage that dedicated for testing)
storage-up-test:
	docker-compose -f docker-compose-pg-storage-test.yml up -d --remove-orphans

# Storage down to deactivate storage for testing
storage-down-test:
	docker-compose -f docker-compose-pg-storage-test.yml down

server:
	go run main.go