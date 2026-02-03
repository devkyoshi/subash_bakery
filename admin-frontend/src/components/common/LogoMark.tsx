import { cn } from "@/lib/utils";

export function LogoMark({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "grid h-24 w-24 place-items-center border border-border",
        className,
      )}
    >
      <span className="text-sm">Logo</span>
    </div>
  );
}
