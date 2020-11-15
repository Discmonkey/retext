.PHONY: frontend
frontend:
	cd pkg/www/retext && npm run build

.PHONY: backend
backend: frontend
	mkdir -p bin
	go build -o bin/server cmd/server/main.go

.PHONY: docker_server
docker_server:
	docker build . -f deployment/Dockerfile -t qode

.PHONY: docker_db_loader
docker_db_loader:
	docker build . -f deployment/db_loader.Dockerfile -t qode_db_loader

.PHONY: tag
tag:
	git tag $(version)
	echo "package version\n\nconst Version string = \"$(version)\"" > pkg/version/version.go

client-apis:
	java -jar third_party/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l go -D models -o pkg/swagger
	java -jar third_party/swagger-codegen-cli-3.0.20.jar generate -i swagger.yaml -l typescript-angular --additional-properties modelPropertyNaming=snake_case -D models -o pkg/www/retext/src/