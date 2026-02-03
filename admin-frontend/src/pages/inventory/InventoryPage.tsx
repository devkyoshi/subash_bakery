import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { inventoryService } from "@/services/inventory.service";
import type { StockLevel } from "@/types/inventory.types";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Search,
  Package,
  TrendingDown,
  TrendingUp,
  AlertTriangle,
} from "lucide-react";

export function InventoryPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [stockLevels, setStockLevels] = useState<StockLevel[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetchStockLevels();
  }, [user?.organization_id]);

  const fetchStockLevels = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockLevels({
        search,
        organization_id: user.organization_id,
      });
      setStockLevels(response.data.data?.data || []);
    } catch (error) {
      console.error("Failed to fetch stock levels", error);
    } finally {
      setLoading(false);
    }
  };

  const getStockStatus = (stock: StockLevel) => {
    if (stock.quantity_on_hand === 0) {
      return {
        label: "Out of Stock",
        variant: "destructive" as const,
        icon: AlertTriangle,
      };
    }
    if (stock.quantity_available < 10) {
      return {
        label: "Low Stock",
        variant: "secondary" as const,
        icon: TrendingDown,
      };
    }
    return { label: "In Stock", variant: "default" as const, icon: TrendingUp };
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Inventory</h2>
          <p className="text-muted-foreground">
            View and manage stock levels across all locations
          </p>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => navigate("/app/inventory/movements")}
          >
            View Movements
          </Button>
          <Button onClick={() => navigate("/app/inventory/adjustments/new")}>
            Stock Adjustment
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Stock Levels</CardTitle>
          <CardDescription>
            Current inventory levels for all products
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-6">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search products..."
                className="pl-8"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && fetchStockLevels()}
              />
            </div>
            <Button variant="outline" onClick={fetchStockLevels}>
              Search
            </Button>
          </div>

          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Product</TableHead>
                  <TableHead>Location</TableHead>
                  <TableHead className="text-right">On Hand</TableHead>
                  <TableHead className="text-right">Available</TableHead>
                  <TableHead className="text-right">Allocated</TableHead>
                  <TableHead className="text-right">In Transit</TableHead>
                  <TableHead className="text-right">Value</TableHead>
                  <TableHead>Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={8} className="h-24 text-center">
                      Loading inventory...
                    </TableCell>
                  </TableRow>
                ) : stockLevels.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={8} className="h-24 text-center">
                      No stock levels found
                    </TableCell>
                  </TableRow>
                ) : (
                  stockLevels.map((stock) => {
                    const status = getStockStatus(stock);
                    const StatusIcon = status.icon;
                    return (
                      <TableRow key={stock.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center gap-2">
                            <Package className="h-4 w-4 text-muted-foreground" />
                            <div>
                              <div>
                                {stock.product_name ||
                                  `Product ${stock.product_id.slice(-6)}`}
                              </div>
                              {stock.sku && (
                                <div className="text-xs text-muted-foreground">
                                  {stock.sku}
                                </div>
                              )}
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <span className="text-sm">
                            {stock.location_name ||
                              `Location ${stock.location_id.slice(-6)}`}
                          </span>
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {stock.quantity_on_hand.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right">
                          {stock.quantity_available.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right">
                          {stock.quantity_allocated.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right">
                          {stock.quantity_in_transit.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-right">
                          ${stock.total_value.toFixed(2)}
                        </TableCell>
                        <TableCell>
                          <Badge variant={status.variant} className="gap-1">
                            <StatusIcon className="h-3 w-3" />
                            {status.label}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    );
                  })
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
