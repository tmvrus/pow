
lint:
	golangci-lint run --config .golangci.yml

test:
	go test ./...

start:
	docker inspect pow-net || docker network create --driver bridge pow-net
	docker build -t pow-server -q -f Dockerfile.server .
	docker build -t pow-client -q -f Dockerfile.client .
	docker run -d --rm --net=pow-net --name=pow-server -p 2222:22222 -v $(PWD):/app pow-server

client:
	docker run --net=pow-net --name pow-client -e SERVER_ADDRESS='pow-server:22222'  --rm -v $(PWD):/app pow-client
