
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
    docker-compose down
    docker-compose up --build
    docker-compose down
    ;;
  run-docker-detached)
    docker-compose down
    docker-compose up --build -d
    ;;
  test)
    go run functional-test/test.go
    ;;
esac
