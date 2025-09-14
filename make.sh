
case "$1" in
  start)
    go run cmd/main.go
    ;;
  format)
    go fmt ./...
    ;;
  build)
    docker build -t netzero .
    ;;
  run-docker)
    docker stop netzero-local
    docker rm netzero-local
    docker run -p 8080:8080 --name netzero-local netzero
    ;;
  stop-docker)
  docker stop 


esac
