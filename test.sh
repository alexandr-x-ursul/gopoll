#!/usr/local/bin/bash
#curl -k -X "POST" -d '{"question":"Are you happy?", "answers":["yes but no", "no"]}' localhost:8080/api/poll

for x in $(seq 105);do
  curl -k -X "POST" -d '{"answer":"Damn right she is!"}' localhost:8080/api/poll/3 &
done

#curl -k -X "DELETE" localhost:8080/api/poll/2
