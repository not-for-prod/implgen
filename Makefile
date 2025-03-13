linter:
	golangci-lint --config .golangci.yaml run

test:
	go run main.go --src ./example/in/aboba.go --dst ./example/out --with-otel