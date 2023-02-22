all: build

build: build/macos build/freebsd

build/macos:
	go build -o ./bin/scraping_bar

build/freebsd:
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/freebsd-amd64/scraping_bar

deploy: deploy/foo

deploy/foo:
	scp -p ./bin/freebsd-amd64/scraping_bar foo@foo.sakura.ne.jp:~/work/