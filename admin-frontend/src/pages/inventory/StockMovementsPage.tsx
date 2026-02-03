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
  Search,
  ArrowRight,
  ArrowLeft,
  RefreshCw,
  Package,
  CalendarDays,
} from "lucide-react";
import { format } from "date-fns";

export function StockMovementsPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetchMovements();
  }, [user?.organization_id]);

  const fetchMovements = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockMovements(
        user.organization_id,
      );
      setMovements(response.data.data || []);
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

      <Card>
        <CardHeader>
          <CardTitle>All Movements</CardTitle>
          <CardDescription>
            A complete history of stock movements
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-6">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search movements..."
                className="pl-8"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
            <Button variant="outline" onClick={fetchMovements}>
              Search
            </Button>
          </div>

          <div className="rounded-md border">
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
                    const typeInfo = getMovementTypeBadge(
                      movement.movement_type,
                    );
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
                          Product {movement.product_id.slice(-6)}
                        </TableCell>
                        <TableCell>
                          {movement.from_location_id
                            ? `Loc ${movement.from_location_id.slice(-6)}`
                            : "-"}
                        </TableCell>
                        <TableCell>
                          {movement.to_location_id
                            ? `Loc ${movement.to_location_id.slice(-6)}`
                            : "-"}
                        </TableCell>
                        <TableCell className="text-right font-medium">
                          {movement.quantity.toFixed(2)} {movement.uom}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground">
                          {movement.reference_no ||
                            movement.reference_type ||
                            "-"}
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
