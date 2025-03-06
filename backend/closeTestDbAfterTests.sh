#!/bin/sh
# This script is for stoping test database after running tests
go test ./tests/...

TEST_EXIT_CODE=$?

curl --unix-socket /var/run/docker.sock -X POST http://localhost/containers/swiftcodeapi-test_db-1/stop

exit $TEST_EXIT_CODE
