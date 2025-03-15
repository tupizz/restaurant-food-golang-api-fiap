# Configure the AWS provider with the specified region
provider "aws" {
  region = var.aws_region # Use the region defined in variables.tf
}

# Configure Terraform backend to store state in S3
# This allows team collaboration and state persistence
terraform {
  backend "s3" {
    bucket = "fiap-tf-state-bucket"           # S3 bucket to store Terraform state
    key    = "auth-service/terraform.tfstate" # Path within the bucket
    region = "us-east-1"                      # Region where the S3 bucket is located
  }

  # Define required provider versions for compatibility
  required_providers {
    aws = {
      source  = "hashicorp/aws" # AWS provider source
      version = "~> 4.0"        # Compatible with version 4.x
    }
    random = {
      source  = "hashicorp/random" # Random provider for generating unique values
      version = "~> 3.0"           # Compatible with version 3.x
    }
  }
}

# Create a ZIP archive of the Lambda code for deployment
# This packages the compiled code in the dist directory
data "archive_file" "lambda_zip" {
  type        = "zip"                       # Archive type
  source_dir  = "../dist"                   # Directory containing compiled code
  output_path = "${path.module}/lambda.zip" # Where to save the ZIP file
}

# Generate a random string to append to resource names
# This ensures unique resource names across deployments
resource "random_string" "suffix" {
  length  = 6     # 6 characters long
  special = false # No special characters
  upper   = false # No uppercase letters
}

# Create an IAM role for the Lambda function
# This role defines what AWS services the Lambda can access
resource "aws_iam_role" "lambda_exec" {
  name = "lambda_exec_role_${random_string.suffix.result}" # Unique role name

  # Trust policy allowing Lambda service to assume this role
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

  # Create new role before destroying old one to avoid downtime
  lifecycle {
    create_before_destroy = true
  }
}

# Attach the basic Lambda execution policy to the role
# This allows Lambda to create logs in CloudWatch
resource "aws_iam_role_policy_attachment" "lambda_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

# Define common tags for resources
# Tags help with resource organization and cost tracking
locals {
  tags = {
    Project     = "AuthService"
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}

# Create the Lambda function resource
# This is the serverless function that will handle authentication requests
resource "aws_lambda_function" "my_lambda" {
  function_name = var.lambda_name              # Name from variables.tf
  role          = aws_iam_role.lambda_exec.arn # IAM role for permissions
  handler       = "index.handler"              # Entry point in the code
  runtime       = var.lambda_runtime           # Runtime environment (Node.js version)

  # Reference to the deployment package
  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = filebase64sha256(data.archive_file.lambda_zip.output_path) # For detecting changes

  # Environment variables available to the Lambda function
  environment {
    variables = {
      NODE_ENV             = "production"                               # Runtime environment
      COGNITO_USER_POOL_ID = aws_cognito_user_pool.user_pool.id         # Reference to Cognito User Pool
      COGNITO_CLIENT_ID    = aws_cognito_user_pool_client.app_client.id # Reference to Cognito App Client
    }
  }

  # Performance configuration
  memory_size = var.lambda_memory  # RAM allocation
  timeout     = var.lambda_timeout # Maximum execution time
  tags        = local.tags         # Resource tags
}

# Create an API Gateway to expose the Lambda function as a REST API
resource "aws_apigatewayv2_api" "lambda_api" {
  name          = "auth-service-api" # API name
  protocol_type = "HTTP"             # HTTP API (vs WebSocket)

  # Configure Cross-Origin Resource Sharing (CORS)
  # This allows web browsers to make requests to your API
  cors_configuration {
    allow_origins = ["https://your-frontend-domain.com"]        # Allowed origins
    allow_methods = ["GET", "POST", "PUT", "DELETE", "OPTIONS"] # Allowed HTTP methods
    allow_headers = ["Content-Type", "Authorization"]           # Allowed headers
    max_age       = 300                                         # Browser can cache CORS response for 300 seconds
  }
  tags = local.tags
}

# Create a stage for the API Gateway
# A stage is a named reference to a deployment of the API
resource "aws_apigatewayv2_stage" "lambda_stage" {
  api_id      = aws_apigatewayv2_api.lambda_api.id
  name        = "$default" # Default stage
  auto_deploy = true       # Automatically deploy changes

  # Configure access logging for API requests
  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gateway_logs.arn # Log destination
    format = jsonencode({                                           # Log format with request details
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

# Create an integration between API Gateway and Lambda
# This defines how API Gateway forwards requests to Lambda
resource "aws_apigatewayv2_integration" "lambda_integration" {
  api_id             = aws_apigatewayv2_api.lambda_api.id
  integration_type   = "AWS_PROXY"                              # Lambda proxy integration
  integration_method = "POST"                                   # Method used to invoke Lambda
  integration_uri    = aws_lambda_function.my_lambda.invoke_arn # Lambda invocation URI
}

# Note: Routes are defined in routes.tf (generated file)

# Grant API Gateway permission to invoke the Lambda function
# This is required for the API Gateway to call the Lambda
resource "aws_lambda_permission" "api_gateway" {
  statement_id  = "AllowExecutionFromAPIGateway" # Permission identifier
  action        = "lambda:InvokeFunction"        # Permission to invoke Lambda
  function_name = aws_lambda_function.my_lambda.function_name
  principal     = "apigateway.amazonaws.com"                               # Service granted permission
  source_arn    = "${aws_apigatewayv2_api.lambda_api.execution_arn}/*/*/*" # API Gateway ARN pattern
}

# Output the API Gateway URL for reference
output "api_gateway_url" {
  value = aws_apigatewayv2_stage.lambda_stage.invoke_url # URL to invoke the API
}

# Create a custom IAM policy for Lambda logging
# This allows more detailed control over logging permissions
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging_policy_${random_string.suffix.result}"
  description = "IAM policy for logging from a lambda"

  # Policy document defining permissions
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      Resource = "arn:aws:logs:*:*:*", # Access to all log groups
      Effect   = "Allow"
    }]
  })

  lifecycle {
    create_before_destroy = true
  }
}

# Attach the logging policy to the Lambda execution role
resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

# Create a CloudWatch Log Group for API Gateway logs
resource "aws_cloudwatch_log_group" "api_gateway_logs" {
  name              = "/aws/apigateway/${aws_apigatewayv2_api.lambda_api.name}_${random_string.suffix.result}"
  retention_in_days = 30 # Keep logs for 30 days

  lifecycle {
    create_before_destroy = true
  }
}

# Create a CloudWatch Log Group for Lambda logs
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${aws_lambda_function.my_lambda.function_name}"
  retention_in_days = 30 # Keep logs for 30 days

  lifecycle {
    prevent_destroy = false # Allow destruction of log group
    ignore_changes = [
      tags, # Ignore changes to tags
    ]
  }
}

# Attach the Cognito policy to the Lambda execution role
# This allows Lambda to interact with Cognito for authentication
resource "aws_iam_role_policy_attachment" "lambda_cognito_policy" {
  role       = aws_iam_role.lambda_exec.name
  policy_arn = aws_iam_policy.lambda_cognito_policy.arn # Policy defined in cognito.tf
}
