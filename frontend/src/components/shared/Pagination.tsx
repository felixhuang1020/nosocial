import { Button } from '@/components/ui/button';

interface PaginationProps {
  page: number;
  total: number;
  pageSize: number;
  onPageChange: (page: number) => void;
}

export function Pagination({ page, total, pageSize, onPageChange }: PaginationProps) {
  const totalPages = Math.ceil(total / pageSize);

  if (totalPages <= 1) return null;

  return (
    <div className="flex items-center justify-between py-4 px-2">
      <span className="text-sm text-text-muted">
        共 {total} 条，第 {page}/{totalPages} 页
      </span>
      <div className="flex items-center gap-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(page - 1)}
          disabled={page <= 1}
          className="border-border bg-transparent text-text-secondary hover:text-text-primary hover:bg-surface-elevated disabled:opacity-40"
        >
          上一页
        </Button>
        {generatePageNumbers(page, totalPages).map((p) => (
          <Button
            key={p}
            variant={p === page ? 'default' : 'outline'}
            size="sm"
            onClick={() => onPageChange(p)}
            className={
              p === page
                ? ''
                : 'border-border bg-transparent text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
            }
          >
            {p}
          </Button>
        ))}
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(page + 1)}
          disabled={page >= totalPages}
          className="border-border bg-transparent text-text-secondary hover:text-text-primary hover:bg-surface-elevated disabled:opacity-40"
        >
          下一页
        </Button>
      </div>
    </div>
  );
}

function generatePageNumbers(current: number, total: number): number[] {
  if (total <= 5) return Array.from({ length: total }, (_, i) => i + 1);
  const start = Math.max(1, Math.min(current - 2, total - 4));
  const end = Math.min(total, start + 4);
  return Array.from({ length: end - start + 1 }, (_, i) => start + i);
}
