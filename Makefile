build:
	go build

test:
	go test ./rpc ./exporter ./config -v

coverage:
	go test ./rpc ./exporter ./config -coverprofile=coverage.out

image:
	docker build -t tezos-prometheus-exporter .
