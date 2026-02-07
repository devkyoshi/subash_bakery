import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { procurementService } from "@/services/procurement.service";
import { PurchaseOrder, POStatus, Supplier } from "@/types/procurement.types";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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
  ArrowLeft,
  CalendarDays,
  FileText,
  Building2,
  CheckCircle2,
  XCircle,
  Send,
  PackageCheck,
} from "lucide-react";
import { format } from "date-fns";
import { Separator } from "@/components/ui/separator";
import { useAuth } from "@/contexts/AuthContext";
import { useToast } from "@/components/ui/use-toast";
import { formatCurrency } from "@/lib/utils";

export function PurchaseOrderDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();

  const [order, setOrder] = useState<PurchaseOrder | null>(null);
  const [supplier, setSupplier] = useState<Supplier | null>(null);
  const [loading, setLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    if (id) {
      fetchOrder(id);
    }
  }, [id]);

  const fetchOrder = async (orderId: string) => {
    try {
      setLoading(true);
      const data = await procurementService.getPurchaseOrder(orderId);
      setOrder(data);

      // Fetch supplier details if we have supplier_id
      if (data.supplier_id) {
        fetchSupplier(data.supplier_id);
      }
    } catch (error) {
      console.error("Failed to fetch order details", error);
      toast({
        title: "Error",
        description: "Failed to load purchase order details",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const fetchSupplier = async (supplierId: string) => {
    try {
      const data = await procurementService.getSupplier(supplierId);
      setSupplier(data);
    } catch (error) {
      console.error("Failed to fetch supplier", error);
    }
  };

  const handleStatusUpdate = async (newStatus: POStatus) => {
    if (!order?.id) return;
    try {
      setActionLoading(true);
      await procurementService.updatePOStatus(order.id, newStatus);
      toast({
        title: "Success",
        description: `Order status updated to ${newStatus}`,
      });
      fetchOrder(order.id); // Refresh data
    } catch (error) {
      console.error("Failed to update status", error);
      toast({
        title: "Error",
        description: "Failed to update order status",
        variant: "destructive",
      });
    } finally {
      setActionLoading(false);
    }
  };

  const handleApprove = async () => {
    if (!order?.id) return;
    try {
      setActionLoading(true);
      await procurementService.approvePurchaseOrder(order.id);
      toast({
        title: "Success",
        description: "Purchase order approved successfully",
      });
      fetchOrder(order.id);
    } catch (error) {
      console.error("Failed to approve order", error);
      toast({
        title: "Error",
        description: "Failed to approve order",
        variant: "destructive",
      });
    } finally {
      setActionLoading(false);
    }
  };

  const getStatusBadgeVariant = (status: POStatus) => {
    switch (status) {
      case POStatus.Draft:
        return "secondary";
      case POStatus.Sent:
        return "outline";
      case POStatus.Confirmed:
        return "default";
      case POStatus.Received:
        return "default";
      case POStatus.PartiallyReceived:
        return "secondary";
      case POStatus.Cancelled:
        return "destructive";
      default:
        return "outline";
    }
  };

  if (loading) {
    return <div className="p-8 text-center">Loading order details...</div>;
  }

  if (!order) {
    return (
      <div className="p-8 text-center">
        <h3 className="text-lg font-medium text-destructive">
          Purchase Order not found
        </h3>
        <Button
          variant="outline"
          onClick={() => navigate("/app/procurement/orders")}
          className="mt-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Orders
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/app/procurement/orders")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">
              {order.po_number || "Order Details"}
            </h2>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant={getStatusBadgeVariant(order.status)}>
                {order.status?.toUpperCase()}
              </Badge>
              <span className="text-sm text-muted-foreground flex items-center gap-1 ml-2">
                <CalendarDays className="h-3 w-3" />
                {format(new Date(order.order_date), "PPP")}
              </span>
            </div>
          </div>
        </div>

        <div className="flex gap-2">
          {order.status === POStatus.Draft && (
            <>
              <Button
                variant="outline"
                onClick={() => navigate(`/app/procurement/orders/${id}/edit`)}
              >
                Edit Order
              </Button>
              <Button onClick={handleApprove} disabled={actionLoading}>
                <CheckCircle2 className="mr-2 h-4 w-4" /> Approve
              </Button>
            </>
          )}

          {order.status === POStatus.Confirmed && (
            <Button
              variant="outline"
              onClick={() => handleStatusUpdate(POStatus.Sent)}
              disabled={actionLoading}
            >
              <Send className="mr-2 h-4 w-4" /> Mark as Sent
            </Button>
          )}

          {order.status === POStatus.Sent && (
            <Button
              onClick={() =>
                navigate(`/app/procurement/grn/new?po_id=${order.id}`)
              }
              disabled={actionLoading}
            >
              <PackageCheck className="mr-2 h-4 w-4" /> Receive Goods
            </Button>
          )}

          {(order.status === POStatus.Draft ||
            order.status === POStatus.Sent) && (
            <Button
              variant="destructive"
              onClick={() => handleStatusUpdate(POStatus.Cancelled)}
              disabled={actionLoading}
            >
              <XCircle className="mr-2 h-4 w-4" /> Cancel
            </Button>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Items</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product</TableHead>
                    <TableHead className="text-right">Quantity</TableHead>
                    <TableHead className="text-right">Unit Price</TableHead>
                    <TableHead className="text-right">Total</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {order.items.map((item, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <div className="font-medium">
                          {/* Ideally we fetch product name or robust PO object has it. 
                               Using product_id fallback if needed. */}
                          Product {item.product_id?.slice(-6)}
                        </div>
                        {item.description && (
                          <div className="text-xs text-muted-foreground">
                            {item.description}
                          </div>
                        )}
                        {item.sku && (
                          <div className="text-xs text-muted-foreground">
                            SKU: {item.sku}
                          </div>
                        )}
                      </TableCell>
                      <TableCell className="text-right">
                        {item.quantity}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(item.unit_price)}
                      </TableCell>
                      <TableCell className="text-right">
                        {formatCurrency(
                          item.line_total || item.quantity * item.unit_price,
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                  <TableRow>
                    <TableCell colSpan={3} className="text-right font-medium">
                      Subtotal
                    </TableCell>
                    <TableCell className="text-right font-medium">
                      {formatCurrency(order.total_amount)}
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {order.notes && (
            <Card>
              <CardHeader>
                <CardTitle>Notes</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {order.notes}
                </p>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Vendor Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {supplier ? (
                <>
                  <div className="font-medium flex items-center gap-2">
                    <Building2 className="h-4 w-4 text-muted-foreground" />
                    {supplier.company_name}
                  </div>
                  <div className="text-sm text-muted-foreground space-y-1">
                    <p>{supplier.contact_person}</p>
                    <p>{supplier.email}</p>
                    <p>{supplier.phone}</p>
                  </div>
                  {supplier.address && (
                    <>
                      <Separator />
                      <div className="text-sm text-muted-foreground mt-2">
                        <p>{supplier.address.street}</p>
                        <p>
                          {supplier.address.city}, {supplier.address.state}
                        </p>
                        <p>{supplier.address.country}</p>
                      </div>
                    </>
                  )}
                  <Button
                    variant="link"
                    className="p-0 h-auto mt-2"
                    onClick={() =>
                      navigate(`/app/procurement/suppliers/${supplier.id}`)
                    }
                  >
                    View Vendor Profile
                  </Button>
                </>
              ) : (
                <div className="text-sm text-muted-foreground">
                  Loading vendor details...
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Order Summary</CardTitle>
            </CardHeader>
            <CardContent className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Reference</span>
                <span className="font-medium">
                  {order.reference_number || "-"}
                </span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-muted-foreground">Expected Date</span>
                <span className="font-medium">
                  {order.expected_date
                    ? format(new Date(order.expected_date), "PPP")
                    : "-"}
                </span>
              </div>
              <Separator className="my-2" />
              <div className="flex justify-between font-medium">
                <span>Total</span>
                <span>{formatCurrency(order.total_amount)}</span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
