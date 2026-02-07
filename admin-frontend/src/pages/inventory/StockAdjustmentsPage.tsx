import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { inventoryService } from "@/services/inventory.service";
import type {
  StockAdjustment,
  AdjustmentStatus,
} from "@/types/inventory.types";
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
  MoreHorizontal,
  FileText,
  CalendarDays,
  Eye,
  Edit,
  Trash2,
  Plus,
  Filter,
  X,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { format } from "date-fns";
import { useToast } from "@/components/ui/use-toast";
import { locationService } from "@/services/location.service";
import { Location } from "@/types/product.types";

export function StockAdjustmentsPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();
  const [adjustments, setAdjustments] = useState<StockAdjustment[]>([]);
  const [loading, setLoading] = useState(true);
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [locationFilter, setLocationFilter] = useState<string>("all");
  const [locations, setLocations] = useState<Location[]>([]);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;

  useEffect(() => {
    fetchLocations();
  }, [user?.organization_id]);

  useEffect(() => {
    fetchAdjustments();
  }, [user?.organization_id, statusFilter, locationFilter, page]);

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

  const fetchAdjustments = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockAdjustments(
        user.organization_id,
        {
          status: statusFilter !== "all" ? statusFilter : undefined,
          location_id: locationFilter !== "all" ? locationFilter : undefined,
          page,
          limit,
        },
      );
      setAdjustments(response.data.data.data || []);
      // @ts-ignore - types might be mismatching but runtime is flat
      setTotal(response.data.data.pagination?.total || 0);
    } catch (error) {
      console.error("Failed to fetch adjustments", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Are you sure you want to delete this adjustment?")) return;
    try {
      await inventoryService.deleteStockAdjustment(id);
      toast({
        title: "Success",
        description: "Adjustment deleted successfully",
      });
      fetchAdjustments();
    } catch (error) {
      console.error("Failed to delete adjustment", error);
      toast({
        title: "Error",
        description: "Failed to delete adjustment",
        variant: "destructive",
      });
    }
  };

  const getStatusBadgeVariant = (status: AdjustmentStatus) => {
    switch (status) {
      case "draft":
        return "secondary";
      case "pending":
        return "outline";
      case "approved":
        return "default";
      case "rejected":
        return "destructive";
      default:
        return "outline";
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            Stock Adjustments
          </h2>
          <p className="text-muted-foreground">
            Manage inventory adjustments and corrections
          </p>
        </div>
        <Button onClick={() => navigate("/app/inventory/adjustments/new")}>
          <Plus className="mr-2 h-4 w-4" /> New Adjustment
        </Button>
      </div>

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <Select
                value={statusFilter}
                onValueChange={(val) => {
                  setStatusFilter(val);
                }}
              >
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="draft">Draft</SelectItem>
                  <SelectItem value="pending">Pending</SelectItem>
                  <SelectItem value="approved">Approved</SelectItem>
                  <SelectItem value="rejected">Rejected</SelectItem>
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

              {(statusFilter !== "all" || locationFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setStatusFilter("all");
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

      {/* Adjustments Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Adjustment No</TableHead>
              <TableHead>Date</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Reason</TableHead>
              <TableHead>Items</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="w-[50px]"></TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  Loading adjustments...
                </TableCell>
              </TableRow>
            ) : adjustments.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  No adjustments found
                </TableCell>
              </TableRow>
            ) : (
              adjustments.map((adjustment) => (
                <TableRow key={adjustment.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      <FileText className="h-4 w-4 text-muted-foreground" />
                      {adjustment.adjustment_no || adjustment.id.slice(-6)}
                    </div>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <CalendarDays className="h-4 w-4 text-muted-foreground" />
                      {format(
                        new Date(adjustment.adjustment_date),
                        "MMM dd, yyyy",
                      )}
                    </div>
                  </TableCell>
                  <TableCell>
                    {adjustment.location_name || adjustment.location_id}
                  </TableCell>
                  <TableCell>{adjustment.reason}</TableCell>
                  <TableCell>{adjustment.items.length}</TableCell>
                  <TableCell>
                    <Badge variant={getStatusBadgeVariant(adjustment.status)}>
                      {adjustment.status.toUpperCase()}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <span className="sr-only">Open menu</span>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>Actions</DropdownMenuLabel>
                        <DropdownMenuItem
                          onClick={() =>
                            navigate(
                              `/app/inventory/adjustments/${adjustment.id}`,
                            )
                          }
                        >
                          <Eye className="mr-2 h-4 w-4" /> View Details
                        </DropdownMenuItem>
                        {adjustment.status === "draft" && (
                          <>
                            <DropdownMenuItem
                              onClick={() =>
                                navigate(
                                  `/app/inventory/adjustments/${adjustment.id}/edit`,
                                )
                              }
                            >
                              <Edit className="mr-2 h-4 w-4" /> Edit
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => handleDelete(adjustment.id)}
                              className="text-destructive"
                            >
                              <Trash2 className="mr-2 h-4 w-4" /> Delete
                            </DropdownMenuItem>
                          </>
                        )}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {!loading && adjustments.length > 0 && (
          <div className="flex items-center justify-between border-t border-border px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} adjustments
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
