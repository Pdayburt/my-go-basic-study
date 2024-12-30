.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi noah/webook:v0.0.1
	@docker build -t noah/webook:v0.0.1 .
