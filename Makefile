docker:
	docker build -t aldrismcloud/rwaldi-go-sse:latest .
	docker push aldrismcloud/rwaldi-go-sse:latest
test:
	docker run aldrismcloud/rwaldi-go-sse:latest
compose:
	docker compose up
build:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o gen/sse.exe .
	CGO_ENABLED=0 GOOS=linux go build -o gen/sse .; ls -lh gen/sse