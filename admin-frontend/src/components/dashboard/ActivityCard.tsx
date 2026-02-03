import { Card } from "@/components/ui/card";
import type { ActivityItem } from "@/types/dashboard";
import { cn } from "@/lib/utils";
import { BarChart3, ClipboardList, Package, ShoppingCart, UserRound } from "lucide-react";

function toneBg(tone: ActivityItem["tone"]) {
  switch (tone) {
    case "success":
      return "bg-success text-success-foreground";
    case "brand":
      return "bg-brand text-brand-foreground";
    case "warning":
      return "bg-warning text-warning-foreground";
    default:
      return "bg-info text-info-foreground";
  }
}

const iconByTitle: Record<string, React.ComponentType<{ className?: string }>> = {
  "New order received": ShoppingCart,
  "Inventory updated": Package,
  "New user registered": UserRound,
  "Report generated": BarChart3,
  "Payment received": ClipboardList,
};

export function ActivityCard({ items }: { items: ActivityItem[] }) {
  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-start justify-between">
        <div>
          <div className="text-sm font-semibold">Recent Activities</div>
          <div className="mt-1 text-xs text-muted-foreground">Latest system updates</div>
        </div>
        <button className="text-muted-foreground">⋮</button>
      </div>

      <div className="mt-5 space-y-4">
        {items.map((a) => {
          const Icon = iconByTitle[a.title] ?? ClipboardList;
          return (
            <div key={a.id} className="flex items-start justify-between gap-4">
              <div className="flex items-start gap-3">
                <div className={cn("grid h-10 w-10 place-items-center rounded-xl", toneBg(a.tone))}>
                  <Icon className="h-5 w-5" />
                </div>
                <div>
                  <div className="text-sm font-medium">{a.title}</div>
                  <div className="text-xs text-muted-foreground">{a.description}</div>
                </div>
              </div>
              <div className="text-xs text-muted-foreground">{a.time}</div>
            </div>
          );
        })}
      </div>

      <div className="mt-6 border-t border-border pt-4 text-center">
        <button className="text-xs font-medium text-brand hover:underline">View all activities →</button>
      </div>
    </Card>
  );
}
