interface DisplayProps {
  value: string;
  expression?: string;
}

/** Keeps the display on one line by shrinking the font as digits are added. */
function valueSizeClass(value: string): string {
  const significantDigits = value.replace(/[-.]/g, "").length;
  if (significantDigits > 12) return "display-value--xs";
  if (significantDigits > 9) return "display-value--sm";
  return "";
}

export function Display({ value, expression }: DisplayProps) {
  return (
    <div className="display">
      <div className="display-expression">{expression ?? "\u00A0"}</div>
      <div className={`display-value ${valueSizeClass(value)}`.trim()} data-testid="display-value">
        {value}
      </div>
    </div>
  );
}