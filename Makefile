start:
	go run .

dev:
	go run . -dev

rebuild: build-client
	go run . -dev

build-client:
	cd web && pnpm run build