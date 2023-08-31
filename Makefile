.PHONY: send build

send:
	echo "hello" | nc -w1 -u localhost 2601

build:
	go build -o bin/udp-test .

docker:
	docker build -t nb-lpetera-u:5000/udp-test:latest . && \
	docker push nb-lpetera-u:5000/udp-test:latest