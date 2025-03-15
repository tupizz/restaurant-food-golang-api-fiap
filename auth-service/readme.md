# Auth Service

This service provides a complete authentication and user management system built on AWS Cognito, Lambda, and API Gateway. It handles user registration, authentication, profile management, and more using CPF as the primary identifier.

## Features

- User registration and authentication
- CPF-based login
- Profile management
- Password reset and change
- Token refresh
- User deletion
- Secure token handling

## Architecture

- **AWS Cognito**: User directory and authentication provider
- **AWS Lambda**: Serverless backend for API handling
- **API Gateway**: HTTP API endpoints
- **Terraform**: Infrastructure as Code for AWS resources

## Prerequisites

- Node.js 18.x or higher
- AWS CLI configured with appropriate credentials
- Terraform CLI
- jq (for testing scripts)

## Environment Setup

Set up the required AWS credentials as environment variables or in your AWS CLI configuration:

```bash
# Configure AWS CLI
aws configure

# Or set environment variables
export AWS_ACCESS_KEY_ID="your-aws-access-key-id"
export AWS_SECRET_ACCESS_KEY="your-aws-secret-access-key"
export AWS_REGION="us-east-1"
```

For GitHub Actions, set these as repository secrets:

```bash
gh secret set AWS_ACCESS_KEY_ID --body "your-aws-access-key-id" -R your-username/your-repo
gh secret set AWS_SECRET_ACCESS_KEY --body "your-aws-secret-access-key" -R your-username/your-repo
gh secret set AWS_REGION --body "your-aws-region" -R your-username/your-repo
```

## Installation

Clone the repository and install dependencies:

```bash
git clone <repository-url>
cd auth-service
npm install
```

## Development

Build the TypeScript project:

```bash
npm run build
```

## Deployment

Deploy the service to AWS:

```bash
# Deploy with existing dependencies
npm run deploy

# Deploy and install dependencies first
npm run deploy:install
```

The deployment process:

1. Builds the TypeScript project
2. Generates Terraform route configurations
3. Initializes Terraform
4. Plans and applies Terraform changes
5. Outputs the API Gateway URL

## Testing

Test the deployed API endpoints:

```bash
# Set the API URL (replace with your actual API Gateway URL)
export API_URL="https://your-api-gateway-url.execute-api.us-east-1.amazonaws.com"

# Run the test script
./bin/test-routes.sh
```

The test script will:

1. Register a new user with a valid CPF
2. Test login with CPF
3. Get user profile
4. Update user profile
5. Test token refresh
6. Change password
7. Test logout
8. Delete the user

## Troubleshooting

### Invalid Refresh Token

If you encounter `NotAuthorizedException: Invalid Refresh Token` errors:

- The refresh token may have expired
- The token might have been invalidated by a password change or logout
- Re-authenticate to get a new token

### Finding User by CPF

Use the `findCPF.ts` script to check if a user exists:

```bash
# Run with a specific CPF
npx ts-node ./bin/findCPF.ts 12345678901

# Run with default CPF
npx ts-node ./bin/findCPF.ts
```

## API Endpoints

| Method | Endpoint              | Description                     |
| ------ | --------------------- | ------------------------------- |
| POST   | /auth/register        | Register a new user             |
| POST   | /auth/login           | Login with username/password    |
| POST   | /auth/login/cpf       | Login with CPF/password         |
| GET    | /auth/profile         | Get authenticated user profile  |
| GET    | /auth/user/cpf/{cpf}  | Get user by CPF                 |
| PUT    | /auth/profile         | Update user profile             |
| PUT    | /auth/user/cpf        | Update user CPF                 |
| DELETE | /auth/user/{username} | Delete user                     |
| POST   | /auth/token           | Refresh authentication token    |
| POST   | /auth/forgot-password | Initiate password reset         |
| POST   | /auth/reset-password  | Complete password reset         |
| POST   | /auth/change-password | Change password (authenticated) |
| POST   | /auth/verify          | Verify user attribute           |
| POST   | /auth/logout          | Logout user                     |

# Utilities

## Setting Secrets to GH repo

```sh
# Set AWS_ACCESS_KEY_ID secret for your repository

gh secret set AWS_ACCESS_KEY_ID --body "your-aws-access-key-id" -R tupizz/restaurant-food-golang-api-fiap

# Set AWS_SECRET_ACCESS_KEY secret for your repository

gh secret set AWS_SECRET_ACCESS_KEY --body "your-aws-secret-access-key" -R tupizz/restaurant-food-golang-api-fiap

# Set AWS_REGION secret for your repository

gh secret set AWS_REGION --body "your-aws-region" -R tupizz/restaurant-food-golang-api-fiap
```
