import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { ReportItem } from "@/types/dashboard";
import { BarChart3, Download, FileSpreadsheet, FileText, Eye } from "lucide-react";

function toneIcon(tone: ReportItem["tone"]) {
  switch (tone) {
    case "success":
      return FileSpreadsheet;
    case "brand":
      return FileText;
    default:
      return BarChart3;
  }
}

function toneBg(tone: ReportItem["tone"]) {
  switch (tone) {
    case "success":
      return "bg-success text-success-foreground";
    case "brand":
      return "bg-brand text-brand-foreground";
    default:
      return "bg-info text-info-foreground";
  }
}

export function ReportsCard({ items }: { items: ReportItem[] }) {
  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-start justify-between">
        <div>
          <div className="text-sm font-semibold">Reports Management</div>
          <div className="mt-1 text-xs text-muted-foreground">Generate and download reports</div>
        </div>

        <div className="grid h-10 w-10 place-items-center rounded-xl bg-brand text-brand-foreground">
          <FileText className="h-5 w-5" />
        </div>
      </div>

      <div className="mt-5 space-y-3">
        {items.map((item, idx) => {
          const Icon = toneIcon(item.tone);
          const active = idx === 0;
          return (
            <div
              key={item.id}
              className={cn(
                "flex items-center justify-between rounded-2xl border bg-background px-4 py-3",
                active ? "border-brand/45 ring-2 ring-brand/20" : "border-border",
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
              <div className="h-5 w-5 rounded-full border border-border" />
            </div>
          );
        })}
      </div>

      <div className="mt-5 flex items-center gap-3">
        <Button className="h-12 flex-1 rounded-xl bg-brand text-brand-foreground hover:bg-brand/90">
          <Download className="mr-2 h-4 w-4" />
          Download
        </Button>
        <Button
          variant="outline"
          size="icon"
          className="h-12 w-12 rounded-xl bg-elevated"
          aria-label="Preview"
        >
          <Eye className="h-4 w-4" />
        </Button>
      </div>
    </Card>
  );
}
