#!/usr/bin/bash

PORT=3000
ADDRESS=localhost:${PORT}

# Start by deleting all existing todos on the server
curl -iL -w "\n" -X DELETE ${ADDRESS}/todos/deleteAll

# Add some todos
curl -iL -X POST -H "Content-Type: application/json" -d '{"description":"Write Code","tags":["coding"],"due":"2023-01-03T15:04:05+00:00"}' ${SERVERADDR}/todo/

# Get todos by tag
curl -iL -w "\n" ${ADDRESS}/tag/coding/

# Get todos by due
curl -iL -w "\n" ${ADDRESS}/due/2023/01/02