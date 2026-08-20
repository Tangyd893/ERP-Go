import {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
  useState,
  type KeyboardEvent,
} from "react";
import { Button, Input } from "@erp/shared-react";
import { cn } from "@erp/shared-react";

export interface ScanInputHandle {
  focus: () => void;
}

interface ScanInputProps {
  placeholder?: string;
  hint?: string;
  disabled?: boolean;
  value?: string;
  clearAfterScan?: boolean;
  onChange?: (value: string) => void;
  onScan: (value: string) => void;
  className?: string;
}

const ScanInput = forwardRef<ScanInputHandle, ScanInputProps>(function ScanInput(
  {
    placeholder = "扫描或输入…",
    hint = "",
    disabled = false,
    value,
    clearAfterScan = true,
    onChange,
    onScan,
    className,
  },
  ref
) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [localValue, setLocalValue] = useState(value ?? "");

  useEffect(() => {
    if (value !== undefined) {
      setLocalValue(value);
    }
  }, [value]);

  useImperativeHandle(ref, () => ({
    focus() {
      requestAnimationFrame(() => inputRef.current?.focus());
    },
  }));

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  function updateValue(next: string) {
    setLocalValue(next);
    onChange?.(next);
  }

  function triggerScan() {
    const trimmed = localValue.trim();
    if (!trimmed || disabled) return;
    onScan(trimmed);
    if (clearAfterScan) {
      updateValue("");
    }
  }

  function handleKeyDown(e: KeyboardEvent<HTMLInputElement>) {
    if (e.key === "Enter") {
      e.preventDefault();
      triggerScan();
    }
  }

  return (
    <div className={cn("mb-3 rounded-lg bg-white p-4 shadow-sm", className)}>
      {hint ? (
        <p className="mb-2.5 text-[13px] text-pda-text-secondary">{hint}</p>
      ) : null}
      <div className="flex gap-2">
        <Input
          ref={inputRef}
          value={localValue}
          placeholder={placeholder}
          disabled={disabled}
          autoComplete="off"
          autoCorrect="off"
          spellCheck={false}
          inputMode="text"
          className="min-h-touch flex-1 text-base tracking-wide"
          onChange={(e) => updateValue(e.target.value)}
          onKeyDown={handleKeyDown}
        />
        <Button
          type="button"
          disabled={!localValue.trim() || disabled}
          className="min-h-touch shrink-0 px-4"
          onClick={triggerScan}
        >
          确认
        </Button>
      </div>
    </div>
  );
});

export default ScanInput;
