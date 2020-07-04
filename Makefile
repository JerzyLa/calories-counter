clean:
		@rm -rf out
		@mkdir -p out

generate:
		go generate ./...

lint:
		go fmt ./...
		golangci-lint run

build: clean
		go build -o out/

test: generate
		go test `go list ./... | grep -v  calories-counter/test` --cover

e2e_test:
		go test calories-counter/test

start: build
		./out/calories-counter
