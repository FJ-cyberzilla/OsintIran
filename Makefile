.PHONY: build run test clean deploy logs

build:
	@docker build -t my-app:latest .

run:
	@docker run -p 8080:8080 my-app:latest

dev:
	@docker build -t my-app:dev .
	@docker run -p 8080:8080 -v $(PWD):/app my-app:dev

deploy:
	@./deploy.sh

logs:
	@docker logs -f my-app

clean:
	@docker system prune -f

test:
	@docker build -t my-app:test .
	@docker run --rm my-app:test go test ./...

status:
	@docker ps -f name=my-app
