import { useEffect, useReducer } from "react";
import { ApiError, calculate } from "../api/calculatorClient";
import type { Operation } from "../types";

const INITIAL_DISPLAY = "0";
const MAX_DIGITS = 15;

interface PendingCalculation {
  operation: Operation;
  operands: number[];
  // The operation to queue once this one resolves, if any (used for
  // chaining, e.g. "2 + 3 ×" evaluates the "+" and then waits on "×").
  nextOperation: Operation | null;
}

interface CalculatorState {
  display: string;
  previousOperand: number | null;
  operation: Operation | null;
  // True when the next digit pressed should start a fresh number rather
  // than append to what's on screen (e.g. right after an operator or a
  // result).
  overwrite: boolean;
  error: string | null;
  isLoading: boolean;
  // Non-null while an effect is meant to call the API. The reducer only
  // ever *records* that a calculation should happen; it never performs it,
  // since reducers must stay pure. A useEffect below does the actual call.
  pendingCalculation: PendingCalculation | null;
}

const initialState: CalculatorState = {
  display: INITIAL_DISPLAY,
  previousOperand: null,
  operation: null,
  overwrite: true,
  error: null,
  isLoading: false,
  pendingCalculation: null,
};

type Action =
  | { type: "INPUT_DIGIT"; digit: string }
  | { type: "INPUT_DECIMAL" }
  | { type: "CLEAR" }
  | { type: "CHOOSE_OPERATION"; operation: Operation }
  | { type: "EVALUATE" }
  | { type: "CALCULATION_SUCCEEDED"; result: number; nextOperation: Operation | null }
  | { type: "CALCULATION_FAILED"; message: string };

/** Formats an API result for display, trimming floating-point noise. */
function formatResult(value: number): string {
  const rounded = Number(value.toPrecision(10));
  return rounded.toString();
}

/**
 * reducer is the single source of truth for state transitions. Using a
 * reducer (rather than reading a closed-over state variable and calling
 * setState) guarantees that when several actions are dispatched before a
 * re-render happens, each one is applied on top of the *result* of the
 * previous one, not a shared stale snapshot — which is what this hook
 * actually needs when e.g. two digits are pressed in quick succession.
 */
function reducer(state: CalculatorState, action: Action): CalculatorState {
  switch (action.type) {
    case "INPUT_DIGIT": {
      if (state.isLoading) return state;

      if (state.overwrite) {
        return { ...state, display: action.digit, overwrite: false, error: null };
      }

      const digitCount = state.display.replace(/[-.]/g, "").length;
      if (digitCount >= MAX_DIGITS) return state;

      const display = state.display === "0" ? action.digit : state.display + action.digit;
      return { ...state, display, error: null };
    }

    case "INPUT_DECIMAL": {
      if (state.isLoading) return state;

      if (state.overwrite) {
        return { ...state, display: "0.", overwrite: false, error: null };
      }

      if (state.display.includes(".")) return state;
      return { ...state, display: state.display + ".", error: null };
    }

    case "CLEAR":
      return initialState;

    case "CHOOSE_OPERATION": {
      if (state.isLoading) return state;

      const currentValue = Number(state.display);

      // Square root is unary: apply it immediately to whatever's on
      // screen, it doesn't wait for a second operand.
      if (action.operation === "sqrt") {
        return {
          ...state,
          isLoading: true,
          error: null,
          pendingCalculation: { operation: "sqrt", operands: [currentValue], nextOperation: null },
        };
      }

      // No pending calculation yet: remember the current value and wait
      // for the second operand.
      if (state.previousOperand === null) {
        return {
          ...state,
          previousOperand: currentValue,
          operation: action.operation,
          overwrite: true,
          error: null,
        };
      }

      // An operator was already chosen but no new number typed yet
      // (e.g. the user changes their mind from "+" to "-"): just swap it.
      if (state.overwrite) {
        return { ...state, operation: action.operation };
      }

      // A full expression is pending (e.g. "5 +" and "3" was just typed):
      // evaluate it now, then queue the next operation on the result.
      return {
        ...state,
        isLoading: true,
        error: null,
        pendingCalculation: {
          operation: state.operation as Operation,
          operands: [state.previousOperand, currentValue],
          nextOperation: action.operation,
        },
      };
    }

    case "EVALUATE": {
      if (state.isLoading) return state;
      if (state.operation === null || state.previousOperand === null || state.overwrite) return state;

      const currentValue = Number(state.display);
      return {
        ...state,
        isLoading: true,
        error: null,
        pendingCalculation: {
          operation: state.operation,
          operands: [state.previousOperand, currentValue],
          nextOperation: null,
        },
      };
    }

    case "CALCULATION_SUCCEEDED":
      return {
        ...state,
        display: formatResult(action.result),
        previousOperand: action.nextOperation ? action.result : null,
        operation: action.nextOperation,
        overwrite: true,
        error: null,
        isLoading: false,
        pendingCalculation: null,
      };

    case "CALCULATION_FAILED":
      return { ...initialState, error: action.message };

    default:
      return state;
  }
}

/**
 * useCalculator owns all calculator state and talks to the backend for
 * every arithmetic operation — this hook does no math itself, it only
 * manages what the user has typed and what to send/show next.
 */
export function useCalculator() {
  const [state, dispatch] = useReducer(reducer, initialState);

  // Runs the actual API call whenever the reducer records a pending
  // calculation. Kept separate from the reducer because reducers must
  // stay pure; this is the one place a real side effect happens.
  useEffect(() => {
    if (!state.pendingCalculation) return;

    const { operation, operands, nextOperation } = state.pendingCalculation;
    let cancelled = false;

    calculate(operation, operands)
      .then((result) => {
        if (!cancelled) dispatch({ type: "CALCULATION_SUCCEEDED", result, nextOperation });
      })
      .catch((err: unknown) => {
        if (cancelled) return;
        const message = err instanceof ApiError ? err.message : "Something went wrong.";
        dispatch({ type: "CALCULATION_FAILED", message });
      });

    return () => {
      cancelled = true;
    };
  }, [state.pendingCalculation]);

  return {
    display: state.display,
    previousOperand: state.previousOperand,
    operation: state.operation,
    error: state.error,
    isLoading: state.isLoading,
    inputDigit: (digit: string) => dispatch({ type: "INPUT_DIGIT", digit }),
    inputDecimal: () => dispatch({ type: "INPUT_DECIMAL" }),
    chooseOperation: (operation: Operation) => dispatch({ type: "CHOOSE_OPERATION", operation }),
    evaluate: () => dispatch({ type: "EVALUATE" }),
    clear: () => dispatch({ type: "CLEAR" }),
  };
}