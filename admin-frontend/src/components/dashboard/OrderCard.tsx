import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import type { OrderItem } from "@/types/dashboard";
import { cn } from "@/lib/utils";
import { Flag, Plus } from "lucide-react";

function statusClass(status: OrderItem["status"]) {
  switch (status) {
    case "Pending":
      return "text-warning";
    case "Processing":
      return "text-info";
    case "Shipped":
      return "text-success";
  }
}

function flagClass(flag: OrderItem["flag"]) {
  switch (flag) {
    case "green":
      return "text-success";
    case "orange":
      return "text-warning";
    default:
      return "text-destructive";
  }
}

export function OrderCard({
  items,
  title = "Order Fulfillment",
  onApprove,
  onReject,
  onNewOrder,
}: {
  items: OrderItem[];
  title?: string;
  onApprove?: (id: string) => void;
  onReject?: (id: string) => void;
  onNewOrder?: () => void;
}) {
  return (
    <Card className="rounded-[18px] bg-elevated p-6 shadow-none">
      <div className="flex items-start justify-between">
        <div>
          <div className="text-sm font-semibold">{title}</div>
          <div className="mt-1 text-xs text-muted-foreground">
            Manage pending items
          </div>
        </div>
        <Button
          className="h-10 rounded-xl bg-brand px-4 text-brand-foreground hover:bg-brand/90"
          onClick={onNewOrder}
        >
          <Plus className="mr-2 h-4 w-4" />
          New Order
        </Button>
      </div>

      <div className="mt-5 space-y-3">
        {items.length === 0 ? (
          <div className="py-8 text-center text-sm text-muted-foreground">
            No pending items
          </div>
        ) : (
          items.map((o) => (
            <div
              key={o.id}
              className="rounded-2xl border border-border bg-background px-4 py-3"
            >
              <div className="flex items-start justify-between gap-4">
                <div className="flex items-start gap-2">
                  <Flag className={cn("mt-0.5 h-4 w-4", flagClass(o.flag))} />
                  <div>
                    <div className="text-xs font-semibold">{o.id}</div>
                    <div className="text-xs text-muted-foreground">
                      {o.customer}
                    </div>
                    <div className="text-[10px] text-muted-foreground">
                      {o.itemCount} items • {o.date}
                    </div>
                  </div>
                </div>
                <div className="text-right flex flex-col items-end gap-2">
                  <div>
                    <div className="text-xs font-semibold">{o.amount}</div>
                    <div
                      className={cn(
                        "text-xs font-medium",
                        statusClass(o.status),
                      )}
                    >
                      {o.status}
                    </div>
                  </div>
                  {onApprove && onReject && (
                    <div className="flex gap-2">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 px-2 text-xs text-destructive hover:text-destructive hover:bg-destructive/10"
                        onClick={() => onReject(o.id)}
                      >
                        Reject
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-6 px-2 text-xs text-success hover:text-success hover:bg-success/10"
                        onClick={() => onApprove(o.id)}
                      >
                        Approve
                      </Button>
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      <div className="mt-4 flex items-center justify-between text-xs text-muted-foreground">
        <span>Showing {items.length} items</span>
        <button className="text-brand hover:underline">View all →</button>
      </div>
    </Card>
  );
}
