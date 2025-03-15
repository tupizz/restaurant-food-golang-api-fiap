export interface SuccessResponse<T> {
  success: boolean;
  data: T;
}
export interface ErrorResponse {
  error: string;
  message: string;
}

export interface ApiResponse {
  statusCode: number;
  headers: {
    [header: string]: string;
  };
  body: string;
  isBase64Encoded: boolean;
}
