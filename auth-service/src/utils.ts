import { ApiResponse, ErrorResponse, SuccessResponse } from "./types";

// Helper function to create a response
export function createResponse<T>(
  statusCode: number,
  body: SuccessResponse<T> | ErrorResponse
): ApiResponse {
  return {
    statusCode,
    headers: {
      "Content-Type": "application/json",
      "Access-Control-Allow-Origin": "*", // This should be configured based on your CORS needs
      "Access-Control-Allow-Credentials": "true",
    },
    body: JSON.stringify(body),
    isBase64Encoded: false,
  };
}
