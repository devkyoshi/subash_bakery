import React from "react";
import { AlertTriangle, ArrowRight, Package } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import type { LowStockItem } from "@/services/dashboard.service";

interface InventoryAlertsProps {
  items: LowStockItem[];
  loading?: boolean;
}

export function InventoryAlerts({ items, loading }: InventoryAlertsProps) {
  // Hardcoded threshold for visualization (10 as per backend)
  const THRESHOLD = 10;

  if (loading) {
    return (
      <Card className="h-full border-none shadow-none bg-card/50">
        <CardHeader>
          <div className="h-6 w-32 bg-muted animate-pulse rounded" />
        </CardHeader>
        <CardContent className="space-y-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="flex items-center gap-4">
              <div className="h-10 w-10 rounded-full bg-muted animate-pulse" />
              <div className="space-y-2 flex-1">
                <div className="h-4 w-1/3 bg-muted animate-pulse rounded" />
                <div className="h-3 w-1/2 bg-muted animate-pulse rounded" />
              </div>
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="h-full border shadow-none bg-gradient-to-br from-card to-card/50">
      <CardHeader className="flex flex-row items-center justify-between pb-2">
        <div className="space-y-1">
          <CardTitle className="text-lg font-semibold flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-orange-500" />
            Inventory Alerts
          </CardTitle>
          <p className="text-sm text-muted-foreground">
            Items below reorder level
          </p>
        </div>
        <Badge
          variant="outline"
          className="bg-orange-50 text-orange-600 border-orange-200"
        >
          {items.length} Critical
        </Badge>
      </CardHeader>
      <CardContent className="pt-4">
        {items.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-center text-muted-foreground">
            <Package className="h-12 w-12 mb-3 opacity-20" />
            <p>No inventory alerts</p>
            <p className="text-xs">Stock levels are healthy</p>
          </div>
        ) : (
          <div className="space-y-6">
            {items.map((item) => {
              // Calculate percentage for progress bar (capped at 100%)
              const percentage = Math.min(
                (item.quantity_available / THRESHOLD) * 100,
                100,
              );

              return (
                <div key={item.product_id} className="space-y-3">
                  <div className="flex items-start justify-between">
                    <div className="space-y-0.5">
                      <div className="font-medium text-sm">
                        {item.product_name || "Unknown Product"}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        SKU: {item.sku || "N/A"} • Zone:{" "}
                        {item.warehouse_zone || "N/A"}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-bold text-orange-600">
                        {item.quantity_available}{" "}
                        <span className="text-xs font-normal text-muted-foreground">
                          / {THRESHOLD}
                        </span>
                      </div>
                      <div className="text-xs text-muted-foreground">
                        Available
                      </div>
                    </div>
                  </div>

                  <Progress
                    value={percentage}
                    className="h-1.5 bg-orange-100"
                  />

                  <div className="flex justify-end">
                    <Button
                      variant="link"
                      size="sm"
                      className="h-auto p-0 text-xs text-primary"
                    >
                      Restock <ArrowRight className="h-3 w-3 ml-1" />
                    </Button>
                  </div>
                </div>
              );
            })}
          </div>
        )}

        {items.length > 0 && (
          <Button className="w-full mt-6" variant="outline">
            View All Stock Levels
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
