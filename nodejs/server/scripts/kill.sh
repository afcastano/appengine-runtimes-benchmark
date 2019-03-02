#!/usr/bin/env bash
server_path_portion="iecl-scheduler/server"
server_tcp_ports="8080,4000,9000"
kill_options="$@"

if [[ -n "${kill_options}" ]]; then
    echo "Using options for kill: ${kill_options}"
fi

echo "Killing npx server processes for this project ..."
for pid in $(ps | grep "npx server" | grep -v grep | awk '{print $1;}'); do
    # Double check that each "npx server" is actually for our project
    grep_result=$(lsof -p ${pid} | grep DIR | grep "${server_path_portion}")
    if [[ -n "${grep_result}" ]]; then
        echo "  - killing pid ${pid}"
        kill ${kill_options} ${pid}
    fi
done
echo "...done"
echo

echo "Killing processes listening to tcp ports ${server_tcp_ports} ..."
for pid in $(lsof -i tcp:${server_tcp_ports} -t); do
    echo "  - killing pid ${pid}"
    kill ${kill_options} ${pid}
done
echo "...done"

echo "Killing node processes executing gae-node-nestjs/server-start.js"
for pid in $(ps | grep node | grep -v grep | grep "gae-node-nestjs/server-start.js" | grep "${server_path_portion}" | awk '{print $1;}'); do
    echo "  - killing pid ${pid}"
    kill ${kill_options} ${pid}
done
echo "...done"
