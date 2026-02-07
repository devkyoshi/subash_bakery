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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Search,
  MoreHorizontal,
  FileCheck,
  CalendarDays,
  Truck,
  Eye,
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

export function GRNListPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [grns, setGrns] = useState<GoodsReceiptNote[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const limit = 10;

  useEffect(() => {
    fetchGRNs();
  }, [user?.organization_id, statusFilter, page]);

  const fetchGRNs = async () => {
    if (!user?.organization_id) return;
    try {
      setLoading(true);
      const response = await procurementService.getGRNs(user.organization_id, {
        search,
        status:
          statusFilter !== "all" ? (statusFilter as GRNStatus) : undefined,
        page,
        limit,
      });
      setGrns(response.data.data || []);
      setTotal(response.data.pagination?.total || 0);
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

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search GRNs..."
                  className="h-10 pl-10"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && fetchGRNs()}
                />
              </div>
              <Button variant="secondary" onClick={fetchGRNs}>
                Search
              </Button>
            </div>

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
                  <SelectItem value="received">Received</SelectItem>
                  <SelectItem value="inspected">Inspected</SelectItem>
                  <SelectItem value="accepted">Accepted</SelectItem>
                  <SelectItem value="rejected">Rejected</SelectItem>
                </SelectContent>
              </Select>

              {(search || statusFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setSearch("");
                    setStatusFilter("all");
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

      {/* GRN Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>GRN Number</TableHead>
              <TableHead>PO Number</TableHead>
              <TableHead>Supplier</TableHead>
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
                  No GRNs found. Receive goods from a Purchase Order to create
                  one.
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
                  <TableCell>{grn.supplier_name || "-"}</TableCell>
                  <TableCell>
                    <div>
                      <div className="font-medium">
                        {grn.received_by_name || grn.received_by}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {format(new Date(grn.receipt_date), "MMM dd, yyyy")}
                      </div>
                    </div>
                  </TableCell>
                  <TableCell>
                    <Badge variant={getStatusBadgeVariant(grn.status)}>
                      {grn.status?.toUpperCase()}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    {grn.qc_status ? (
                      <Badge variant="outline">
                        {grn.qc_status?.toUpperCase()}
                      </Badge>
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
                        <DropdownMenuItem onClick={() => navigate(`${grn.id}`)}>
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

        {!loading && grns.length > 0 && (
          <div className="flex items-center justify-between border-t px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} GRNs
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
