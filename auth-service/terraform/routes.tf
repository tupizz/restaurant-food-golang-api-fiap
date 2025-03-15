# This file is generated automatically. Do not edit directly.

resource "aws_apigatewayv2_route" "route_0" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/register"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_1" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "GET /"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_2" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/login"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_3" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/login/cpf"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_4" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "GET /auth/user/cpf/{cpf}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_5" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "GET /auth/profile"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_6" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "PUT /auth/profile"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_7" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "PUT /auth/user/cpf"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_8" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "DELETE /auth/user/{username}"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_9" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/token"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_10" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/forgot-password"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_11" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/reset-password"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_12" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/change-password"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_13" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/verify"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

resource "aws_apigatewayv2_route" "route_14" {
  api_id    = aws_apigatewayv2_api.lambda_api.id
  route_key = "POST /auth/logout"
  target    = "integrations/${aws_apigatewayv2_integration.lambda_integration.id}"
}

