build:
	go build

test:
	go test ./rpc ./exporter ./config ./process -v

coverage:
	go test ./rpc ./exporter ./config ./process -coverprofile=coverage.out

image:
	docker build -t tezos-prometheus-exporter .
