#!/bin/bash

if [ -z "$1" ]; then
  echo "Uso: $0 <numero_de_requisicoes> [host]"
  exit 1
fi

N=$1
DEFAULT_HOST="http://localhost:8080"
HOST=${2:-$DEFAULT_HOST}
URL="$HOST/api/v1/admin/orders?pageSize=5&page=2"

HEADERS="Accept: application/json"

make_request() {
  HTTP_CODE=$(curl -L -s -o /dev/null -w "%{http_code}" -H "$HEADERS" "$URL")

  if [ "$HTTP_CODE" -eq 200 ]; then
    echo "Request successful: HTTP $HTTP_CODE"
  else
    echo "Request failed: HTTP $HTTP_CODE"
  fi
}

export URL HEADERS
export -f make_request

seq $N | xargs -I {} -n1 -P$N bash -c 'make_request'
