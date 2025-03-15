import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";
import { createResponse } from "../utils";

export async function registerUser(
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> {
  try {
    // Parse the request body
    const body = event.body ? JSON.parse(event.body) : {};

    return createResponse(201, {
      success: true,
      data: body,
    });
  } catch (error: any) {
    console.error("Registration error:", error);
    return createResponse(error.statusCode || 500, {
      error: error.code || "RegistrationError",
      message: error.message || "An error occurred during registration",
    });
  }
}
