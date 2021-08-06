build:
	go build

test:
	go test ./rpc ./exporter -v

coverage:
	go test ./rpc ./exporter -coverprofile=coverage.out

image:
	docker build -t tezos-prometheus-exporter .
