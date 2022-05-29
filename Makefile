storage:
	docker-compose -f docker-compose-mysql-storage.yml up -d

server:
	go run main.go