import { APIGatewayProxyEvent, APIGatewayProxyResult } from "aws-lambda";
import { createResponse } from "../utils";

export async function index(
  event: APIGatewayProxyEvent
): Promise<APIGatewayProxyResult> {
  try {
    return createResponse(200, {
      success: true,
      data: "Hello World",
    });
  } catch (error: any) {
    console.error("Index error:", error);
    return createResponse(error.statusCode || 500, {
      error: error.code || "IndexError",
      message: error.message || "An error occurred during index",
    });
  }
}
