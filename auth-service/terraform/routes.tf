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

