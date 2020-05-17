binaries = gohome

.PHONY: clean gohome

all: $(binaries)

gohome: cmd/gohome/main.go
	go build -o $@ $<

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
	docker rm gohome; true
	docker run -d --name gohome -p 8053:8080 --env-file ./env.list ericvolp12/gohome