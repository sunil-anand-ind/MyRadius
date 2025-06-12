This is a custom Radius server implementation of Accounting Server which stores the Accounting Information in the form of record into Redis persistent storage.

## radius-control-plane

* Responsible for handling Radius Accounting Information (START, STOP, INTERIM-UPDATE).
* Connects to the Redis Server for storing Accounting Records.
* Radius Server Listens on Port 1813 Port. 
* Parses Accounting information and performs authentication based on username and secret.
* Marshals the Accounting Information in the form of JSON data.
* Stores the JSON data with a expected key format radius:acct:{username}:{acct_session_id}:{timestamp}

## redis-control-plane-logger

* Responsible for subscribing to the keyspace notifications onto redis.
* Watch the key pattern which contains radius:acct:*
* Displays the information retrieved from the key patters with timestamp onto custom logger file.

## simulator
* Uses opensource free-radius package to trigger Accounting traffic towards radius control plane

This project is implemented in GO using the opensource radius and redis packages.
This project is containerized and can be able to deploy using a single docker-compose.yml.
