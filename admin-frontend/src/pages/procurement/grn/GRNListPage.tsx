import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { procurementService } from "@/services/procurement.service";
import { GoodsReceiptNote, GRNStatus } from "@/types/procurement.types";
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
  FileCheck,
  CalendarDays,
  Truck,
  Eye,
} from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { format } from "date-fns";

export function GRNListPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [grns, setGrns] = useState<GoodsReceiptNote[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");

  useEffect(() => {
    fetchGRNs();
  }, [user?.organization_id]);

  const fetchGRNs = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await procurementService.getGRNs(user.organization_id, {
        search,
      });
      setGrns(response.data.data || []);
    } catch (error) {
      console.error("Failed to fetch GRNs", error);
    } finally {
      setLoading(false);
    }
  };

  const getStatusBadgeVariant = (status: GRNStatus) => {
    switch (status) {
      case GRNStatus.Draft:
        return "secondary";
      case GRNStatus.Received:
        return "default";
      case GRNStatus.Inspected:
        return "outline";
      case GRNStatus.Accepted:
        return "default"; // success equivalent
      case GRNStatus.Rejected:
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
            Goods Receipt Notes
          </h2>
          <p className="text-muted-foreground">
            Manage received goods and inspections.
          </p>
        </div>
        {/* GRNs are usually created from a PO, so maybe no direct "Create" button here, 
            or it redirects to PO list to select one? 
            For now, let's keep it clean or add a button that goes to PO list with a hint.
        */}
        <Button onClick={() => navigate("/app/procurement/orders")}>
          <Truck className="mr-2 h-4 w-4" /> Receive from PO
        </Button>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All GRNs</CardTitle>
          <CardDescription>A list of all goods receipt notes.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 mb-6">
            <div className="relative flex-1 max-w-sm">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Search GRNs..."
                className="pl-8"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && fetchGRNs()}
              />
            </div>
            <Button variant="outline" onClick={fetchGRNs}>
              Search
            </Button>
          </div>

          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>GRN Number</TableHead>
                  <TableHead>PO Number</TableHead>
                  <TableHead>Date Received</TableHead>
                  <TableHead>Received By</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>QC</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {loading ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      Loading GRNs...
                    </TableCell>
                  </TableRow>
                ) : grns.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={7} className="h-24 text-center">
                      No GRNs found. Receive goods from a Purchase Order to
                      create one.
                    </TableCell>
                  </TableRow>
                ) : (
                  grns.map((grn) => (
                    <TableRow key={grn.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <FileCheck className="h-4 w-4 text-muted-foreground" />
                          {grn.grn_number}
                        </div>
                      </TableCell>
                      <TableCell>{grn.po_number || "-"}</TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <CalendarDays className="h-4 w-4 text-muted-foreground" />
                          {format(new Date(grn.receipt_date), "MMM dd, yyyy")}
                        </div>
                      </TableCell>
                      <TableCell>{grn.received_by}</TableCell>
                      <TableCell>
                        <Badge variant={getStatusBadgeVariant(grn.status)}>
                          {grn.status?.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {grn.qc_status ? (
                          <Badge variant="outline">{grn.qc_status?.toUpperCase()}</Badge>
                        ) : (
                          "-"
                        )}
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
                              onClick={() => navigate(`${grn.id}`)}
                            >
                              <Eye className="mr-2 h-4 w-4" /> View Details
                            </DropdownMenuItem>
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
