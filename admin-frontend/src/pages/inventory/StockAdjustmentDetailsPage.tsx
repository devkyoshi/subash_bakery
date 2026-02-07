import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { inventoryService } from "@/services/inventory.service";
import type { StockAdjustment } from "@/types/inventory.types";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import { ArrowLeft, CheckCircle2, XCircle, CalendarDays } from "lucide-react";
import { format } from "date-fns";
import { Separator } from "@/components/ui/separator";
import { useToast } from "@/components/ui/use-toast";
import { formatCurrency } from "@/lib/utils";

export function StockAdjustmentDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const [adjustment, setAdjustment] = useState<StockAdjustment | null>(null);
  const [loading, setLoading] = useState(true);
  const [rejectDialogOpen, setRejectDialogOpen] = useState(false);
  const [rejectReason, setRejectReason] = useState("");
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    if (id) {
      fetchAdjustment(id);
    }
  }, [id]);

  const fetchAdjustment = async (adjustmentId: string) => {
    try {
      setLoading(true);
      const response = await inventoryService.getStockAdjustment(adjustmentId);
      setAdjustment(response.data.data);
    } catch (error) {
      console.error("Failed to fetch adjustment", error);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async () => {
    if (!adjustment?.id) return;
    try {
      setActionLoading(true);
      await inventoryService.approveStockAdjustment(adjustment.id);
      toast({
        title: "Success",
        description: "Adjustment approved successfully",
      });
      fetchAdjustment(adjustment.id);
    } catch (error) {
      console.error("Failed to approve adjustment", error);
      toast({
        title: "Error",
        description: "Failed to approve adjustment",
        variant: "destructive",
      });
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async () => {
    if (!adjustment?.id || !rejectReason.trim()) return;
    try {
      setActionLoading(true);
      await inventoryService.rejectStockAdjustment(adjustment.id, rejectReason);
      toast({
        title: "Success",
        description: "Adjustment rejected",
      });
      setRejectDialogOpen(false);
      fetchAdjustment(adjustment.id);
    } catch (error) {
      console.error("Failed to reject adjustment", error);
      toast({
        title: "Error",
        description: "Failed to reject adjustment",
        variant: "destructive",
      });
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return <div className="p-8 text-center">Loading adjustment details...</div>;
  }

  if (!adjustment) {
    return (
      <div className="p-8 text-center">
        <h3 className="text-lg font-medium text-destructive">
          Adjustment not found
        </h3>
        <Button
          variant="outline"
          onClick={() => navigate("/app/inventory/adjustments")}
          className="mt-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Adjustments
        </Button>
      </div>
    );
  }

  const getStatusBadgeVariant = (status: string) => {
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
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/app/inventory/adjustments")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">
              {adjustment.adjustment_no ||
                `Adjustment ${adjustment.id.slice(-6)}`}
            </h2>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant={getStatusBadgeVariant(adjustment.status)}>
                {adjustment.status.toUpperCase()}
              </Badge>
              <span className="text-sm text-muted-foreground flex items-center gap-1 ml-2">
                <CalendarDays className="h-3 w-3" />
                {format(new Date(adjustment.adjustment_date), "PPP")}
              </span>
            </div>
          </div>
        </div>

        <div className="flex gap-2">
          {(adjustment.status === "draft" ||
            adjustment.status === "pending") && (
            <>
              <Dialog
                open={rejectDialogOpen}
                onOpenChange={setRejectDialogOpen}
              >
                <DialogTrigger asChild>
                  <Button variant="outline">
                    <XCircle className="mr-2 h-4 w-4" /> Reject
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Reject Adjustment</DialogTitle>
                    <DialogDescription>
                      Please provide a reason for rejecting this adjustment.
                    </DialogDescription>
                  </DialogHeader>
                  <Textarea
                    placeholder="Rejection reason..."
                    value={rejectReason}
                    onChange={(e) => setRejectReason(e.target.value)}
                  />
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => setRejectDialogOpen(false)}
                    >
                      Cancel
                    </Button>
                    <Button
                      variant="destructive"
                      onClick={handleReject}
                      disabled={!rejectReason.trim() || actionLoading}
                    >
                      Reject
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
              <Button onClick={handleApprove} disabled={actionLoading}>
                <CheckCircle2 className="mr-2 h-4 w-4" /> Approve
              </Button>
            </>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Adjustment Items</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product</TableHead>
                    <TableHead className="text-right">Expected</TableHead>
                    <TableHead className="text-right">Actual</TableHead>
                    <TableHead className="text-right">Difference</TableHead>
                    <TableHead className="text-right">Cost Impact</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {adjustment.items.map((item, index) => (
                    <TableRow key={index}>
                      <TableCell className="font-medium">
                        <div>
                          {item.product_name ||
                            `Product ${item.product_id.slice(-6)}`}
                        </div>
                        {item.sku && (
                          <div className="text-xs text-muted-foreground">
                            {item.sku}
                          </div>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        {item.expected_qty.toFixed(2)}
                      </TableCell>
                      <TableCell className="text-right">
                        {item.actual_qty.toFixed(2)}
                      </TableCell>
                      <TableCell className="text-right">
                        <span
                          className={
                            item.difference_qty > 0
                              ? "text-green-600"
                              : item.difference_qty < 0
                                ? "text-red-600"
                                : ""
                          }
                        >
                          {item.difference_qty > 0 ? "+" : ""}
                          {item.difference_qty.toFixed(2)}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(item.total_cost || 0)}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {adjustment.notes && (
            <Card>
              <CardHeader>
                <CardTitle>Notes</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {adjustment.notes}
                </p>
              </CardContent>
            </Card>
          )}

          {adjustment.rejected_reason && (
            <Card>
              <CardHeader>
                <CardTitle>Rejection Reason</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-destructive whitespace-pre-wrap">
                  {adjustment.rejected_reason}
                </p>
              </CardContent>
            </Card>
          )}
        </div>

        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <span className="text-sm text-muted-foreground">Location</span>
                <div className="font-medium">
                  {adjustment.location_name ||
                    `Location ${adjustment.location_id.slice(-6)}`}
                </div>
              </div>
              <Separator />
              <div>
                <span className="text-sm text-muted-foreground">Reason</span>
                <div className="font-medium">{adjustment.reason}</div>
                {adjustment.reason_details && (
                  <p className="text-sm text-muted-foreground mt-1">
                    {adjustment.reason_details}
                  </p>
                )}
              </div>
              {adjustment.approved_by && (
                <>
                  <Separator />
                  <div>
                    <span className="text-sm text-muted-foreground">
                      Approved By
                    </span>
                    <div className="font-medium">
                      {adjustment.approved_by_name ||
                        adjustment.approved_by.slice(-6)}
                    </div>
                    {adjustment.approved_at && (
                      <p className="text-sm text-muted-foreground">
                        {format(new Date(adjustment.approved_at), "PPP")}
                      </p>
                    )}
                  </div>
                </>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
