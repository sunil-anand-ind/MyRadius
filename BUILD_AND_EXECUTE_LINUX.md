## LINUX BASED EXECUTION 

## Steps to Build the solution

PRE-REQUISITE: 
 - Install latest go package.

### Radius Control Plane

Move to the folder "radius-control-plane"

```
go mod init radius-control-plane
go mod tidy
go build .
```

The executable named "radius-control-plane" should be generated.


### Redis Control Plane Logger

Move to the folder "redis-control-plane-logger"

```
go mod init redis-control-plane-logger
go mod tidy
go build .
```

The executable named "redis-control-plane-logger" should be generated.

## Steps to Run the solution **

PRE-REQUISITE: 
 - A Linux based platform for execution
 - Install the latest go package
 - Install the radclient for traffic generation
 - Install redis and configure ```notify-keyspace-events "KEA"``` in the "/etc/redis.conf" 
 
### STEP 1 : Open a new terminal window and execute the following. 
This should start the Radius Server on default port 1813 with default user name "testuser" and default secret "testing123".
This should start the Redis Client and connect to the default port 6379 in localhost.

```
cd radius-control-plane
./radius-control-plane
```

### STEP 2 : Open a new terminal window and execute the following.
This should start the Redis Client and connect to the default port 6379 in localhost.

```
cd redis-control-plane-logger
./redis-control-plane-logger
```
### STEP 3: Open a new terminal window and execute the following.
This should trigger Accounting-Request message through the radclient traffic simulator.

```
cd simulator
radclient -x localhost:1813 acct testing123 < acct_start.txt
```

**RESULT:** 
The traffic generator should send the Accounting-Request towards radius-control-plane.
The radius-control-plane verifies the user name and secret, parses the message and send Accounting-Response.
The radius-control-plane prepares the Accounting record in the form of JSON and stores it in the REDIS persistent storage.
The redis-control-plane-logger subscribes for the keyspace-events and retrieve the records.
The redis-control-plane-logger logs the message with the appropriate key and timestamp in the log file "/var/log/radius_updates.log"
