#!/bin/bash

# check if the number of requests is provided
if [ -z "$1" ]; then
  echo "Usage: $0 <number_of_requests>"
  exit 1
fi

# Number of parallel requests
N=$1

# URL of your local server
URL="http://localhost:8080/api/v1/admin/orders?pageSize=5&page=2"

# Optional headers or additional options can be added here (e.g., Authorization)
HEADERS="Accept: application/json"

# Function to make a single request, handling potential errors
make_request() {
  HTTP_CODE=$(curl -L -s -o /dev/null -w "%{http_code}" -H "$HEADERS" "$URL")

  # Check if HTTP code is valid
  if [ "$HTTP_CODE" -eq 200 ]; then
    echo "Request successful: HTTP $HTTP_CODE"
  else
    echo "Request failed: HTTP $HTTP_CODE"
  fi
}

# Export variables and functions to make them available in parallel subshells
export URL HEADERS
export -f make_request

# Generate N requests in parallel
seq $N | xargs -I {} -n1 -P$N bash -c 'make_request'
