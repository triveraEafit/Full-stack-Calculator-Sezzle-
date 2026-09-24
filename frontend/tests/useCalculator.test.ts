import { act, renderHook, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { useCalculator } from "../src/hooks/useCalculator";

vi.mock("../src/api/calculatorClient", async () => {
  const actual = await vi.importActual<typeof import("../src/api/calculatorClient")>(
    "../src/api/calculatorClient"
  );
  return { ...actual, calculate: vi.fn() };
});

import { ApiError, calculate } from "../src/api/calculatorClient";

const mockedCalculate = vi.mocked(calculate);

beforeEach(() => {
  mockedCalculate.mockReset();
});

describe("useCalculator", () => {
  it("builds up a number by pressing digits", () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("4");
      result.current.inputDigit("2");
    });

    expect(result.current.display).toBe("42");
  });

  it("supports a decimal point and ignores a second one", () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("3");
      result.current.inputDecimal();
      result.current.inputDigit("1");
      result.current.inputDecimal(); // ignored: already has a decimal point
      result.current.inputDigit("4");
    });

    expect(result.current.display).toBe("3.14");
  });

  it("performs addition via the API when evaluate is called", async () => {
    mockedCalculate.mockResolvedValueOnce(6.5);
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("4");
      result.current.chooseOperation("add");
      result.current.inputDigit("2");
    });
    act(() => {
      result.current.evaluate();
    });

    await waitFor(() => expect(result.current.display).toBe("6.5"));
    expect(mockedCalculate).toHaveBeenCalledWith("add", [4, 2]);
  });

  it("chains operations by evaluating the pending one first", async () => {
    mockedCalculate.mockResolvedValueOnce(5); // 2 + 3
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("2");
      result.current.chooseOperation("add");
      result.current.inputDigit("3");
      result.current.chooseOperation("multiply");
    });

    await waitFor(() => expect(mockedCalculate).toHaveBeenCalledWith("add", [2, 3]));
    await waitFor(() => expect(result.current.operation).toBe("multiply"));
    expect(result.current.previousOperand).toBe(5);
  });

  it("treats square root as an immediate unary operation", async () => {
    mockedCalculate.mockResolvedValueOnce(4);
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("1");
      result.current.inputDigit("6");
      result.current.chooseOperation("sqrt");
    });

    await waitFor(() => expect(result.current.display).toBe("4"));
    expect(mockedCalculate).toHaveBeenCalledWith("sqrt", [16]);
  });

  it("surfaces the API error message and resets on failure", async () => {
    mockedCalculate.mockRejectedValueOnce(new ApiError("DIVISION_BY_ZERO", "Cannot divide by zero"));
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("5");
      result.current.chooseOperation("divide");
      result.current.inputDigit("0");
    });
    act(() => {
      result.current.evaluate();
    });

    await waitFor(() => expect(result.current.error).toBe("Cannot divide by zero"));
    expect(result.current.display).toBe("0");
  });

  it("clears back to the initial state", () => {
    const { result } = renderHook(() => useCalculator());

    act(() => {
      result.current.inputDigit("9");
      result.current.clear();
    });

    expect(result.current.display).toBe("0");
    expect(result.current.operation).toBeNull();
  });
});