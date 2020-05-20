.PHONY: frontend
frontend:
	cd pkg/www/retext && npm run build

.PHONY: backend
backend: frontend
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

.PHONY: docker
docker:
	docker build . -f deployment/Dockerfile -t retext