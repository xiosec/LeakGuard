#!/bin/bash

ELASTICSEARCH_URL="http://elasticsearch:9200"
INDEX_NAME="pwned-password"

check_index_exists() {
  response=$(curl -s -o /dev/null -w "%{http_code}" "$ELASTICSEARCH_URL/$INDEX_NAME")
  [[ $response == "200" ]]
}

wait_for_index() {
  while ! check_index_exists; do
    echo "Waiting for Elasticsearch index '$INDEX_NAME' to be created..."
    sleep 5
  done

  echo "Elasticsearch index '$INDEX_NAME' is now available."
}

wait_for_index
