import { cn } from "@/lib/utils";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Metric } from "@/types/dashboard";
import { DollarSign, ShoppingCart, UserRound, ClipboardList } from "lucide-react";

const iconById: Record<string, React.ComponentType<{ className?: string }>> = {
  revenue: DollarSign,
  orders: ShoppingCart,
  users: UserRound,
  tasks: ClipboardList,
};

function toneClasses(tone: Metric["tone"]) {
  switch (tone) {
    case "success":
      return "bg-success text-success-foreground";
    case "info":
      return "bg-info text-info-foreground";
    case "warning":
      return "bg-warning text-warning-foreground";
    default:
      return "bg-brand text-brand-foreground";
  }
}

export function StatCard({ metric }: { metric: Metric }) {
  const Icon = iconById[metric.id] ?? DollarSign;
  const deltaOk = metric.deltaVariant === "up";

  return (
    <Card className="rounded-[16px] bg-elevated text-elevated-foreground p-5 shadow-none">
      <div className="flex items-start justify-between">
        <div className={cn("grid h-10 w-10 place-items-center rounded-xl", toneClasses(metric.tone))}>
          <Icon className="h-5 w-5" />
        </div>
        <Badge
          variant="secondary"
          className={cn(
            "rounded-full px-2.5 py-1 text-xs",
            deltaOk ? "bg-success/15 text-success" : "bg-destructive/15 text-destructive",
          )}
        >
          {metric.deltaLabel}
        </Badge>
      </div>

      <div className="mt-4">
        <div className="text-xs text-muted-foreground">{metric.label}</div>
        <div className="mt-1 text-2xl font-semibold tracking-tight">{metric.value}</div>
      </div>
    </Card>
  );
}
