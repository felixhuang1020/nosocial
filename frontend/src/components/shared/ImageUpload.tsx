import { useState, useRef, useCallback } from 'react';
import { Upload, X, ImagePlus, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import { uploadToOSS, validateImageFile } from '@/lib/upload';

interface ImageUploadProps {
  value?: string;
  onChange: (url: string) => void;
  className?: string;
  /** 指定 OSS 目录后启用真实直传（banners/drinks/tarot/reviews/avatars/uploads）；留空则用 DataURL 本地预览 */
  dir?: string;
}

export function ImageUpload({ value, onChange, className, dir }: ImageUploadProps) {
  const [preview, setPreview] = useState<string | null>(value || null);
  const [isDragging, setIsDragging] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = useCallback(
    async (file: File) => {
      setError(null);
      const v = validateImageFile(file);
      if (v) { setError(v); return; }

      if (dir) {
        // 真实 OSS 直传
        setUploading(true);
        // 先展示本地 DataURL 作为即时预览
        const reader = new FileReader();
        reader.onload = (e) => setPreview(e.target?.result as string);
        reader.readAsDataURL(file);
        try {
          const url = await uploadToOSS(file, dir);
          setPreview(url);
          onChange(url);
        } catch (e) {
          setError(e instanceof Error ? e.message : '上传失败');
          setPreview(null);
        } finally {
          setUploading(false);
        }
      } else {
        // 兼容旧模式：DataURL 本地预览
        const reader = new FileReader();
        reader.onload = (e) => {
          const result = e.target?.result as string;
          setPreview(result);
          onChange(result);
        };
        reader.readAsDataURL(file);
      }
    },
    [onChange, dir]
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragging(false);
      const file = e.dataTransfer.files[0];
      if (file) handleFile(file);
    },
    [handleFile]
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (file) handleFile(file);
    },
    [handleFile]
  );

  const clearImage = () => {
    setPreview(null);
    setError(null);
    onChange('');
    if (inputRef.current) inputRef.current.value = '';
  };

  if (preview) {
    return (
      <div className={cn('relative w-40 h-40 rounded-xl overflow-hidden group', className)}>
        <img src={preview} alt="Preview" className="w-full h-full object-cover" />
        {uploading && (
          <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
            <Loader2 className="h-6 w-6 text-white animate-spin" />
          </div>
        )}
        <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <button
            onClick={clearImage}
            className="p-2 rounded-lg bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-colors"
          >
            <X className="h-5 w-5" />
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={cn('flex flex-col gap-1', className)}>
      <div
        onClick={() => inputRef.current?.click()}
        onDrop={handleDrop}
        onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
        onDragLeave={() => setIsDragging(false)}
        className={cn(
          'w-40 h-40 rounded-xl border-2 border-dashed border-border bg-background cursor-pointer',
          'flex flex-col items-center justify-center gap-2 transition-all duration-200',
          'hover:border-primary hover:bg-primary/[0.04]',
          isDragging && 'border-gold bg-gold/[0.05]',
        )}
      >
        <ImagePlus className="h-8 w-8 text-text-muted" />
        <span className="text-xs text-text-muted">点击或拖拽上传</span>
        <span className="text-[10px] text-text-muted/70">≤ 5MB</span>
        <input
          ref={inputRef}
          type="file"
          accept="image/*"
          onChange={handleChange}
          className="hidden"
        />
      </div>
      {error && <span className="text-[11px] text-red-400">{error}</span>}
    </div>
  );
}

interface MultiImageUploadProps {
  value: string[];
  onChange: (urls: string[]) => void;
  /** 指定 OSS 目录后启用真实直传；留空则用 DataURL 本地预览 */
  dir?: string;
  /** 最多多少张 */
  max?: number;
}

export function MultiImageUpload({ value, onChange, dir, max = 9 }: MultiImageUploadProps) {
  const [uploading, setUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const handleFile = async (file: File) => {
    setError(null);
    if (value.length >= max) {
      setError(`最多上传 ${max} 张`);
      return;
    }
    const v = validateImageFile(file);
    if (v) { setError(v); return; }

    if (dir) {
      setUploading(true);
      try {
        const url = await uploadToOSS(file, dir);
        onChange([...value, url]);
      } catch (e) {
        setError(e instanceof Error ? e.message : '上传失败');
      } finally {
        setUploading(false);
      }
    } else {
      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result as string;
        onChange([...value, result]);
      };
      reader.readAsDataURL(file);
    }
  };

  const removeImage = (index: number) => {
    onChange(value.filter((_, i) => i !== index));
  };

  return (
    <div className="flex flex-col gap-1">
      <div className="flex flex-wrap gap-3">
        {value.map((url, index) => (
          <div key={index} className="relative w-32 h-32 rounded-xl overflow-hidden group">
            <img src={url} alt={`upload-${index}`} className="w-full h-full object-cover" />
            <div className="absolute inset-0 bg-black/30 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
              <button
                onClick={(e) => { e.stopPropagation(); removeImage(index); }}
                className="p-1.5 rounded-lg bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
          </div>
        ))}
        {value.length < max && (
          <div
            onClick={() => inputRef.current?.click()}
            className={cn(
              'w-32 h-32 rounded-xl border-2 border-dashed border-border bg-background cursor-pointer',
              'flex flex-col items-center justify-center gap-1.5 transition-all',
              'hover:border-primary hover:bg-primary/[0.04]',
              uploading && 'opacity-60 pointer-events-none'
            )}
          >
            {uploading ? (
              <Loader2 className="h-5 w-5 text-text-muted animate-spin" />
            ) : (
              <Upload className="h-6 w-6 text-text-muted" />
            )}
            <span className="text-[11px] text-text-muted">
              {uploading ? '上传中' : `上传图片 (${value.length}/${max})`}
            </span>
            <span className="text-[10px] text-text-muted/70">≤ 5MB</span>
          </div>
        )}
        <input
          ref={inputRef}
          type="file"
          accept="image/*"
          onChange={(e) => { const f = e.target.files?.[0]; if (f) handleFile(f); if (inputRef.current) inputRef.current.value = ''; }}
          className="hidden"
        />
      </div>
      {error && <span className="text-[11px] text-red-400">{error}</span>}
    </div>
  );
}
