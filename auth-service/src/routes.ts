import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";

// Define the route handler type
export type RouteHandler = (
  event: APIGatewayProxyEvent
) => Promise<APIGatewayProxyResult>;

// Define the route configuration type
export interface RouteConfig {
  method: string;
  path: string;
  handler: RouteHandler;
  description?: string;
}

// Import your handler functions
import { index } from "./handlers";
import {
  changePassword,
  deleteUser,
  forgotPassword,
  getUserByCPF,
  getUserProfile,
  loginUser,
  loginUserByCPF,
  logout,
  refreshToken,
  registerUser,
  resetPassword,
  updateUserCPF,
  updateUserProfile,
  verifyAttribute,
} from "./handlers/auth";

// Define your routes
export const routes: RouteConfig[] = [
  {
    method: "POST",
    path: "/auth/register",
    handler: registerUser,
    description: "Register a new user",
  },
  {
    method: "GET",
    path: "/",
    handler: index,
    description: "Index route",
  },
  // Authentication routes
  {
    method: "POST",
    path: "/auth/login",
    handler: loginUser,
    description: "Login with username/email/phone and password",
  },
  {
    method: "POST",
    path: "/auth/login/cpf",
    handler: loginUserByCPF,
    description: "Login with CPF and password",
  },

  // User management routes
  {
    method: "GET",
    path: "/auth/user/cpf/{cpf}",
    handler: getUserByCPF,
    description: "Get user by CPF",
  },
  {
    method: "GET",
    path: "/auth/profile",
    handler: getUserProfile,
    description: "Get authenticated user's profile",
  },
  {
    method: "PUT",
    path: "/auth/profile",
    handler: updateUserProfile,
    description: "Update user profile",
  },
  {
    method: "PUT",
    path: "/auth/user/cpf",
    handler: updateUserCPF,
    description: "Update user CPF",
  },
  {
    method: "DELETE",
    path: "/auth/user/{username}",
    handler: deleteUser,
    description: "Delete user account",
  },

  // Token management routes
  {
    method: "POST",
    path: "/auth/token",
    handler: refreshToken,
    description: "Refresh authentication token",
  },

  // Password management routes
  {
    method: "POST",
    path: "/auth/forgot-password",
    handler: forgotPassword,
    description: "Initiate forgot password flow",
  },
  {
    method: "POST",
    path: "/auth/reset-password",
    handler: resetPassword,
    description: "Reset password with confirmation code",
  },
  {
    method: "POST",
    path: "/auth/change-password",
    handler: changePassword,
    description: "Change password (authenticated)",
  },

  // Verification routes
  {
    method: "POST",
    path: "/auth/verify",
    handler: verifyAttribute,
    description: "Verify user attribute",
  },

  // Session management
  {
    method: "POST",
    path: "/auth/logout",
    handler: logout,
    description: "Logout user",
  },
];

// Export a function to generate Terraform-compatible route keys
export function generateTerraformRouteKeys(): string[] {
  return routes.map((route) => `${route.method} ${route.path}`);
}
