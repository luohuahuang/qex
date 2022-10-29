# Instruction
* this kafka consumer will try to consume the msgs from QEX web server, and try to parse Jira ticket IDs from branch
  names, commit msgs, comments, and MR titles
* and update the Jira tickets and seatalk system accounts with Jenkins build/deploy status and test results

# Instruction
* `go build -o jenkins-build-consumer`
* `./jenkins-build-consumer`