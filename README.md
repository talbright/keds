# Keds

A prototype for a generic and opinionated CLI plugin framework.

# Running

		$ go run main.go

# Debugging

		$ open http://localhost:8081/debug/requests
		$ open http://localhost:8081/debug/events

# TODO

1. Configuration management
	- [x] setup for viper
	- [x] config for plugin path
2. Plugin lifecycle management
	- [x] server loads plugins
	- [x] plugin registers itself with the server
	- [x] plugin exits when server exits
3. Misc
	- [ ] add ci setup for github project
	- [ ] dependency management (glide or...)
4. Testing framework
	- [x] add ginko test package
	- [ ] tests for client package
	- [ ] tests for plugin package
	- [ ] tests for server package
	- [x] tests for utils package
5. Command line
	- [x] setup cobra
	- [x] plugin registration creates new cobra command
	- [x] plugin invocation when registered cobra command is passed in the args
6. Versioning
	- [ ] semantic version for server/host
	- [ ] semantic version for plugin
7. Autopatching
	- https://github.com/docker/Notary
	- https://godoc.org/github.com/inconshreveable/go-update
