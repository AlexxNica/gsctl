PWD := $(shell pwd)

get-deps:
	go get github.com/go-resty/resty

generate:
	docker run --rm -it \
		-v ${PWD}:/swagger-api/out \
		-v ${PWD}/api-spec:/swagger-api/yaml \
		jimschubert/swagger-codegen-cli generate \
		--input-spec /swagger-api/yaml/oai-spec.yaml \
		--lang go \
		--config /swagger-api/out/swagger-codegen-conf.json \
		--output /swagger-api/out
	gofmt -s -l -w .

validate:
	docker run --rm -it \
		-v ${PWD}/api-spec:/swagger-api/yaml \
		jimschubert/swagger-codegen-cli generate \
		--input-spec /swagger-api/yaml/oai-spec.yaml \
		--lang swagger \
		--output /tmp/

build:
	go build ./...