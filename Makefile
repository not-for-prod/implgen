linter:
	golangci-lint --config .golangci.yaml run

test_basic:
	go run main.go --src ./example/in/aboba.go --dst ./example/out/basic --interface-name Aboba --with-otel

test_repo:
	go run main.go repo --src ./example/in/aboba.go --dst ./example/out/repo --interface-name AbobaRepository
	