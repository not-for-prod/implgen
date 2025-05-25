test:
	rm -rf example/out/*
	go run main.go basic --src example/in/interface.go \
		--dst example/out/ \
		--interface-name TestInterface \
		--impl-name Test \
		--impl-package test \
		--tracer-name my-brilliant-tracer