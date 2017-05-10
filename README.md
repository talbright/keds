# Keds

A prototype for a generic CLI plugin framework.

# Running

		$ go run main.go

# Debugging

		$ open http://localhost:8081/debug/requests
		$ open http://localhost:8081/debug/events

# TODO

- Fixes
	- [x] fix inconsistent use of interfaces, concrete, and proto gen (plugins and events)
	- [x] better use of "addressing" for event source/target...maybe use URIs?
	- [x] prevent event source from sending back to itself (loop)
	- [ ] take advantage of cancelations via Context
	- [x] clean-up and create a "Base" plugin type that's easier to extend
	- [ ] rename event/event_bus to event/bus (go conventions)
	- [x] send config file location to plugins so they can read it (env is already inherited)
- Configuration management
	- [x] setup for viper
	- [x] config for plugin path
	- [ ] debug mode flag
- Plugin lifecycle management
	- [x] server loads plugins
	- [x] plugin registers itself with the server
	- [x] plugin exits when server exits
	- [x] capture plugin stdout/stderr and log to console
	- [x] plugin signal to terminate coordinates exit of server
- Misc
	- [ ] dependency management (glide or ?)
- Testing
	- [ ] add ci setup for github project
	- [x] add ginko test package
	- [ ] tests for client package
	- [ ] tests for plugin package
	- [ ] tests for server package
	- [x] tests for utils package
- Command line
	- [x] setup cobra
	- [x] plugin registration creates new cobra command
	- [x] plugin invocation when registered cobra command is passed in the args
- Versioning
	- [ ] semantic version for server/host
	- [ ] semantic version for plugin
- Autopatching
	- https://github.com/docker/Notary
	- https://godoc.org/github.com/inconshreveable/go-update
- Example Plugin
	- [x] signal termination event
	- [x] integration with cobra
- Notifications Plugin
	- [ ] create bare plugin to listen to events with filtering support
	- [ ] publish notifications to os/x (if generic go library available to handle os specifics use that)
	- [ ] publish notifications to slack channel
	see https://github.com/deckarep/gosx-notifier, https://github.com/0xAX/notificator
