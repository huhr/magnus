install:
	go build -o magnus

vendor:
	@go get -u github.com/BurntSushi/toml
	@go get -u github.com/huhr/simplelog

clean:
	go clean

