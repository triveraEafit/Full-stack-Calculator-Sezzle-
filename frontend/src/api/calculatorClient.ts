import type { ApiErrorBody, CalculateResponse, Operation } from "../types";

// Falls back to the local Go server's default port. Override with a
// VITE_API_BASE_URL env var when the backend runs somewhere else.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

/**
 * ApiError represents either a rejection from the backend (its `code`
 * matches the backend's error codes, e.g. "DIVISION_BY_ZERO") or a failure
 * to reach/parse a response at all. `message` is always safe to show
 * directly to the user.
 */
export class ApiError extends Error {
  readonly code: string;

  constructor(code: string, message: string) {
    super(message);
    this.name = "ApiError";
    this.code = code;
  }
}

/**
 * calculate sends one operation to the backend and returns the numeric
 * result, or throws an ApiError describing what went wrong.
 */
export async function calculate(operation: Operation, operands: number[]): Promise<number> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE_URL}/api/v1/calculate`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ operation, operands }),
    });
  } catch {
    throw new ApiError("NETWORK_ERROR", "Could not reach the server. Is it running?");
  }

  let body: unknown;
  try {
    body = await response.json();
  } catch {
    throw new ApiError("INVALID_RESPONSE", "The server returned an unreadable response.");
  }

  if (!response.ok) {
    const errorBody = (body as { error?: ApiErrorBody }).error;
    throw new ApiError(
      errorBody?.code ?? "UNKNOWN_ERROR",
      errorBody?.message ?? "Something went wrong."
    );
  }

  return (body as CalculateResponse).result;
}