import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { AlertTriangle, CheckCircle } from 'lucide-react';

interface ConfirmDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  title: string;
  description: string;
  type?: 'danger' | 'warning' | 'success';
  confirmText?: string;
  cancelText?: string;
  onConfirm: () => void;
}

export function ConfirmDialog({
  open,
  onOpenChange,
  title,
  description,
  type = 'warning',
  confirmText = '确认',
  cancelText = '取消',
  onConfirm,
}: ConfirmDialogProps) {
  const icons = {
    danger: AlertTriangle,
    warning: AlertTriangle,
    success: CheckCircle,
  };
  const colors = {
    danger: 'text-red-500',
    warning: 'text-amber-500',
    success: 'text-green-500',
  };
  const btnColors = {
    danger: 'bg-red-500 hover:bg-red-600',
    warning: 'bg-amber-500 hover:bg-amber-600',
    success: 'bg-green-500 hover:bg-green-600',
  };

  const Icon = icons[type];

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[420px] bg-surface-secondary border-gray-100">
        <DialogHeader className="flex flex-col items-center text-center gap-4">
          <Icon className={`h-12 w-12 ${colors[type]}`} />
          <DialogTitle className="text-xl text-text-primary">{title}</DialogTitle>
          <DialogDescription className="text-text-secondary">{description}</DialogDescription>
        </DialogHeader>
        <div className="flex justify-center gap-3 mt-6">
          <Button
            variant="outline"
            onClick={() => onOpenChange(false)}
            className="min-w-[100px] border-border bg-transparent text-text-primary hover:bg-surface-elevated hover:text-text-primary"
          >
            {cancelText}
          </Button>
          <Button
            onClick={() => { onConfirm(); onOpenChange(false); }}
            className={`min-w-[100px] text-white ${btnColors[type]}`}
          >
            {confirmText}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
