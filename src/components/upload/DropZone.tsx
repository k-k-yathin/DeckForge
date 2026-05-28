import { useCallback, useState, type DragEvent } from 'react';

interface DropZoneProps {
  onFileSelect: (file: File) => void;
  accept?: string;
  disabled?: boolean;
}

export function DropZone({
  onFileSelect,
  accept = '.pdf,.docx,.txt',
  disabled,
}: DropZoneProps) {
  const [dragOver, setDragOver] = useState(false);

  const handleDrop = useCallback(
    (e: DragEvent) => {
      e.preventDefault();
      setDragOver(false);
      if (disabled) return;
      const file = e.dataTransfer.files[0];
      if (file) onFileSelect(file);
    },
    [disabled, onFileSelect]
  );

  return (
    <div
      onDragOver={(e) => {
        e.preventDefault();
        setDragOver(true);
      }}
      onDragLeave={() => setDragOver(false)}
      onDrop={handleDrop}
      className={`relative rounded-xl border-2 border-dashed p-12 text-center transition ${
        dragOver
          ? 'border-brand-500 bg-brand-50'
          : 'border-slate-200 bg-white hover:border-brand-300'
      } ${disabled ? 'opacity-50 pointer-events-none' : ''}`}
    >
      <div className="text-4xl">📄</div>
      <p className="mt-4 text-sm font-medium text-slate-700">
        Drag & drop your file here
      </p>
      <p className="mt-1 text-xs text-slate-500">PDF, DOCX, or TXT — max 10MB</p>
      <label className="mt-6 inline-block cursor-pointer">
        <span className="rounded-lg bg-brand-600 px-4 py-2 text-sm font-medium text-white hover:bg-brand-700">
          Browse files
        </span>
        <input
          type="file"
          className="hidden"
          accept={accept}
          disabled={disabled}
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) onFileSelect(file);
          }}
        />
      </label>
    </div>
  );
}
