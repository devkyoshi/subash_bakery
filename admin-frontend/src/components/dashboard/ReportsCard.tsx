import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import type { ReportItem } from "@/types/dashboard";
import { BarChart3, FileSpreadsheet, FileText, ArrowRight, Boxes, PackageSearch, AlertTriangle } from "lucide-react";
import { useNavigate } from "react-router-dom";

function toneIcon(tone: ReportItem["tone"]) {
  switch (tone) {
    case "success":
      return Boxes;
    case "brand":
      return PackageSearch;
    case "warning":
      return AlertTriangle;
    default:
      return BarChart3;
  }
}

function toneBg(tone: ReportItem["tone"]) {
  switch (tone) {
    case "success":
      return "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400";
    case "brand":
      return "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400";
    case "warning":
      return "bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400";
    default:
      return "bg-info text-info-foreground";
  }
}

export function ReportsCard({ items }: { items: ReportItem[] }) {
  const navigate = useNavigate();

  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-start justify-between">
        <div>
          <div className="text-sm font-semibold">Reports Management</div>
          <div className="mt-1 text-xs text-muted-foreground">View and export reports</div>
        </div>

        <div className="grid h-10 w-10 place-items-center rounded-xl bg-brand text-brand-foreground">
          <FileText className="h-5 w-5" />
        </div>
      </div>

      <div className="mt-5 space-y-3">
        {items.map((item) => {
          const Icon = toneIcon(item.tone);
          return (
            <div
              key={item.id}
              onClick={() => item.route && navigate(item.route)}
              className={cn(
                "flex items-center justify-between rounded-2xl border bg-background px-4 py-3 transition-colors",
                item.route
                  ? "cursor-pointer hover:border-brand/45 hover:bg-muted/50"
                  : "opacity-60",
              )}
            >
              <div className="flex items-center gap-3">
                <div className={cn("grid h-10 w-10 place-items-center rounded-xl", toneBg(item.tone))}>
                  <Icon className="h-5 w-5" />
                </div>
                <div>
                  <div className="text-sm font-medium">{item.title}</div>
                  <div className="text-xs text-muted-foreground">{item.subtitle}</div>
                </div>
              </div>
              {item.route && (
                <ArrowRight className="h-4 w-4 text-muted-foreground" />
              )}
            </div>
          );
        })}
      </div>
    </Card>
  );
}
