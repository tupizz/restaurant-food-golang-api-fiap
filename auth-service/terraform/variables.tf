# AWS region for deploying resources
# This defines where all AWS resources will be created
variable "aws_region" {
  description = "The AWS region to deploy resources"
  type        = string
  default     = "us-east-1" # Default to US East (N. Virginia)
}

# Environment name (dev, staging, prod)
# Used for tagging and naming resources appropriately
variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "prod" # Default to production
}

# Lambda function name
# The name of the Lambda function that will handle authentication
variable "lambda_name" {
  description = "Name of the Lambda function"
  type        = string
  default     = "auth-service-lambda" # Default name
}

# Lambda runtime environment
# The programming language runtime for the Lambda function
variable "lambda_runtime" {
  description = "Runtime for the Lambda function"
  type        = string
  default     = "nodejs18.x" # Node.js 18.x runtime
}

# Lambda memory allocation
# The amount of memory allocated to the Lambda function
variable "lambda_memory" {
  description = "Memory allocation for the Lambda function in MB"
  type        = number
  default     = 256 # 256 MB of memory
}

# Lambda timeout
# Maximum execution time for the Lambda function
variable "lambda_timeout" {
  description = "Timeout for the Lambda function in seconds"
  type        = number
  default     = 30 # 30 seconds timeout
}

# Log retention period
# Number of days to keep CloudWatch logs
variable "log_retention_days" {
  description = "Number of days to retain logs"
  type        = number
  default     = 30 # Keep logs for 30 days
}
