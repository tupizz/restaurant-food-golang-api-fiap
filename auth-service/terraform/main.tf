provider "aws" {
  region = var.aws_region
}

# Create an archive of the built Lambda code.
data "archive_file" "lambda_zip" {
  type        = "zip"
  source_dir  = "../dist"
  output_path = "${path.module}/lambda.zip"
}

resource "random_string" "suffix" {
  length  = 6
  special = false
  upper   = false
}

resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role_${random_string.suffix.result}"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "lambda.amazonaws.com"
      }
    }]
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

locals {
  tags = {
    Project     = "AuthService"
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}

resource "aws_lambda_function" "my_lambda" {
  function_name = var.lambda_name
  role          = aws_iam_role.lambda_exec.arn
  handler       = "index.handler"
  runtime       = var.lambda_runtime

  # Points to the zip created from the Lambda build.
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path)

  # Add environment variables
  environment {
    variables = {
      NODE_ENV = "production"
      # Add other environment variables your Lambda needs
      # LOG_LEVEL = "info"
      # API_VERSION = "v1"
    }
  }

  # Configure Lambda performance
  memory_size = var.lambda_memory
  timeout     = var.lambda_timeout
  tags        = local.tags
}

# API Gateway
resource "aws_apigatewayv2_api" "lambda_api" {
  name          = "auth-service-api"
  protocol_type = "HTTP"

  # Add CORS configuration
  cors_configuration {
    allow_origins = ["https://your-frontend-domain.com"] # Adjust as needed
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
    allow_headers = ["Content-Type", "Authorization"]
    max_age       = 300
  }
  tags = local.tags
}

resource "aws_apigatewayv2_stage" "lambda_stage" {
  api_id      = aws_apigatewayv2_api.lambda_api.id
  name        = "$default"
  auto_deploy = true

  # Add access logging
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_logs.arn
    format = jsonencode({
      requestId          = "$context.requestId"
      ip                 = "$context.identity.sourceIp"
      requestTime        = "$context.requestTime"
      httpMethod         = "$context.httpMethod"
      routeKey           = "$context.routeKey"
      status             = "$context.status"
      protocol           = "$context.protocol"
      responseLength     = "$context.responseLength"
      path               = "$context.path"
      integrationLatency = "$context.integrationLatency"
      responseLatency    = "$context.responseLatency"
    })
  }
}

resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.lambda_api.id
  integration_type   = "AWS_PROXY"
  integration_method = "POST"
  integration_uri    = aws_lambda_function.my_lambda.invoke_arn
}

# The routes are now defined in routes.tf (generated file)

# Lambda permission to allow API Gateway to invoke the function
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.my_lambda.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.lambda_api.execution_arn}/*/*/*"
}

# Output the API Gateway URL
output "api_gateway_url" {
  value = aws_apigatewayv2_stage.lambda_stage.invoke_url
}

# Temporarily comment this out
# resource "aws_cloudwatch_log_group" "lambda_logs" {
#   name              = "/aws/lambda/${aws_lambda_function.my_lambda.function_name}"
#   retention_in_days = 30
#   
#   lifecycle {
#     prevent_destroy = false
#     ignore_changes = [
#       tags,
#     ]
#   }
# }

# Add CloudWatch Logs permissions to the Lambda role
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging_policy_${random_string.suffix.result}"
  description = "IAM policy for logging from a lambda"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      Resource = "arn:aws:logs:*:*:*",
      Effect   = "Allow"
    }]
  })

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

# Configure API Gateway logging
resource "aws_cloudwatch_log_group" "api_gateway_logs" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.lambda_api.name}_${random_string.suffix.result}"
  retention_in_days = 30

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.my_lambda.function_name}"
  retention_in_days = 30

  lifecycle {
    prevent_destroy = false
    ignore_changes = [
      tags,
    ]
  }
}
