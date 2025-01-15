#!/bin/bash

DEFAULT_HOST="http://localhost:8080"
HOST=${1:-$DEFAULT_HOST}
URL="$HOST/api/v1/admin/products"

HEADERS="Content-Type: application/json"

PRODUCTS=(
  '{"name": "X-Burguer", "price": 10.00, "description": "Hambúrguer de carne bovina", "category": "lanche", "images": ["https://placehold.co/600x400/png"]}'
  '{"name": "Coca-Cola", "price": 5.00, "description": "Refrigerante de cola", "category": "bebida", "images": ["https://placehold.co/600x400/png"]}'
  '{"name": "Pudim", "price": 8.00, "description": "Sobremesa de pudim", "category": "sobremesa", "images": ["https://placehold.co/600x400/png"]}'
  '{"name": "Batata frita", "price": 3.00, "description": "Acompanhamento de batata frita", "category": "acompanhamento", "images": ["https://placehold.co/600x400/png"]}'
)

make_request() {
  local product_data=$1

  HTTP_CODE=$(curl -L -s -o /dev/null -w "%{http_code}" -H "$HEADERS" -d "$product_data" "$URL")

  if [ "$HTTP_CODE" -eq 201 ]; then
    echo "Produto cadastrado com sucesso: HTTP $HTTP_CODE"
  else
    echo "Falha ao cadastrar produto: HTTP $HTTP_CODE"
  fi
}

for product in "${PRODUCTS[@]}"; do
  make_request "$product"
done
