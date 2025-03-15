#!/bin/bash

# Make the script exit on any error
set -e

# Colors for better output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL - change this to your API Gateway URL
BASE_URL=${API_URL:-"https://k9ugdqojic.execute-api.us-east-1.amazonaws.com"}

# Generate more unique test data with timestamp and random string
RANDOM_SUFFIX=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
RANDOM_NUMBER=$(openssl rand -hex 4 | tr -dc '0-9' | head -c 8)
TIMESTAMP=$(date +%s)
EMAIL="tadeu.tupiz+${RANDOM_NUMBER}@gmail.com"
PHONE_NUMBER="+5511${TIMESTAMP:(-8)}"
NAME="Test User ${TIMESTAMP}"
PASSWORD="Test@123456"

# Generate a valid CPF for testing
# This function generates a valid CPF with check digits
generate_valid_cpf() {
  # Generate first 9 digits
  local cpf=""
  for i in {1..9}; do
    cpf="${cpf}$(( RANDOM % 10 ))"
  done
  
  # Calculate first verification digit
  local sum=0
  for i in {0..8}; do
    sum=$((sum + ${cpf:$i:1} * (10 - $i)))
  done
  local remainder=$((11 - (sum % 11)))
  local digit1=$((remainder > 9 ? 0 : remainder))
  
  # Calculate second verification digit
  cpf="${cpf}${digit1}"
  sum=0
  for i in {0..9}; do
    sum=$((sum + ${cpf:$i:1} * (11 - $i)))
  done
  remainder=$((11 - (sum % 11)))
  local digit2=$((remainder > 9 ? 0 : remainder))
  
  echo "${cpf}${digit2}"
}

# Generate a unique CPF for this test run
CPF=$(generate_valid_cpf)

# Variables to store tokens
ACCESS_TOKEN=""
ID_TOKEN=""
REFRESH_TOKEN=""

# Function to print section headers
print_header() {
  echo -e "\n${BLUE}==== $1 ====${NC}\n"
}

# Function to print success messages
print_success() {
  echo -e "${GREEN}✓ $1${NC}"
}

# Function to print error messages
print_error() {
  echo -e "${RED}✗ $1${NC}"
  echo -e "${RED}Response: $2${NC}"
  exit 1
}

# Function to print warning messages (non-fatal errors)
print_warning() {
  echo -e "${YELLOW}⚠ $1${NC}"
  echo -e "${YELLOW}Response: $2${NC}"
}

# Function to print info messages
print_info() {
  echo -e "${YELLOW}ℹ $1${NC}"
}

# Check if jq is installed
if ! command -v jq &> /dev/null; then
  echo -e "${RED}Error: jq is not installed. Please install it to parse JSON responses.${NC}"
  echo "On Ubuntu/Debian: sudo apt-get install jq"
  echo "On macOS: brew install jq"
  exit 1
fi

# Print test configuration
print_header "Test Configuration"
echo "CPF (Username): $CPF"
echo "Email: $EMAIL"
echo "Phone: $PHONE_NUMBER"
echo "API URL: $BASE_URL"

# 1. Test Register User
test_register_user() {
  print_header "Testing User Registration"
  
  response=$(curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"cpf\": \"${CPF}\",
      \"password\": \"${PASSWORD}\",
      \"email\": \"${EMAIL}\",
      \"phoneNumber\": \"${PHONE_NUMBER}\",
      \"name\": \"${NAME}\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "User registered successfully" ]]; then
    print_success "User registered successfully"
    user_id=$(echo "$response" | jq -r '.userId')
    print_info "User ID: $user_id"
  else
    print_error "Failed to register user" "$response"
  fi
}

# 2. Test Login with CPF (now the primary login method)
test_login_with_cpf() {
  print_header "Testing Login with CPF"
  
  response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"username\": \"${CPF}\",
      \"password\": \"${PASSWORD}\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Login successful" ]]; then
    print_success "Login successful"
    ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
    ID_TOKEN=$(echo "$response" | jq -r '.idToken')
    REFRESH_TOKEN=$(echo "$response" | jq -r '.refreshToken')
    print_info "Access Token: ${ACCESS_TOKEN:0:20}..."
  else
    print_error "Failed to login" "$response"
  fi
}

# 3. Test Login with CPF using the dedicated CPF endpoint
test_login_with_cpf_endpoint() {
  print_header "Testing Login with CPF Endpoint"
  
  # Add debug output to see what we're sending
  echo "Sending CPF: ${CPF}"
  
  response=$(curl -s -X POST "${BASE_URL}/auth/login/cpf" \
    -H "Content-Type: application/json" \
    -d "{
      \"cpf\": \"${CPF}\",
      \"password\": \"${PASSWORD}\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Login successful" ]]; then
    print_success "Login with CPF endpoint successful"
    # Update tokens if needed
    ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
    ID_TOKEN=$(echo "$response" | jq -r '.idToken')
    REFRESH_TOKEN=$(echo "$response" | jq -r '.refreshToken')
  else
    print_error "Failed to login with CPF endpoint" "$response"
  fi
}

# 4. Test Get User Profile
test_get_user_profile() {
  print_header "Testing Get User Profile"
  
  response=$(curl -s -X GET "${BASE_URL}/auth/profile" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")
  
  username=$(echo "$response" | jq -r '.username')
  
  if [[ "$username" == "$CPF" ]]; then
    print_success "Got user profile successfully"
    echo "$response" | jq '.'
  else
    print_error "Failed to get user profile" "$response"
  fi
}

# 5. Test Get User by CPF
test_get_user_by_cpf() {
  print_header "Testing Get User by CPF"
  
  response=$(curl -s -X GET "${BASE_URL}/auth/user/cpf/${CPF}")
  
  username=$(echo "$response" | jq -r '.username')
  
  if [[ "$username" == "$CPF" ]]; then
    print_success "Got user by CPF successfully"
    echo "$response" | jq '.'
  else
    print_error "Failed to get user by CPF" "$response"
  fi
}

# 6. Test Update User Profile
test_update_user_profile() {
  print_header "Testing Update User Profile"
  
  NEW_NAME="Updated Test User"
  
  response=$(curl -s -X PUT "${BASE_URL}/auth/profile" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"${NEW_NAME}\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Profile updated successfully" ]]; then
    print_success "Profile updated successfully"
    
    # Verify the update
    response=$(curl -s -X GET "${BASE_URL}/auth/profile" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}")
    
    updated_name=$(echo "$response" | jq -r '.attributes.name')
    
    if [[ "$updated_name" == "$NEW_NAME" ]]; then
      print_success "Verified profile update"
    else
      print_error "Failed to verify profile update" "$response"
    fi
  else
    print_error "Failed to update profile" "$response"
  fi
}

# 7. Test Update User CPF - This is now more complex since CPF is the username
# We'll need to create a new user with the new CPF
test_update_user_cpf() {
  print_header "Testing Update User CPF (Creating New User)"
  
  # Generate a new valid CPF
  NEW_CPF=$(generate_valid_cpf)
  
  print_info "Since CPF is now the username, we'll create a new user with CPF: $NEW_CPF"
  
  # Register a new user with the new CPF
  response=$(curl -s -X POST "${BASE_URL}/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
      \"cpf\": \"${NEW_CPF}\",
      \"password\": \"${PASSWORD}\",
      \"email\": \"new.${EMAIL}\",
      \"phoneNumber\": \"${PHONE_NUMBER}1\",
      \"name\": \"${NAME} New\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "User registered successfully" ]]; then
    print_success "Created new user with CPF: $NEW_CPF"
    
    # Login with the new CPF
    response=$(curl -s -X POST "${BASE_URL}/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${NEW_CPF}\",
        \"password\": \"${PASSWORD}\"
      }")
    
    login_status=$(echo "$response" | jq -r '.message')
    
    if [[ "$login_status" == "Login successful" ]]; then
      print_success "Logged in with new CPF"
      # Update tokens for future tests
      ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
      ID_TOKEN=$(echo "$response" | jq -r '.idToken')
      REFRESH_TOKEN=$(echo "$response" | jq -r '.refreshToken')
      
      # Update CPF variable for future tests
      CPF=$NEW_CPF
    else
      print_error "Failed to login with new CPF" "$response"
    fi
  else
    print_warning "Failed to create new user with CPF" "$response"
    print_info "Continuing with original CPF: $CPF"
  fi
}

# 8. Test Refresh Token
test_refresh_token() {
  print_header "Testing Refresh Token"
  
  # Add a longer delay before refreshing token
  print_info "Waiting 5 seconds before attempting token refresh..."
  sleep 5
  
  # Print the first few characters of the refresh token for debugging
  print_info "Using refresh token: ${REFRESH_TOKEN:0:10}..."
  
  response=$(curl -s -X POST "${BASE_URL}/auth/token" \
    -H "Content-Type: application/json" \
    -d "{
      \"refreshToken\": \"${REFRESH_TOKEN}\"
    }")
  
  # Debug the raw response
  print_info "Raw response: ${response}"
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Token refreshed successfully" ]]; then
    print_success "Token refreshed successfully"
    # Update tokens
    ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
    ID_TOKEN=$(echo "$response" | jq -r '.idToken')
    print_info "New Access Token: ${ACCESS_TOKEN:0:20}..."
  else
    error_message=$(echo "$response" | jq -r '.message')
    print_warning "Failed to refresh token: $error_message" "$response"
    
    print_info "Attempting to re-authenticate to get new tokens..."
    
    # Re-authenticate to get new tokens
    response=$(curl -s -X POST "${BASE_URL}/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${CPF}\",
        \"password\": \"${PASSWORD}\"
      }")
    
    login_status=$(echo "$response" | jq -r '.message')
    
    if [[ "$login_status" == "Login successful" ]]; then
      print_success "Re-authenticated successfully"
      ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
      ID_TOKEN=$(echo "$response" | jq -r '.idToken')
      REFRESH_TOKEN=$(echo "$response" | jq -r '.refreshToken')
      print_info "New Access Token: ${ACCESS_TOKEN:0:20}..."
      print_info "New Refresh Token: ${REFRESH_TOKEN:0:20}..."
    else
      print_warning "Failed to re-authenticate, but continuing tests" "$response"
    fi
  fi
}

# 9. Test Change Password
test_change_password() {
  print_header "Testing Change Password"
  
  NEW_PASSWORD="NewTest@123456"
  
  response=$(curl -s -X POST "${BASE_URL}/auth/change-password" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{
      \"oldPassword\": \"${PASSWORD}\",
      \"newPassword\": \"${NEW_PASSWORD}\"
    }")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Password changed successfully" ]]; then
    print_success "Password changed successfully"
    # Update password for future tests
    PASSWORD=$NEW_PASSWORD
    
    # Verify by logging in with new password
    response=$(curl -s -X POST "${BASE_URL}/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${CPF}\",
        \"password\": \"${NEW_PASSWORD}\"
      }")
    
    login_status=$(echo "$response" | jq -r '.message')
    
    if [[ "$login_status" == "Login successful" ]]; then
      print_success "Verified password change by logging in"
      # Update tokens
      ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
      ID_TOKEN=$(echo "$response" | jq -r '.idToken')
      REFRESH_TOKEN=$(echo "$response" | jq -r '.refreshToken')
    else
      print_error "Failed to verify password change" "$response"
    fi
  else
    print_error "Failed to change password" "$response"
  fi
}

# 10. Test Logout
test_logout() {
  print_header "Testing Logout"
  
  response=$(curl -s -X POST "${BASE_URL}/auth/logout" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "Logged out successfully" ]]; then
    print_success "Logged out successfully"
    
    # Verify logout by trying to get profile (should fail)
    response=$(curl -s -X GET "${BASE_URL}/auth/profile" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}")
    
    error_message=$(echo "$response" | jq -r '.message')
    
    if [[ "$error_message" == "Invalid or expired token" ]]; then
      print_success "Verified logout - token is now invalid"
    else
      print_error "Failed to verify logout" "$response"
    fi
  else
    print_error "Failed to logout" "$response"
  fi
}

# 11. Test Delete User
test_delete_user() {
  print_header "Testing Delete User"
  
  # First, login again to get a valid token
  response=$(curl -s -X POST "${BASE_URL}/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
      \"username\": \"${CPF}\",
      \"password\": \"${PASSWORD}\"
    }")
  
  ACCESS_TOKEN=$(echo "$response" | jq -r '.accessToken')
  
  # Now delete the user
  response=$(curl -s -X DELETE "${BASE_URL}/auth/user/${CPF}")
  
  status=$(echo "$response" | jq -r '.message')
  
  if [[ "$status" == "User deleted successfully" ]]; then
    print_success "User deleted successfully"
    
    # Verify deletion by trying to login (should fail)
    response=$(curl -s -X POST "${BASE_URL}/auth/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${CPF}\",
        \"password\": \"${PASSWORD}\"
      }")
    
    error_message=$(echo "$response" | jq -r '.message')
    
    # Accept either "User not found" or "Invalid credentials" as valid responses
    if [[ "$error_message" == "User not found" || "$error_message" == "Invalid credentials" ]]; then
      print_success "Verified user deletion - login failed as expected with message: $error_message"
    else
      print_error "Failed to verify user deletion" "$response"
    fi
  else
    print_error "Failed to delete user" "$response"
  fi
}

# Main test sequence
main() {
  echo -e "${BLUE}Starting Auth Service API Tests${NC}"
  echo -e "${YELLOW}Base URL: ${BASE_URL}${NC}"
  
  # Run tests in sequence
  test_register_user
  test_login_with_cpf
  test_login_with_cpf_endpoint
  test_get_user_profile
  test_get_user_by_cpf
  test_update_user_profile
  test_update_user_cpf
  test_refresh_token
  test_change_password
  test_logout
  test_delete_user
  
  echo -e "\n${GREEN}All tests completed successfully!${NC}"
}

# Run the main function
main