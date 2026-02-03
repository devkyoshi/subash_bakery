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
  Search,
  MoreHorizontal,
  FileText,
  CalendarDays,
  Eye,
  Edit,
  Trash2,
  Plus,
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

export function StockAdjustmentsPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();
  const [adjustments, setAdjustments] = useState<StockAdjustment[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetchAdjustments();
  }, [user?.organization_id]);

  const fetchAdjustments = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await inventoryService.getStockAdjustments(
        user.organization_id,
      );
      setAdjustments(response.data.data || []);
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

      <Card>
        <CardHeader>
          <CardTitle>All Adjustments</CardTitle>
          <CardDescription>A list of all stock adjustments</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-6">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search adjustments..."
                className="pl-8"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </div>
            <Button variant="outline" onClick={fetchAdjustments}>
              Search
            </Button>
          </div>

          <div className="rounded-md border">
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
                        <Badge
                          variant={getStatusBadgeVariant(adjustment.status)}
                        >
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
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
