# Running

		$ go run main.go
		$ cd plugin/example && go run main.go

# Debugging

		$ open http://localhost:8081/debug/requests
		$ open http://localhost:8081/debug/events

# TODO

1. Configuration management
  - [ ] setup for viper
2. Plugin lifecycle management
	- [ ] server boots plugins
	- [x] plugin registers itself with the server
	- [ ] plugins detect server disconnect and exit
3. Command line
  - [ ] pass args from server to plugins
4. Autopatching
	- https://github.com/docker/Notary
	- https://godoc.org/github.com/inconshreveable/go-update
