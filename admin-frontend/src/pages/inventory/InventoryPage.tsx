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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Search,
  Package,
  TrendingDown,
  TrendingUp,
  AlertTriangle,
  Filter,
  X,
} from "lucide-react";
import { formatCurrency } from "@/lib/utils";
import { locationService } from "@/services/location.service";
import { Location } from "@/types/product.types";

export function InventoryPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [stockLevels, setStockLevels] = useState<StockLevel[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [locationFilter, setLocationFilter] = useState<string>("all");
  const [locations, setLocations] = useState<Location[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;

  useEffect(() => {
    fetchLocations();
  }, [user?.organization_id]);

  useEffect(() => {
    fetchStockLevels();
  }, [user?.organization_id, locationFilter, page]);

  const fetchLocations = async () => {
    if (!user?.organization_id) return;
    try {
      const locations = await locationService.getOrganizationLocations(
        user.organization_id,
      );
      setLocations(locations || []);
    } catch (error) {
      console.error("Failed to fetch locations", error);
    }
  };

  const fetchStockLevels = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockLevels({
        search,
        organization_id: user.organization_id,
        location_id: locationFilter !== "all" ? locationFilter : undefined,
        page,
        limit,
      });
      setStockLevels(response.data.data.data || []);
      // @ts-ignore - types might be mismatching but runtime is flat
      setTotal(response.data.data.pagination?.total || 0);
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

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search products..."
                  className="h-10 pl-10"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && fetchStockLevels()}
                />
              </div>
              <Button variant="secondary" onClick={fetchStockLevels}>
                Search
              </Button>
            </div>

            <div className="flex items-center gap-2">
              <Select
                value={locationFilter}
                onValueChange={(val) => {
                  setLocationFilter(val);
                }}
              >
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Location" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Locations</SelectItem>
                  {locations.map((location) => (
                    <SelectItem key={location.id} value={location.id}>
                      {location.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              {(search || locationFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setSearch("");
                    setLocationFilter("all");
                  }}
                >
                  <X className="mr-2 h-4 w-4" />
                  Clear
                </Button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Stock Levels Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
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
                      {formatCurrency(stock.total_value)}
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

        {!loading && stockLevels.length > 0 && (
          <div className="flex items-center justify-between border-t border-border px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} items
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page + 1)}
                disabled={page * limit >= total}
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
