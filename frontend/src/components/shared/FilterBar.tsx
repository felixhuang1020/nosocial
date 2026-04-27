import type { ReactNode } from 'react';
import { Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';

interface FilterBarProps {
  searchPlaceholder?: string;
  searchValue?: string;
  onSearchChange?: (value: string) => void;
  children?: ReactNode;
  className?: string;
}

export function FilterBar({
  searchPlaceholder = '搜索...',
  searchValue,
  onSearchChange,
  children,
  className,
}: FilterBarProps) {
  return (
    <div className={cn('flex flex-wrap items-center gap-3 rounded-xl bg-surface-secondary border border-gray-100 p-4', className)}>
      {onSearchChange && (
        <div className="relative flex-1 min-w-[200px]">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-text-muted pointer-events-none" />
          <Input
            placeholder={searchPlaceholder}
            value={searchValue || ''}
            onChange={(e) => onSearchChange(e.target.value)}
            className="pl-9 bg-background border-border text-text-primary placeholder:text-text-muted focus:border-gold focus:ring-gold/15"
          />
        </div>
      )}
      {children}
    </div>
  );
}

interface FilterTabsProps {
  options: { value: string | number; label: string }[];
  value: string | number;
  onChange: (value: string | number) => void;
  className?: string;
}

export function FilterTabs({ options, value, onChange, className }: FilterTabsProps) {
  return (
    <div className={cn('flex items-center gap-1 rounded-lg bg-background p-1', className)}>
      {options.map((option) => (
        <button
          key={option.value}
          onClick={() => onChange(option.value)}
          className={cn(
            'px-3 py-1.5 rounded-md text-sm font-medium transition-all duration-200',
            value === option.value
              ? 'bg-gold text-white'
              : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
          )}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
}
