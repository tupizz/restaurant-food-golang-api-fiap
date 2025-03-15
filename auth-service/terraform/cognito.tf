# Generate a random string for Cognito resource names
# This ensures unique resource names across deployments
resource "random_string" "cognito_suffix" {
  length  = 6     # 6 characters long
  special = false # No special characters
  upper   = false # No uppercase letters
}

# Create an AWS Cognito User Pool
# This is the user directory for authentication and user management
resource "aws_cognito_user_pool" "user_pool" {
  name = "auth-user-pool-${random_string.cognito_suffix.result}" # Unique pool name

  # Change from alias_attributes to username_attributes
  # This specifies that username is required and not an alias
  # We'll use CPF as the username directly
  username_attributes = [] # Empty means username is required

  # Configure username case sensitivity
  username_configuration {
    case_sensitive = false # Usernames (CPFs) are not case sensitive
  }

  # Define password policy requirements
  password_policy {
    minimum_length                   = 8    # Minimum 8 characters
    require_lowercase                = true # Must include lowercase
    require_numbers                  = true # Must include numbers
    require_symbols                  = true # Must include symbols
    require_uppercase                = true # Must include uppercase
    temporary_password_validity_days = 7    # Temporary passwords valid for 7 days
  }

  # Configure Multi-Factor Authentication (MFA)
  mfa_configuration = "OPTIONAL" # MFA is optional
  software_token_mfa_configuration {
    enabled = true # Enable TOTP (Time-based One-Time Password)
  }

  # Configure account recovery options
  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_phone_number" # Recover via phone
      priority = 1                       # First priority
    }
    recovery_mechanism {
      name     = "verified_email" # Recover via email
      priority = 2                # Second priority
    }
  }

  # Configure email sending
  email_configuration {
    email_sending_account = "COGNITO_DEFAULT" # Use Cognito's email service
  }

  # Configure SMS sending for verification and MFA
  sms_configuration {
    external_id    = "auth-service-sms-${random_string.cognito_suffix.result}"
    sns_caller_arn = aws_iam_role.cognito_sms_role.arn # IAM role for SMS
  }

  # Configure attributes that are automatically verified
  auto_verified_attributes = ["phone_number"] # Auto-verify phone numbers

  # Keep CPF as a custom attribute for validation and additional info
  # But we'll also use it as the username
  schema {
    name                     = "cpf"    # Attribute name
    attribute_data_type      = "String" # Data type
    developer_only_attribute = false    # Visible to app
    mutable                  = false    # Cannot be changed
    required                 = false    # Not required (Cognito limitation)
    string_attribute_constraints {
      min_length = 11 # CPF is 11 digits
      max_length = 11
    }
  }

  # Define name attribute
  schema {
    name                     = "name"
    attribute_data_type      = "String"
    developer_only_attribute = false
    mutable                  = true # Can be changed
    required                 = true # Required attribute
    string_attribute_constraints {
      min_length = 1
      max_length = 100
    }
  }

  # Configure advanced security features
  user_pool_add_ons {
    advanced_security_mode = "AUDIT" # Audit mode for security features
  }

  # Configure device tracking
  device_configuration {
    challenge_required_on_new_device      = true # Require verification for new devices
    device_only_remembered_on_user_prompt = true # Remember device only if user agrees
  }

  # Lambda triggers for custom behavior (commented out)
  # lambda_config {
  #   pre_sign_up = aws_lambda_function.pre_sign_up.arn
  #   custom_message = aws_lambda_function.custom_message.arn
  # }

  # Resource tags
  tags = {
    Name        = "AuthUserPool"
    Environment = "Production"
    ManagedBy   = "Terraform"
  }
}

# Create a Cognito User Pool Client
# This is the app that will interact with the User Pool
resource "aws_cognito_user_pool_client" "app_client" {
  name         = "auth-app-client-${random_string.cognito_suffix.result}" # Client name
  user_pool_id = aws_cognito_user_pool.user_pool.id                       # Reference to User Pool

  # Client secret configuration
  generate_secret = false # No client secret (for public clients)

  # Token validity periods
  refresh_token_validity = 30 # Refresh tokens valid for 30 days
  access_token_validity  = 1  # Access tokens valid for 1 hour
  id_token_validity      = 1  # ID tokens valid for 1 hour

  # Token validity units
  token_validity_units {
    access_token  = "hours"
    id_token      = "hours"
    refresh_token = "days"
  }

  # OAuth 2.0 flow configuration
  allowed_oauth_flows                  = ["implicit", "code"]                    # Allowed flows
  allowed_oauth_flows_user_pool_client = true                                    # Enable OAuth flows
  allowed_oauth_scopes                 = ["phone", "email", "openid", "profile"] # Allowed scopes

  # Callback and logout URLs for OAuth flows
  callback_urls = ["https://example.com/callback", "http://localhost:3000/callback"]
  logout_urls   = ["https://example.com/logout", "http://localhost:3000/logout"]

  # Identity providers
  supported_identity_providers = ["COGNITO"] # Only Cognito (no external providers)

  # Authentication flows
  explicit_auth_flows = [
    "ALLOW_ADMIN_USER_PASSWORD_AUTH", # Allow admin auth API
    "ALLOW_USER_PASSWORD_AUTH",       # Allow user password auth
    "ALLOW_REFRESH_TOKEN_AUTH",       # Allow refresh token usage
    "ALLOW_USER_SRP_AUTH",            # Allow Secure Remote Password
    "ALLOW_CUSTOM_AUTH"               # Allow custom auth flows
  ]

  # Security configuration
  prevent_user_existence_errors = "ENABLED" # Don't reveal if user exists

  # Attributes the client can read
  read_attributes = [
    "email",
    "email_verified",
    "phone_number",
    "phone_number_verified",
    "custom:cpf",
    "name"
  ]

  # Attributes the client can write
  write_attributes = [
    "email",
    "phone_number",
    "custom:cpf",
    "name"
  ]
}

# Create a custom domain for the Cognito User Pool
# This provides a branded URL for the hosted UI
resource "aws_cognito_user_pool_domain" "main" {
  domain       = "auth-service-${random_string.cognito_suffix.result}" # Domain prefix
  user_pool_id = aws_cognito_user_pool.user_pool.id
}

# Create an IAM role for Cognito SMS sending
# This allows Cognito to send SMS messages via Amazon SNS
resource "aws_iam_role" "cognito_sms_role" {
  name = "cognito-sms-role-${random_string.cognito_suffix.result}"

  # Trust policy allowing Cognito to assume this role
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = "sts:AssumeRole",
      Effect = "Allow",
      Principal = {
        Service = "cognito-idp.amazonaws.com"
      }
    }]
  })

  tags = {
    Name      = "CognitoSMSRole"
    ManagedBy = "Terraform"
  }
}

# Create an IAM policy for SMS sending permissions
resource "aws_iam_policy" "cognito_sms_policy" {
  name        = "cognito-sms-policy-${random_string.cognito_suffix.result}"
  description = "Policy for Cognito to send SMS via SNS"

  # Policy allowing SNS publish action
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Action = [
        "sns:publish" # Permission to publish SMS
      ],
      Resource = "*" # To any SNS topic
    }]
  })
}

# Attach the SMS policy to the SMS role
resource "aws_iam_role_policy_attachment" "cognito_sms_attach" {
  role       = aws_iam_role.cognito_sms_role.name
  policy_arn = aws_iam_policy.cognito_sms_policy.arn
}

# Create a Cognito Resource Server for OAuth 2.0 scopes
# This defines custom scopes for API access control
resource "aws_cognito_resource_server" "resource_server" {
  identifier   = "https://api.auth-service.com" # Resource server identifier
  name         = "Auth Service API"             # Resource server name
  user_pool_id = aws_cognito_user_pool.user_pool.id

  # Define custom scopes
  scope {
    scope_name        = "read" # Read scope
    scope_description = "Read access to API"
  }

  scope {
    scope_name        = "write" # Write scope
    scope_description = "Write access to API"
  }

  scope {
    scope_name        = "admin" # Admin scope
    scope_description = "Admin access to API"
  }
}

# Create a Cognito Identity Pool
# This allows authenticated users to access AWS services directly
resource "aws_cognito_identity_pool" "identity_pool" {
  identity_pool_name               = "auth_identity_pool_${random_string.cognito_suffix.result}"
  allow_unauthenticated_identities = false # No guest access
  allow_classic_flow               = false # Use enhanced flow

  # Configure Cognito User Pool as an identity provider
  cognito_identity_providers {
    client_id               = aws_cognito_user_pool_client.app_client.id
    provider_name           = "cognito-idp.${var.aws_region}.amazonaws.com/${aws_cognito_user_pool.user_pool.id}"
    server_side_token_check = false # No additional token validation
  }

  tags = {
    Name      = "AuthIdentityPool"
    ManagedBy = "Terraform"
  }
}

# Create an IAM role for authenticated users
# This defines what AWS services authenticated users can access
resource "aws_iam_role" "authenticated" {
  name = "cognito_authenticated_role_${random_string.cognito_suffix.result}"

  # Trust policy allowing authenticated users to assume this role
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Effect = "Allow",
      Principal = {
        Federated = "cognito-identity.amazonaws.com"
      },
      Action = "sts:AssumeRoleWithWebIdentity",
      Condition = {
        StringEquals = {
          "cognito-identity.amazonaws.com:aud" = aws_cognito_identity_pool.identity_pool.id
        },
        "ForAnyValue:StringLike" = {
          "cognito-identity.amazonaws.com:amr" = "authenticated"
        }
      }
    }]
  })

  tags = {
    Name      = "CognitoAuthenticatedRole"
    ManagedBy = "Terraform"
  }
}

# Create an IAM policy for authenticated users
# This defines the permissions for authenticated users
resource "aws_iam_policy" "authenticated_policy" {
  name        = "cognito_authenticated_policy_${random_string.cognito_suffix.result}"
  description = "Policy for Cognito authenticated users"

  # Policy document with permissions
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "mobileanalytics:PutEvents", # Mobile analytics
          "cognito-sync:*"             # Cognito Sync
        ],
        Resource = "*"
      },
      {
        Effect = "Allow",
        Action = [
          "execute-api:Invoke" # Invoke API Gateway
        ],
        Resource = "arn:aws:execute-api:${var.aws_region}:*:*/*/GET/*" # GET methods only
      }
    ]
  })
}

# Attach the authenticated policy to the authenticated role
resource "aws_iam_role_policy_attachment" "authenticated_attach" {
  role       = aws_iam_role.authenticated.name
  policy_arn = aws_iam_policy.authenticated_policy.arn
}

# Attach roles to the Identity Pool
# This maps the authenticated role to the identity pool
resource "aws_cognito_identity_pool_roles_attachment" "identity_pool_roles" {
  identity_pool_id = aws_cognito_identity_pool.identity_pool.id

  # Role mapping
  roles = {
    "authenticated" = aws_iam_role.authenticated.arn
  }
}

# Create a Cognito User Group for admins
# This allows role-based access control
resource "aws_cognito_user_group" "admin_group" {
  name         = "admin" # Group name
  user_pool_id = aws_cognito_user_pool.user_pool.id
  description  = "Admin group with elevated privileges"
  precedence   = 1 # Higher precedence (lower number)
}

# Create a Cognito User Group for regular users
resource "aws_cognito_user_group" "user_group" {
  name         = "users" # Group name
  user_pool_id = aws_cognito_user_pool.user_pool.id
  description  = "Regular users group"
  precedence   = 10 # Lower precedence (higher number)
}

# Lambda permissions for Cognito triggers (commented out)
# resource "aws_lambda_permission" "cognito_pre_signup" {
#   statement_id  = "AllowCognitoPreSignUp"
#   action        = "lambda:InvokeFunction"
#   function_name = aws_lambda_function.pre_sign_up.function_name
#   principal     = "cognito-idp.amazonaws.com"
#   source_arn    = aws_cognito_user_pool.user_pool.arn
# }

# Create an IAM policy for Lambda to interact with Cognito
# This allows the Lambda function to call Cognito APIs
resource "aws_iam_policy" "lambda_cognito_policy" {
  name        = "lambda-cognito-policy-${random_string.cognito_suffix.result}"
  description = "IAM policy for Lambda to interact with Cognito"

  # Policy document with Cognito permissions
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [{
      Action = [
        "cognito-idp:AdminCreateUser",             # Create users
        "cognito-idp:AdminGetUser",                # Get user details
        "cognito-idp:AdminInitiateAuth",           # Authenticate users
        "cognito-idp:AdminUpdateUserAttributes",   # Update user attributes
        "cognito-idp:AdminDeleteUser",             # Delete users
        "cognito-idp:AdminRespondToAuthChallenge", # Respond to auth challenges
        "cognito-idp:ListUsers",                   # List users
        "cognito-idp:AdminAddUserToGroup",         # Add users to groups
        "cognito-idp:AdminRemoveUserFromGroup",    # Remove users from groups
        "cognito-idp:AdminListGroupsForUser",      # List user's groups
        "cognito-idp:AdminSetUserPassword",        # Set user passwords
        "cognito-idp:GetUser",                     # Get authenticated user
        "cognito-idp:GlobalSignOut",               # Sign out users
        "cognito-idp:VerifyUserAttribute",         # Verify user attributes
        "cognito-idp:ForgotPassword",              # Initiate forgot password
        "cognito-idp:ConfirmForgotPassword",       # Confirm forgot password
        "cognito-idp:ChangePassword"               # Change password
      ],
      Resource = aws_cognito_user_pool.user_pool.arn, # Specific User Pool
      Effect   = "Allow"
    }]
  })
}

# Output the Cognito User Pool ID
output "cognito_user_pool_id" {
  value       = aws_cognito_user_pool.user_pool.id
  description = "The ID of the Cognito User Pool"
}

# Output the Cognito App Client ID
output "cognito_app_client_id" {
  value       = aws_cognito_user_pool_client.app_client.id
  description = "The ID of the Cognito App Client"
}

# Output the Cognito domain URL
output "cognito_domain" {
  value       = "https://${aws_cognito_user_pool_domain.main.domain}.auth.${var.aws_region}.amazoncognito.com"
  description = "The domain name of the Cognito User Pool"
}

# Output the Identity Pool ID
output "identity_pool_id" {
  value       = aws_cognito_identity_pool.identity_pool.id
  description = "The ID of the Cognito Identity Pool"
}
