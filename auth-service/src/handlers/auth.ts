import {
  AdminCreateUserCommand,
  AdminDeleteUserCommand,
  AdminGetUserCommand,
  AdminInitiateAuthCommand,
  AdminRespondToAuthChallengeCommand,
  AdminSetUserPasswordCommand,
  AdminUpdateUserAttributesCommand,
  ChangePasswordCommand,
  CognitoIdentityProviderClient,
  ConfirmForgotPasswordCommand,
  ForgotPasswordCommand,
  GetUserCommand,
  GlobalSignOutCommand,
  VerifyUserAttributeCommand,
} from "@aws-sdk/client-cognito-identity-provider";
import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";

// Initialize the Cognito client
const cognitoClient = new CognitoIdentityProviderClient({
  region: process.env.AWS_REGION || "us-east-1",
});

// Constants
const USER_POOL_ID = process.env.COGNITO_USER_POOL_ID || "";
const CLIENT_ID = process.env.COGNITO_CLIENT_ID || "";

// Helper function to parse request body
const parseBody = (event: APIGatewayProxyEvent) => {
  try {
    return event.body ? JSON.parse(event.body) : {};
  } catch (error) {
    throw new Error("Invalid request body");
  }
};

// Helper function for API responses
const response = (statusCode: number, body: any): APIGatewayProxyResult => {
  return {
    statusCode,
    headers: {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Credentials": true,
    },
    body: JSON.stringify(body),
  };
};

// Extract token from Authorization header
const extractToken = (event: APIGatewayProxyEvent): string | null => {
  const authHeader = event.headers.Authorization || event.headers.authorization;
  if (!authHeader) return null;

  const parts = authHeader.split(" ");
  if (parts.length !== 2 || parts[0] !== "Bearer") return null;

  return parts[1];
};

// CPF validation
const isValidCPF = (cpf: string): boolean => {
  // Remove non-numeric characters
  const cleanCPF = cpf.replace(/\D/g, "");

  // Check if it has 11 digits
  if (cleanCPF.length !== 11) return false;

  // Check if all digits are the same (invalid CPF)
  if (/^(\d)\1+$/.test(cleanCPF)) return false;

  // Calculate first verification digit
  let sum = 0;
  for (let i = 0; i < 9; i++) {
    sum += parseInt(cleanCPF.charAt(i)) * (10 - i);
  }
  let remainder = 11 - (sum % 11);
  let digit1 = remainder > 9 ? 0 : remainder;

  // Calculate second verification digit
  sum = 0;
  for (let i = 0; i < 10; i++) {
    sum += parseInt(cleanCPF.charAt(i)) * (11 - i);
  }
  remainder = 11 - (sum % 11);
  let digit2 = remainder > 9 ? 0 : remainder;

  // Check if calculated verification digits match the provided ones
  return (
    digit1 === parseInt(cleanCPF.charAt(9)) &&
    digit2 === parseInt(cleanCPF.charAt(10))
  );
};

// Format phone number to E.164 format
const formatPhoneNumber = (phoneNumber: string): string => {
  if (!phoneNumber) return "";

  // If already in E.164 format, return as is
  if (phoneNumber.startsWith("+")) return phoneNumber;

  // Remove non-numeric characters
  const cleanNumber = phoneNumber.replace(/\D/g, "");

  // Add Brazil country code if not present
  if (cleanNumber.length === 11 || cleanNumber.length === 10) {
    return `+55${cleanNumber}`;
  }

  // If already has country code
  if (
    cleanNumber.startsWith("55") &&
    (cleanNumber.length === 13 || cleanNumber.length === 12)
  ) {
    return `+${cleanNumber}`;
  }

  // Return original if format is unknown
  return `+${cleanNumber}`;
};

// Register a new user
export const registerUser = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { password, email, phoneNumber, name, cpf } = parseBody(event);

    // CPF is now required as it will be the username
    if (!cpf || !password) {
      return response(400, { message: "CPF and password are required" });
    }

    // Validate CPF
    if (!isValidCPF(cpf)) {
      return response(400, { message: "Invalid CPF format" });
    }

    // Format phone number if provided
    const formattedPhone = phoneNumber
      ? formatPhoneNumber(phoneNumber)
      : undefined;

    // Prepare user attributes
    const userAttributes = [];

    if (email) {
      userAttributes.push({ Name: "email", Value: email });
      userAttributes.push({ Name: "email_verified", Value: "true" });
    }

    if (formattedPhone) {
      userAttributes.push({ Name: "phone_number", Value: formattedPhone });
      userAttributes.push({ Name: "phone_number_verified", Value: "true" });
    }

    if (name) {
      userAttributes.push({ Name: "name", Value: name });
    }

    // Still store CPF as custom attribute for validation purposes
    userAttributes.push({ Name: "custom:cpf", Value: cpf });

    // Create user command - use CPF as username
    const command = new AdminCreateUserCommand({
      UserPoolId: USER_POOL_ID,
      Username: cpf, // CPF is now the username
      TemporaryPassword: password,
      MessageAction: "SUPPRESS",
      UserAttributes: userAttributes,
    });

    const result = await cognitoClient.send(command);

    // Set permanent password
    if (result.User) {
      const setPasswordCommand = new AdminSetUserPasswordCommand({
        UserPoolId: USER_POOL_ID,
        Username: cpf, // CPF is the username
        Password: password,
        Permanent: true,
      });

      await cognitoClient.send(setPasswordCommand);
    }

    return response(201, {
      message: "User registered successfully",
      userId: result.User?.Username,
      userSub: result.User?.Attributes?.find((attr) => attr.Name === "sub")
        ?.Value,
    });
  } catch (error: any) {
    console.error("Error in registerUser:", error);

    if (error.name === "UsernameExistsException") {
      return response(409, { message: "User with this CPF already exists" });
    }

    if (error.name === "InvalidParameterException") {
      return response(400, { message: error.message });
    }

    return response(500, {
      message: error.message || "Error registering user",
    });
  }
};

// Login with username/email/phone and password
export const loginUser = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { username, password } = parseBody(event);

    if (!username || !password) {
      return response(400, { message: "Username and password are required" });
    }

    // If username looks like a CPF, validate it
    if (username.length === 11 && /^\d+$/.test(username)) {
      if (!isValidCPF(username)) {
        return response(400, { message: "Invalid CPF format" });
      }
    }

    const command = new AdminInitiateAuthCommand({
      UserPoolId: USER_POOL_ID,
      ClientId: CLIENT_ID,
      AuthFlow: "ADMIN_USER_PASSWORD_AUTH",
      AuthParameters: {
        USERNAME: username,
        PASSWORD: password,
      },
    });

    const result = await cognitoClient.send(command);

    // Handle new password required challenge
    if (result.ChallengeName === "NEW_PASSWORD_REQUIRED") {
      return response(200, {
        message: "New password required",
        session: result.Session,
        challengeName: result.ChallengeName,
      });
    }

    // Return tokens on successful authentication
    return response(200, {
      message: "Login successful",
      accessToken: result.AuthenticationResult?.AccessToken,
      idToken: result.AuthenticationResult?.IdToken,
      refreshToken: result.AuthenticationResult?.RefreshToken,
      expiresIn: result.AuthenticationResult?.ExpiresIn,
    });
  } catch (error: any) {
    console.error("Error in loginUser:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid credentials" });
    }

    if (error.name === "UserNotFoundException") {
      return response(404, { message: "User not found" });
    }

    return response(500, { message: error.message || "Error during login" });
  }
};

// Login with CPF and password
export const loginUserByCPF = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { cpf, password } = parseBody(event);

    console.log(`Attempting to login with CPF: ${cpf}`);

    if (!cpf || !password) {
      return response(400, { message: "CPF and password are required" });
    }

    if (!isValidCPF(cpf)) {
      return response(400, { message: "Invalid CPF format" });
    }

    // Since CPF is the username, we can authenticate directly
    const command = new AdminInitiateAuthCommand({
      UserPoolId: USER_POOL_ID,
      ClientId: CLIENT_ID,
      AuthFlow: "ADMIN_USER_PASSWORD_AUTH",
      AuthParameters: {
        USERNAME: cpf, // CPF is the username
        PASSWORD: password,
      },
    });

    try {
      const result = await cognitoClient.send(command);

      // Handle new password required challenge
      if (result.ChallengeName === "NEW_PASSWORD_REQUIRED") {
        return response(200, {
          message: "New password required",
          session: result.Session,
          challengeName: result.ChallengeName,
        });
      }

      // Return tokens on successful authentication
      return response(200, {
        message: "Login successful",
        accessToken: result.AuthenticationResult?.AccessToken,
        idToken: result.AuthenticationResult?.IdToken,
        refreshToken: result.AuthenticationResult?.RefreshToken,
        expiresIn: result.AuthenticationResult?.ExpiresIn,
      });
    } catch (error: any) {
      console.error("Error in CPF authentication:", error);

      if (error.name === "NotAuthorizedException") {
        return response(401, { message: "Invalid credentials" });
      }

      if (error.name === "UserNotFoundException") {
        return response(404, { message: "User not found with this CPF" });
      }

      throw error; // Re-throw for the outer catch block
    }
  } catch (error: any) {
    console.error("Error in loginUserByCPF:", error);
    return response(500, {
      message: "Error during login with CPF",
      details: error.message,
    });
  }
};

// Get user by CPF
export const getUserByCPF = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const cpf = event.pathParameters?.cpf;

    if (!cpf) {
      return response(400, { message: "CPF is required" });
    }

    if (!isValidCPF(cpf)) {
      return response(400, { message: "Invalid CPF format" });
    }

    const user = await findUserByCPF(cpf);
    return response(200, user);
  } catch (error: any) {
    if (error.name === "UserNotFoundException") {
      return response(404, { message: "User not found with this CPF" });
    }
    console.error("Error in getUserByCPF:", error);
    return response(500, { message: error.message || "Error retrieving user" });
  }
};

// Helper function to find user by CPF
const findUserByCPF = async (cpf: string) => {
  try {
    const command = new AdminGetUserCommand({
      UserPoolId: USER_POOL_ID,
      Username: cpf,
    });

    const user = await cognitoClient.send(command);

    // Extract relevant user information
    const userAttributes = user.UserAttributes?.reduce((acc: any, attr) => {
      // Don't expose sensitive information
      if (!["sub", "password"].includes(attr.Name ?? "")) {
        acc[attr.Name?.replace("custom:", "") ?? ""] = attr.Value;
      }
      return acc;
    }, {});

    return {
      username: user.Username,
      attributes: userAttributes,
      userStatus: user.UserStatus,
      enabled: user.Enabled,
      created: user.UserCreateDate,
      modified: user.UserLastModifiedDate,
    };
  } catch (error: any) {
    throw error;
  }
};

// Get authenticated user's profile
export const getUserProfile = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    console.log("token", token);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const command = new GetUserCommand({
      AccessToken: token,
    });

    const result = await cognitoClient.send(command);

    // Format user attributes
    const userAttributes = result.UserAttributes?.reduce((acc: any, attr) => {
      // Don't expose sensitive information
      if (!["sub", "password"].includes(attr.Name ?? "")) {
        acc[attr.Name?.replace("custom:", "") ?? ""] = attr.Value;
      }
      return acc;
    }, {});

    return response(200, {
      username: result.Username,
      attributes: userAttributes,
    });
  } catch (error: any) {
    console.error("Error in getUserProfile:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid or expired token" });
    }

    return response(500, {
      message: error.message || "Error retrieving user profile",
    });
  }
};

// Update user profile
export const updateUserProfile = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const { name, email, phoneNumber } = parseBody(event);

    // Get current user
    const getUserCommand = new GetUserCommand({
      AccessToken: token,
    });

    const user = await cognitoClient.send(getUserCommand);
    const username = user.Username;

    // Prepare attributes to update
    const userAttributes = [];

    if (name) {
      userAttributes.push({ Name: "name", Value: name });
    }

    if (email) {
      userAttributes.push({ Name: "email", Value: email });
    }

    if (phoneNumber) {
      const formattedPhone = formatPhoneNumber(phoneNumber);
      userAttributes.push({ Name: "phone_number", Value: formattedPhone });
    }

    if (userAttributes.length === 0) {
      return response(400, { message: "No attributes to update" });
    }

    // Update user attributes
    const updateCommand = new AdminUpdateUserAttributesCommand({
      UserPoolId: USER_POOL_ID,
      Username: username,
      UserAttributes: userAttributes,
    });

    await cognitoClient.send(updateCommand);

    return response(200, { message: "Profile updated successfully" });
  } catch (error: any) {
    console.error("Error in updateUserProfile:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid or expired token" });
    }

    if (error.name === "InvalidParameterException") {
      return response(400, { message: error.message });
    }

    return response(500, {
      message: error.message || "Error updating profile",
    });
  }
};

// Update user CPF
export const updateUserCPF = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const { cpf } = parseBody(event);

    if (!cpf) {
      return response(400, { message: "CPF is required" });
    }

    if (!isValidCPF(cpf)) {
      return response(400, { message: "Invalid CPF format" });
    }

    // Check if CPF is already in use
    const existingUser = await findUserByCPF(cpf);

    if (existingUser) {
      // Get current user
      const getUserCommand = new GetUserCommand({
        AccessToken: token,
      });

      const currentUser = await cognitoClient.send(getUserCommand);

      // If CPF is already used by another user
      if (existingUser.username !== currentUser.Username) {
        return response(409, {
          message: "CPF is already in use by another user",
        });
      }
    }

    // Get current user
    const getUserCommand = new GetUserCommand({
      AccessToken: token,
    });

    const user = await cognitoClient.send(getUserCommand);
    const username = user.Username;

    // Update CPF
    const updateCommand = new AdminUpdateUserAttributesCommand({
      UserPoolId: USER_POOL_ID,
      Username: username,
      UserAttributes: [
        {
          Name: "custom:cpf",
          Value: cpf,
        },
      ],
    });

    await cognitoClient.send(updateCommand);

    return response(200, { message: "CPF updated successfully" });
  } catch (error: any) {
    console.error("Error in updateUserCPF:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid or expired token" });
    }

    return response(500, { message: error.message || "Error updating CPF" });
  }
};

// Delete user
export const deleteUser = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const username = event.pathParameters?.username;

    if (!username) {
      return response(400, { message: "Username is required" });
    }

    // Check if user exists
    try {
      const getUserCommand = new AdminGetUserCommand({
        UserPoolId: USER_POOL_ID,
        Username: username,
      });

      await cognitoClient.send(getUserCommand);
    } catch (error: any) {
      if (error.name === "UserNotFoundException") {
        return response(404, { message: "User not found" });
      }
      throw error;
    }

    // Delete user
    const deleteCommand = new AdminDeleteUserCommand({
      UserPoolId: USER_POOL_ID,
      Username: username,
    });

    await cognitoClient.send(deleteCommand);

    return response(200, { message: "User deleted successfully" });
  } catch (error: any) {
    console.error("Error in deleteUser:", error);
    return response(500, { message: error.message || "Error deleting user" });
  }
};

// Refresh token
export const refreshToken = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { refreshToken } = parseBody(event);

    if (!refreshToken) {
      return response(400, { message: "Refresh token is required" });
    }

    const command = new AdminInitiateAuthCommand({
      UserPoolId: USER_POOL_ID,
      ClientId: CLIENT_ID,
      AuthFlow: "REFRESH_TOKEN_AUTH",
      AuthParameters: {
        REFRESH_TOKEN: refreshToken,
      },
    });

    const result = await cognitoClient.send(command);

    return response(200, {
      message: "Token refreshed successfully",
      accessToken: result.AuthenticationResult?.AccessToken,
      idToken: result.AuthenticationResult?.IdToken,
      expiresIn: result.AuthenticationResult?.ExpiresIn,
    });
  } catch (error: any) {
    console.error("Error in refreshToken:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid refresh token" });
    }

    return response(500, {
      message: error.message || "Error refreshing token",
    });
  }
};

// Forgot password
export const forgotPassword = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { username } = parseBody(event);

    if (!username) {
      return response(400, { message: "Username is required" });
    }

    const command = new ForgotPasswordCommand({
      ClientId: CLIENT_ID,
      Username: username,
    });

    await cognitoClient.send(command);

    return response(200, { message: "Password reset code sent successfully" });
  } catch (error: any) {
    console.error("Error in forgotPassword:", error);

    if (error.name === "UserNotFoundException") {
      // For security reasons, don't reveal that the user doesn't exist
      return response(200, {
        message: "If the account exists, a password reset code has been sent",
      });
    }

    if (error.name === "InvalidParameterException") {
      return response(400, { message: error.message });
    }

    return response(500, { message: "Error initiating password reset" });
  }
};

// Reset password with confirmation code
export const resetPassword = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { username, password, confirmationCode } = parseBody(event);

    if (!username || !password || !confirmationCode) {
      return response(400, {
        message: "Username, password, and confirmation code are required",
      });
    }

    const command = new ConfirmForgotPasswordCommand({
      ClientId: CLIENT_ID,
      Username: username,
      Password: password,
      ConfirmationCode: confirmationCode,
    });

    await cognitoClient.send(command);

    return response(200, { message: "Password reset successfully" });
  } catch (error: any) {
    console.error("Error in resetPassword:", error);

    if (error.name === "CodeMismatchException") {
      return response(400, { message: "Invalid confirmation code" });
    }

    if (error.name === "ExpiredCodeException") {
      return response(400, { message: "Confirmation code has expired" });
    }

    if (error.name === "InvalidPasswordException") {
      return response(400, { message: error.message });
    }

    return response(500, {
      message: error.message || "Error resetting password",
    });
  }
};

// Change password (authenticated)
export const changePassword = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const { oldPassword, newPassword } = parseBody(event);

    if (!oldPassword || !newPassword) {
      return response(400, {
        message: "Old password and new password are required",
      });
    }

    const command = new ChangePasswordCommand({
      AccessToken: token,
      PreviousPassword: oldPassword,
      ProposedPassword: newPassword,
    });

    await cognitoClient.send(command);

    return response(200, { message: "Password changed successfully" });
  } catch (error: any) {
    console.error("Error in changePassword:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid credentials or token" });
    }

    if (error.name === "InvalidPasswordException") {
      return response(400, { message: error.message });
    }

    return response(500, {
      message: error.message || "Error changing password",
    });
  }
};

// Verify user attribute
export const verifyAttribute = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const { attributeName, code } = parseBody(event);

    if (!attributeName || !code) {
      return response(400, {
        message: "Attribute name and verification code are required",
      });
    }

    const command = new VerifyUserAttributeCommand({
      AccessToken: token,
      AttributeName: attributeName,
      Code: code,
    });

    await cognitoClient.send(command);

    return response(200, { message: `${attributeName} verified successfully` });
  } catch (error: any) {
    console.error("Error in verifyAttribute:", error);

    if (error.name === "CodeMismatchException") {
      return response(400, { message: "Invalid verification code" });
    }

    if (error.name === "ExpiredCodeException") {
      return response(400, { message: "Verification code has expired" });
    }

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid or expired token" });
    }

    return response(500, {
      message: error.message || "Error verifying attribute",
    });
  }
};

// Logout user
export const logout = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const token = extractToken(event);

    if (!token) {
      return response(401, { message: "Authentication required" });
    }

    const command = new GlobalSignOutCommand({
      AccessToken: token,
    });

    await cognitoClient.send(command);

    return response(200, { message: "Logged out successfully" });
  } catch (error: any) {
    console.error("Error in logout:", error);

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid or expired token" });
    }

    return response(500, { message: error.message || "Error during logout" });
  }
};

// Respond to auth challenge (for NEW_PASSWORD_REQUIRED)
export const respondToAuthChallenge = async (
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> => {
  try {
    const { username, session, newPassword } = parseBody(event);

    if (!username || !session || !newPassword) {
      return response(400, {
        message: "Username, session, and new password are required",
      });
    }

    const command = new AdminRespondToAuthChallengeCommand({
      UserPoolId: USER_POOL_ID,
      ClientId: CLIENT_ID,
      ChallengeName: "NEW_PASSWORD_REQUIRED",
      ChallengeResponses: {
        USERNAME: username,
        NEW_PASSWORD: newPassword,
      },
      Session: session,
    });

    const result = await cognitoClient.send(command);

    return response(200, {
      message: "Password set successfully",
      accessToken: result.AuthenticationResult?.AccessToken,
      idToken: result.AuthenticationResult?.IdToken,
      refreshToken: result.AuthenticationResult?.RefreshToken,
      expiresIn: result.AuthenticationResult?.ExpiresIn,
    });
  } catch (error: any) {
    console.error("Error in respondToAuthChallenge:", error);

    if (error.name === "InvalidPasswordException") {
      return response(400, { message: error.message });
    }

    if (error.name === "NotAuthorizedException") {
      return response(401, { message: "Invalid session" });
    }

    return response(500, {
      message: error.message || "Error responding to challenge",
    });
  }
};
