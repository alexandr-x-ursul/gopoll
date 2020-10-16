#!/usr/local/bin/bash
curl -k -X "POST" -d '{"question":"Are you happy?", "answers":["yes but no", "no"]}' localhost:8080/api/poll

for x in $(seq 97);do
  curl -k -X "POST" -d '{"id":1, "answer":"yes but no"}' localhost:8080/api/poll/1 &
done

#curl -k -X "DELETE" localhost:8080/api/poll/2
