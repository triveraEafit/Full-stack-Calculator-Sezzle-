import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "../src/components/Calculator";

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

describe("Calculator", () => {
  it("performs an addition end-to-end through the UI", async () => {
    mockedCalculate.mockResolvedValueOnce(7);
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "4" }));
    await user.click(screen.getByRole("button", { name: "+" }));
    await user.click(screen.getByRole("button", { name: "3" }));
    await user.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() => expect(screen.getByTestId("display-value")).toHaveTextContent("7"));
    expect(mockedCalculate).toHaveBeenCalledWith("add", [4, 3]);
  });

  it("shows an error banner when the API rejects the request", async () => {
    mockedCalculate.mockRejectedValueOnce(new ApiError("DIVISION_BY_ZERO", "Cannot divide by zero"));
    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "5" }));
    await user.click(screen.getByRole("button", { name: "÷" }));
    await user.click(screen.getByRole("button", { name: "0" }));
    await user.click(screen.getByRole("button", { name: "equals" }));

    await waitFor(() => expect(screen.getByRole("alert")).toHaveTextContent("Cannot divide by zero"));
  });

  it("disables the keypad while a calculation is in flight", async () => {
    let resolveCalculation: (value: number) => void = () => {};
    mockedCalculate.mockImplementationOnce(
      () =>
        new Promise((resolve) => {
          resolveCalculation = resolve;
        })
    );

    const user = userEvent.setup();
    render(<Calculator />);

    await user.click(screen.getByRole("button", { name: "9" }));
    await user.click(screen.getByRole("button", { name: "+" }));
    await user.click(screen.getByRole("button", { name: "1" }));
    await user.click(screen.getByRole("button", { name: "equals" }));

    expect(screen.getByRole("button", { name: "9" })).toBeDisabled();

    resolveCalculation(10);
    await waitFor(() => expect(screen.getByRole("button", { name: "9" })).not.toBeDisabled());
  });
});