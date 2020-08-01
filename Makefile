.PHONY: frontend
frontend:
	cd pkg/www/retext && npm run build

.PHONY: backend
backend: frontend
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

.PHONY: docker_server
docker_server:
	docker build . -f deployment/Dockerfile -t retext

.PHONY: docker_db_loader
docker_db_loader:
	docker build . -f deployment/db_loader.Dockerfile -t qode_db_loader

.PHONY: devserve
devserve:
	cd pkg/www/retext && npm run serve

