import type { Operation } from "../types";

interface KeypadProps {
  onDigit: (digit: string) => void;
  onDecimal: () => void;
  onOperation: (operation: Operation) => void;
  onEvaluate: () => void;
  onClear: () => void;
  disabled: boolean;
}

/**
 * Keypad is purely presentational: it renders buttons and calls the
 * callbacks it was given. All calculator state and validation lives in
 * useCalculator, which keeps this component trivial to test and reuse.
 */
export function Keypad({ onDigit, onDecimal, onOperation, onEvaluate, onClear, disabled }: KeypadProps) {
  function digitButton(digit: string) {
    return (
      <button key={digit} type="button" className="btn" onClick={() => onDigit(digit)} disabled={disabled}>
        {digit}
      </button>
    );
  }

  return (
    <div className="keypad">
      <button type="button" className="btn btn-clear" onClick={onClear} disabled={disabled}>
        C
      </button>
      <button type="button" className="btn btn-function" onClick={() => onOperation("sqrt")} disabled={disabled}>
        √
      </button>
      <button type="button" className="btn btn-operator" onClick={() => onOperation("power")} disabled={disabled}>
        ^
      </button>
      <button type="button" className="btn btn-operator" onClick={() => onOperation("divide")} disabled={disabled}>
        ÷
      </button>

      {["7", "8", "9"].map(digitButton)}
      <button type="button" className="btn btn-operator" onClick={() => onOperation("multiply")} disabled={disabled}>
        ×
      </button>

      {["4", "5", "6"].map(digitButton)}
      <button type="button" className="btn btn-operator" onClick={() => onOperation("subtract")} disabled={disabled}>
        −
      </button>

      {["1", "2", "3"].map(digitButton)}
      <button type="button" className="btn btn-operator" onClick={() => onOperation("add")} disabled={disabled}>
        +
      </button>

      <button type="button" className="btn btn-zero" onClick={() => onDigit("0")} disabled={disabled}>
        0
      </button>
      <button type="button" className="btn" onClick={onDecimal} disabled={disabled}>
        .
      </button>
      <button
        type="button"
        className="btn btn-equals"
        onClick={onEvaluate}
        disabled={disabled}
        aria-label="equals"
      >
        =
      </button>
    </div>
  );
}