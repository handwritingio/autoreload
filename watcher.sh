#!/bin/bash

# Watches for .go file changes in the project directory, telling the server to
# restart whenever a change is detected. Should be run outside the docker
# container, because file system events are not properly detected inside the
# boot2docker VM: https://github.com/boot2docker/boot2docker/issues/688
#
# Requires fswatch which can be installed with brew.
#
# This works by opening a tcp connection on $RESET_PORT and sending a few bytes
# (the actual message doesn't matter). If the API was started with
# docker-compose, the autoreloader will be running inside the container and
# waiting on $RESET_PORT for a connection.

set -e

# get_ip will get your docker IP address
# depending on if you run boot2docker or docker-machine
# TODO: what if you're running docker-machine and using
# a machine name other than default?
get_ip() {
    boot2docker ip 2> /dev/null
    if [ $? -eq 0 ]; then
        echo $(boot2docker ip)
    else
        echo $(docker-machine ip default)
    fi
}


RESET_PORT=12345
if [ -n $RESET_PORT ]; then
    RESET_PORT=$1
fi

wait_for_changes() {
    echo 'Waiting for changes'
    fswatch -1 --include '\.go$' --exclude '.*' .
}

reload_server() {
    echo 'Reloading server'
    echo 'reset' | nc -w 1 $(get_ip) $RESET_PORT
}

if [[ $1 = 'reload' ]]; then
    reload_server
    exit 0
fi

# Exit on ctrl-c (without this, ctrl-c would go to fswatch, causing it to
# reload instead of exit):
trap 'exit 0' SIGINT

while true; do
    wait_for_changes
    reload_server
done
