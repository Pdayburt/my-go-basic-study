.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi noah/webook:v0.0.1 || true
	@docker build -t noah/webook:v0.0.1 .

.PHONY: grpc
grpc:
	@buf generate webook/api/proto