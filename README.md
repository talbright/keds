# Running

		$ go run main.go
		$ cd plugin/example && go run main.go

# Debugging

		$ open http://localhost:8081/debug/requests
		$ open http://localhost:8081/debug/events

# TODO

1. Configuration management
  - [x] setup for viper
	- [x] config for plugin path
2. Plugin lifecycle management
	- [ ] server loads plugins
	- [x] plugin registers itself with the server
	- [ ] plugins detect server disconnect and exit
4. Testing framework
	- [ ] add ginko
	- [ ] tests for client package
	- [ ] tests for plugin package
	- [ ] tests for server package
	- [ ] tests for utils package
3. Command line
  - [ ] pass args from server to plugins
4. Versioning
  - [ ] semantic version for server/host
	- [ ] semantic version for plugin
5. Autopatching
	- https://github.com/docker/Notary
	- https://godoc.org/github.com/inconshreveable/go-update
