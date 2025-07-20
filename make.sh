
case "$1" in
  start)
    go run cmd/main.go
    ;;
  format)
    go fmt ./...
esac
