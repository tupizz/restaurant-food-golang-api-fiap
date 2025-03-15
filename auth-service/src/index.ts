import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";
import { routes } from "./routes";
import { createResponse } from "./utils";

export async function handler(
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> {
  console.log("Event received:", JSON.stringify(event));

  try {
    // Extract the HTTP method and path from the event
    const httpMethod = event.httpMethod;
    const path = event.path;

    console.log("--------------------------------");
    console.log("Event received:", JSON.stringify(event));
    console.log("httpMethod:", httpMethod);
    console.log("path:", path);
    console.log("--------------------------------");

    // Find the matching route
    const route = routes.find((r) => {
      // Convert API Gateway path parameters to route pattern
      // e.g., /users/123 should match /users/{userId}
      const routePathParts = r.path.split("/");
      const requestPathParts = path.split("/");

      if (routePathParts.length !== requestPathParts.length) return false;
      if (r.method !== httpMethod) return false;

      for (let i = 0; i < routePathParts.length; i++) {
        const routePart = routePathParts[i];
        if (routePart.startsWith("{") && routePart.endsWith("}")) {
          // This is a path parameter, so it matches any value
          continue;
        }
        if (routePart !== requestPathParts[i]) {
          return false;
        }
      }

      return true;
    });

    if (route) {
      // Call the handler for the matched route
      return await route.handler(event);
    }

    // No matching route found
    return createResponse(404, {
      error: "NotFound",
      message: `No handler found for ${httpMethod} ${path}`,
    });
  } catch (error: any) {
    console.error("Error processing request:", error);
    return createResponse(500, {
      error: "InternalServerError",
      message: error.message || "An unexpected error occurred",
    });
  }
}
