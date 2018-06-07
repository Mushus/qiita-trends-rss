.PHONY: dep, build

init:
	go get -u github.com/golang/dep/cmd/dep
dep:
	dep ensure
run:
	go run ./...
build:
	mkdir -p ./build
	GOARCH=amd64 GOOS=linux GOARM=6 go build -o ./build/qiita-trends-rss-amd64-linux
	GOARCH=arm GOOS=linux GOARM=6 go build -o ./build/qiita-trends-rss-arm-linux
