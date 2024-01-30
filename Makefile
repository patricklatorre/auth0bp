run:
	go run .

rebuild: build-client
	go run .

build-client:
	cd web && pnpm run build