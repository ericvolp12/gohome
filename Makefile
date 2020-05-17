binaries = gohome

.PHONY: clean gohome

all: $(binaries)

gohome: cmd/gohome/main.go
	go build -o $@ $<

gohome-static: cmd/gohome/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $@ $<

clean:
	rm -f $(binaries)

.PHONY: test
test:
	go test -coverprofile cp.out ./...; rm cp.out

.PHONY: build-gohome-image
build-gohome-image:
	docker build -f build/gohome.Dockerfile -t ericvolp12/gohome .

.PHONY: run-gohome
run-gohome:
	docker rm -f gohome; true
	docker run -d --restart=always --name gohome --network host --env-file ./env.list ericvolp12/gohome