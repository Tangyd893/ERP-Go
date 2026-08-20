import { useRef, useState } from "react";
import { Button } from "./ui/button";
import { Upload } from "lucide-react";

export function FileUpload({
  accept,
  onFile,
  label = "选择文件",
  disabled,
}: {
  accept?: string;
  onFile: (file: File) => void | Promise<void>;
  label?: string;
  disabled?: boolean;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [name, setName] = useState("");

  return (
    <div className="flex items-center gap-3">
      <input
        ref={inputRef}
        type="file"
        className="hidden"
        accept={accept}
        onChange={async (e) => {
          const file = e.target.files?.[0];
          if (!file) return;
          setName(file.name);
          await onFile(file);
        }}
      />
      <Button
        type="button"
        variant="outline"
        disabled={disabled}
        onClick={() => inputRef.current?.click()}
      >
        <Upload className="size-4" />
        {label}
      </Button>
      {name ? <span className="truncate text-sm text-muted-foreground">{name}</span> : null}
    </div>
  );
}
