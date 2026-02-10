import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { inventoryService } from "@/services/inventory.service";
import type { StockMovement, MovementType } from "@/types/inventory.types";
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
  ArrowRight,
  ArrowLeft,
  RefreshCw,
  Package,
  CalendarDays,
  Filter,
  X,
} from "lucide-react";
import { format } from "date-fns";
import { locationService } from "@/services/location.service";
import { Location } from "@/types/product.types";

export function StockMovementsPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [locationFilter, setLocationFilter] = useState<string>("all");
  const [typeFilter, setTypeFilter] = useState<string>("all");
  const [locations, setLocations] = useState<Location[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;

  useEffect(() => {
    fetchLocations();
  }, [user?.organization_id]);

  useEffect(() => {
    fetchMovements();
  }, [user?.organization_id, locationFilter, typeFilter, page]);

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

  const fetchMovements = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockMovements(
        user.organization_id,
        {
          location_id: locationFilter !== "all" ? locationFilter : undefined,
          movement_type: typeFilter !== "all" ? typeFilter : undefined,
          page,
          limit,
        },
      );
      const responseData: any = response.data;
      if (Array.isArray(responseData.data)) {
        setMovements(responseData.data);
        setTotal(responseData.data.length);
      } else {
        setMovements(responseData.data?.data || []);
        setTotal(responseData.data?.pagination?.total || 0);
      }
    } catch (error) {
      console.error("Failed to fetch movements", error);
    } finally {
      setLoading(false);
    }
  };

  const getMovementTypeIcon = (type: MovementType) => {
    switch (type) {
      case "in":
        return ArrowRight;
      case "out":
        return ArrowLeft;
      case "transfer":
        return RefreshCw;
      default:
        return Package;
    }
  };

  const getMovementTypeBadge = (type: MovementType) => {
    switch (type) {
      case "in":
        return { label: "In", variant: "default" as const };
      case "out":
        return { label: "Out", variant: "secondary" as const };
      case "transfer":
        return { label: "Transfer", variant: "outline" as const };
      case "adjustment":
        return { label: "Adjustment", variant: "secondary" as const };
      case "return":
        return { label: "Return", variant: "outline" as const };
      case "scrap":
        return { label: "Scrap", variant: "destructive" as const };
      default:
        return { label: type, variant: "outline" as const };
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Stock Movements</h2>
          <p className="text-muted-foreground">
            View all inventory transactions and transfers
          </p>
        </div>
        <Button onClick={() => navigate("/app/inventory")}>
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Inventory
        </Button>
      </div>

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <Select
                value={typeFilter}
                onValueChange={(val) => {
                  setTypeFilter(val);
                }}
              >
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  <SelectItem value="in">In</SelectItem>
                  <SelectItem value="out">Out</SelectItem>
                  <SelectItem value="transfer">Transfer</SelectItem>
                  <SelectItem value="adjustment">Adjustment</SelectItem>
                  <SelectItem value="return">Return</SelectItem>
                  <SelectItem value="scrap">Scrap</SelectItem>
                </SelectContent>
              </Select>

              <Select
                value={locationFilter}
                onValueChange={(val) => {
                  setLocationFilter(val);
                }}
              >
                <SelectTrigger className="h-10 w-[150px]">
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

              {(typeFilter !== "all" || locationFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setTypeFilter("all");
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

      {/* Movements Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Date</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Product</TableHead>
              <TableHead>From</TableHead>
              <TableHead>To</TableHead>
              <TableHead className="text-right">Quantity</TableHead>
              <TableHead>Reference</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  Loading movements...
                </TableCell>
              </TableRow>
            ) : movements.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  No movements found
                </TableCell>
              </TableRow>
            ) : (
              movements.map((movement) => {
                const typeInfo = getMovementTypeBadge(movement.movement_type);
                const Icon = getMovementTypeIcon(movement.movement_type);
                return (
                  <TableRow key={movement.id}>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <CalendarDays className="h-4 w-4 text-muted-foreground" />
                        {format(
                          new Date(movement.movement_date),
                          "MMM dd, yyyy",
                        )}
                      </div>
                    </TableCell>
                    <TableCell>
                      <Badge variant={typeInfo.variant} className="gap-1">
                        <Icon className="h-3 w-3" />
                        {typeInfo.label}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-medium">
                      {movement.product_name ||
                        `Product ${movement.product_id.slice(-6)}`}
                    </TableCell>
                    <TableCell>
                      {movement.from_location_id
                        ? movement.from_location_name ||
                          `Loc ${movement.from_location_id.slice(-6)}`
                        : "-"}
                    </TableCell>
                    <TableCell>
                      {movement.to_location_id
                        ? movement.to_location_name ||
                          `Loc ${movement.to_location_id.slice(-6)}`
                        : "-"}
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {movement.quantity.toFixed(2)} {movement.uom}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {movement.reference_no || movement.reference_type || "-"}
                    </TableCell>
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>

        {!loading && movements.length > 0 && (
          <div className="flex items-center justify-between border-t border-border px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} movements
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
