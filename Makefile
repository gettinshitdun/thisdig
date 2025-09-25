build:
	go build -o main main.go

run: build
	time ./main rtb-sc.media.net sc && rm -rf main
