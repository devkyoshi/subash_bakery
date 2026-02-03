import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { procurementService } from "@/services/procurement.service";
import { PurchaseOrder, POStatus } from "@/types/procurement.types";
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
  Plus,
  Search,
  MoreHorizontal,
  FileText,
  CalendarDays,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { format } from "date-fns";

export function PurchaseOrdersPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [orders, setOrders] = useState<PurchaseOrder[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetchOrders();
  }, [user?.organization_id]);

  const fetchOrders = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await procurementService.getPurchaseOrders(
        user.organization_id,
        {
          search,
        },
      );
      setOrders(response.data.data || []);
    } catch (error) {
      console.error("Failed to fetch purchase orders", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadgeVariant = (status: POStatus) => {
    switch (status) {
      case POStatus.Draft:
        return "secondary";
      case POStatus.Sent:
        return "outline";
      case POStatus.Confirmed:
        return "default"; // info/brand color
      case POStatus.Received:
        return "default"; // success color in many themes, or leave default
      case POStatus.PartiallyReceived:
        return "secondary";
      case POStatus.Cancelled:
        return "destructive";
      default:
        return "outline";
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Purchase Orders</h2>
          <p className="text-muted-foreground">
            Create and manage orders to your suppliers.
          </p>
        </div>
        <Button onClick={() => navigate("new")}>
          <Plus className="mr-2 h-4 w-4" /> Create Order
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Orders</CardTitle>
          <CardDescription>A list of all purchase orders.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-6">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search orders..."
                className="pl-8"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && fetchOrders()}
              />
            </div>
            <Button variant="outline" onClick={fetchOrders}>
              Search
            </Button>
          </div>

          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>PO Number</TableHead>
                  <TableHead>Date</TableHead>
                  <TableHead>Supplier</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Items</TableHead>
                  <TableHead className="text-right">Total Amount</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      Loading orders...
                    </TableCell>
                  </TableRow>
                ) : orders.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      No purchase orders found. Create one to get started.
                    </TableCell>
                  </TableRow>
                ) : (
                  orders.map((order) => (
                    <TableRow key={order.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <FileText className="h-4 w-4 text-muted-foreground" />
                          {order.po_number}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <CalendarDays className="h-4 w-4 text-muted-foreground" />
                          {format(new Date(order.order_date), "MMM dd, yyyy")}
                        </div>
                      </TableCell>
                      <TableCell>
                        {/* We might need to fetch supplier name or include it in aggregation. 
                            For now, assuming supplier_id or extra logic. 
                            Ideally, the backend list endpoint should populate supplier name.
                            Checking type definition... it has supplier_id. 
                            We'll display ID for now or 'Loading...' if we don't have name. 
                            IMPROVEMENT: specific fetch or expanded DTO. 
                        */}
                        <span className="font-medium text-muted-foreground">
                          {/* Placeholder until we have supplier name populated */}
                          Supplier {order.supplier_id.slice(-6)}
                        </span>
                      </TableCell>
                      <TableCell>
                        <Badge variant={getStatusBadgeVariant(order.status)}>
                          {order.status?.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell>{order.items.length} items</TableCell>
                      <TableCell className="text-right font-medium">
                        {new Intl.NumberFormat("en-US", {
                          style: "currency",
                          currency: order.currency || "USD",
                        }).format(order.total_amount)}
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
                              onClick={() => navigate(`${order.id}`)}
                            >
                              View Details
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              onClick={() => navigate(`${order.id}/edit`)}
                              disabled={order.status !== POStatus.Draft}
                            >
                              Edit
                            </DropdownMenuItem>
                            <DropdownMenuSeparator />
                            {order.status === POStatus.Draft && (
                              <DropdownMenuItem className="text-destructive">
                                Cancel Order
                              </DropdownMenuItem>
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
