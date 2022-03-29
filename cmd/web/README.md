# QEX Web Service
* this web app provides a HTTP Post function to receive a real-time test execution result
* this web app provides a HTTP form function to allow user upload a git master log

# Requirement
* Install Kafka server, and update Kafka zookeeper and bootstrap configuration in [config](../config/server.go)
* Install Influx DB, and update influx DB configuration in [config](../config/server.go)
* Install Grafana, and import config in [dashboard](../grafana)
* Update Jira and message web hook URL in [config](../config/server.go)
* Update the AT Sign Off and AT Found Bug JQL in [config](../config/server.go)

# Instruction
* `go build -o qex-web-server`
* `./qex-web-server`
