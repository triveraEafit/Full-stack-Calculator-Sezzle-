// Mirrors the Go backend's request/response shapes exactly, so the frontend
// and backend contract stays in one obvious place.

export type Operation = "add" | "subtract" | "multiply" | "divide" | "power" | "sqrt";

export interface CalculateResponse {
  result: number;
  operation: Operation;
  operands: number[];
}

export interface ApiErrorBody {
  code: string;
  message: string;
}