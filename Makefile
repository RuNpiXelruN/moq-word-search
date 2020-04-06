.PHONY: proto cover test cleanproto cleanbinaries clean binaries build buildrun generate mocks run cover

binaries: ## Builds binaries for osx, linux, and windows, and places them in the ./bin folder.
	cd handler && \
	env GOOS=darwin GOARCH=amd64 go build -v -o grpcWordSearch_osx && \
	env GOOS=linux GOARCH=amd64 go build -v -o grpcWordSearch_linux && \
	env GOOS=windows GOARCH=amd64 go build -v -o grpcWordSearch_windows && \
	mv ./*grpcWordSearch* ../bin/

build: test binaries ## Generates interface mocks, runs test suite, and builds binaries for osx, linux, and windows, and places them in the ./bin folder.

buildrun: build run ## Runs tests, builds binaries and runs osx version of binary

clean: ## Removes all generated proto related files (excluding original .proto file) from ./proto folder, and all binaries from ./bin folder
	rm -r ./proto/{*.go,*.json} && \
	rm -r ./bin/* && \
	rm -r ./mocks/*

cleanbinaries: ## Removes all binaries from ./bin folder
	rm -r ./bin/*

cleanmocks: ## Removes all binaries from ./bin folder
	rm -r ./mocks/*

cleanproto:	## Removes all generated proto related files (excluding original .proto file) from ./proto folder
	rm -r ./proto/{*.go,*.json}

cover: ## Runs test suite and opens html coverage profile in browser
	go test ./... -coverprofile cover.out && \
	go tool cover -html cover.out

generate: proto build run ## Generates proto output, generates interface mocks, runs tests, builds binaries, and runs osx version of binary

help: routes ## Display available commands	
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

mocks: ## Generates interface mocks for testing
	go generate
  
proto: ## Generates proto api, grpc-gateway definition, and swagger definition 
	protoc \
		-I/usr/local/include \
		-I. \
		-I${GOPATH} \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:./proto proto/moqwordsearch.proto  \
		--grpc-gateway_out=./proto \
		--swagger_out=logtostderr=true:./proto

		mv ./proto/proto/* ./proto/
		mv ./proto/github.com/RuNpiXelruN/moq-word-search/* ./proto/
		rm -r ./proto/{proto,github.com}

routes: ## Displays endpoints of word search app
	$(info ) \
	$(info ************ Routes ************) \
	$(info ) \
	$(info curl http://localhost:8090/api/words?term={your_term}) \
	$(info curl http://localhost:8090/api/words/popular) \
	$(info curl -X "POST" "http://localhost:8090/api/words/sawyer") \
	$(info )

run: ## Runs osx version of binary
	cd bin && \
	./grpcWordSearch_osx

test: ## Runs test suite
	go test -v ./...
