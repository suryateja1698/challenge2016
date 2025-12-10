run:
	CGO_ENABLED=0 go build -o distrib ./cmd/app/
	./distrib
