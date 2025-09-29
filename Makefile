build:
	go build -o bin/main pkg/main.go

run: build
	time bin/main rtb-sc.media.net && rm -rf bin/main