import { useCalculator } from "../hooks/useCalculator";
import type { Operation } from "../types";
import { Display } from "./Display";
import { ErrorBanner } from "./ErrorBanner";
import { Keypad } from "./Keypad";

const OPERATION_SYMBOLS: Record<Operation, string> = {
  add: "+",
  subtract: "−",
  multiply: "×",
  divide: "÷",
  power: "^",
  sqrt: "√",
};

export function Calculator() {
  const {
    display,
    previousOperand,
    operation,
    error,
    isLoading,
    inputDigit,
    inputDecimal,
    chooseOperation,
    evaluate,
    clear,
  } = useCalculator();

  const expression =
    previousOperand !== null && operation
      ? `${previousOperand} ${OPERATION_SYMBOLS[operation]}`
      : undefined;

  return (
    <div className="calculator">
      <Display value={isLoading ? "…" : display} expression={expression} />
      <ErrorBanner message={error} />
      <Keypad
        onDigit={inputDigit}
        onDecimal={inputDecimal}
        onOperation={chooseOperation}
        onEvaluate={evaluate}
        onClear={clear}
        disabled={isLoading}
      />
    </div>
  );
}