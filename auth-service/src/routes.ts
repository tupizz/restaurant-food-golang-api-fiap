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
import { registerUser } from "./handlers/auth";

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
];

// Export a function to generate Terraform-compatible route keys
export function generateTerraformRouteKeys(): string[] {
  return routes.map((route) => `${route.method} ${route.path}`);
}
