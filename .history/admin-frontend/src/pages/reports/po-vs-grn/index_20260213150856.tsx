import React, { useEffect, useState, useCallback } from "react";
import {
  Download,
  FileSpreadsheet,
  Printer,
  RefreshCw,
  Calendar as CalendarIcon,
  ArrowUpRight,
  CheckCircle2,
  AlertCircle,
  Clock,
  FileText,
  Loader2,
} from "lucide-react";
import { format, subDays } from "date-fns";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
} from "recharts";
import { useToast } from "@/hooks/use-toast";
import reportService from "@/services/report.service";
import { procurementService } from "@/services/procurement.service";
import type {
  POvsGRNReportResponse,
  POvsGRNComparisonItem,
  VarianceDistribution,
  ActionItem,
  POvsGRNMetrics,
  ReportFilters,
} from "@/types/report.types";

export const POvsGRNPage: React.FC = () => {
  const { toast } = useToast();
  const orgId = localStorage.getItem("organizationId") || "";

  // State
  const [loading, setLoading] = useState(true);
  const [exporting, setExporting] = useState<"pdf" | "excel" | null>(null);
  const [reportData, setReportData] = useState<POvsGRNReportResponse | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [limit] = useState(10);

  // Filters
  const [filters, setFilters] = useState<ReportFilters>({
    start_date: format(subDays(new Date(), 30), "yyyy-MM-dd"),
    end_date: format(new Date(), "yyyy-MM-dd"),
  });
  const [supplierFilter, setSupplierFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");

  // Suppliers list for filter dropdown
  const [suppliers, setSuppliers] = useState<{ id: string; name: string }[]>([]);

  // Fetch suppliers for filter dropdown
  useEffect(() => {
    const fetchSuppliers = async () => {
      if (!orgId) return;
      try {
        const response = await procurementService.getSuppliers(orgId, { limit: 100 });
        const supplierList = (response.data?.data || []).map((s: any) => ({
          id: s.id,
          name: s.company_name,
        }));
        setSuppliers(supplierList);
      } catch {
        // Silently fail — suppliers dropdown will just show "All"
      }
    };
    fetchSuppliers();
  }, [orgId]);

  // Fetch report data
  const fetchReport = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);
    try {
      const params: ReportFilters & { page?: number; limit?: number } = {
        ...filters,
        page,
        limit,
      };
      if (supplierFilter !== "all") params.supplier_id = supplierFilter;
      if (statusFilter !== "all") params.status = statusFilter;

      const response = await reportService.getPOvsGRNComparison(orgId, params);
      setReportData(response.data.data as unknown as POvsGRNReportResponse);
      setTotalPages(response.data.pagination.total_pages || 1);
    } catch (err: any) {
      toast({
        title: "Error",
        description: err?.response?.data?.error?.message || "Failed to fetch report data",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  }, [orgId, filters, page, limit, supplierFilter, statusFilter, toast]);

  useEffect(() => { fetchReport(); }, [fetchReport]);

  const handleApplyFilters = () => { setPage(1); fetchReport(); };

  const handleExportExcel = async () => {
    if (!orgId) return;
    setExporting("excel");
    try {
      const exportFilters: ReportFilters = { ...filters };
      if (supplierFilter !== "all") exportFilters.supplier_id = supplierFilter;
      if (statusFilter !== "all") exportFilters.status = statusFilter;
      const blob = await reportService.exportPOvsGRNExcel(orgId, exportFilters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `PO_vs_GRN_Report_${format(new Date(), "yyyyMMdd_HHmmss")}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
      toast({ title: "Success", description: "Excel report downloaded" });
    } catch {
      toast({ title: "Error", description: "Failed to export Excel", variant: "destructive" });
    } finally { setExporting(null); }
  };

  const handleExportPDF = async () => {
    if (!orgId) return;
    setExporting("pdf");
    try {
      const exportFilters: ReportFilters = { ...filters };
      if (supplierFilter !== "all") exportFilters.supplier_id = supplierFilter;
      if (statusFilter !== "all") exportFilters.status = statusFilter;
      const blob = await reportService.exportPOvsGRNPDF(orgId, exportFilters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `PO_vs_GRN_Report_${format(new Date(), "yyyyMMdd_HHmmss")}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
      toast({ title: "Success", description: "PDF report downloaded" });
    } catch {
      toast({ title: "Error", description: "Failed to export PDF", variant: "destructive" });
    } finally { setExporting(null); }
  };

  // Derived data
  const metrics: POvsGRNMetrics = reportData?.metrics || {
    total_pos: 0, completed_pos: 0, partial_pos: 0, pending_pos: 0, excess_pos: 0,
    total_variance: 0, total_po_value: 0, total_grn_value: 0, variance_percent: 0, completed_percent: 0,
  };
  const items: POvsGRNComparisonItem[] = reportData?.items || [];
  const varianceData: VarianceDistribution[] = reportData?.variance_distribution || [];
  const actionItems: ActionItem[] = reportData?.action_items || [];
  const matchedPct = varianceData.find((v) => v.name === "Matched")?.value || 0;

  const METRICS_CARDS = [
    { title: "Total POs", value: metrics.total_pos.toString(), change: "", icon: FileText },
    { title: "Completed", value: metrics.completed_pos.toString(), subtext: `${metrics.completed_percent}% of total`, icon: CheckCircle2, color: "text-green-600 dark:text-green-400" },
    { title: "Partial", value: metrics.partial_pos.toString(), change: "", icon: AlertCircle, color: "text-amber-600 dark:text-amber-400" },
    { title: "Pending", value: metrics.pending_pos.toString(), change: "", icon: Clock, color: "text-red-600 dark:text-red-400" },
    { title: "Total Variance", value: `$${metrics.total_variance.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`, change: `${metrics.variance_percent}% of Total Spend`, icon: ArrowUpRight, color: "text-green-600 dark:text-green-400" },
  ];

  return (
    <div className="space-y-6 pt-6 pb-12 w-full max-w-[1600px] mx-auto px-4 sm:px-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground border-l-4 border-primary pl-3">
            PO vs GRN Comparison
          </h1>
          <p className="text-muted-foreground mt-1 pl-4 flex items-center gap-2 text-sm">
            <CalendarIcon className="h-4 w-4" />
            Report Generated {format(new Date(), "dd MMM yyyy, hh:mm a")}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button
            variant="outline"
            className="gap-2 text-destructive border-destructive/20 hover:bg-destructive/10"
            onClick={handleExportPDF}
            disabled={exporting === "pdf" || loading}
          >
            {exporting === "pdf" ? <Loader2 className="h-4 w-4 animate-spin" /> : <Download className="h-4 w-4" />} PDF
          </Button>
          <Button
            variant="outline"
            className="gap-2 text-green-600 dark:text-green-400 border-green-200 dark:border-green-800 hover:bg-green-50 dark:hover:bg-green-900/20"
            onClick={handleExportExcel}
            disabled={exporting === "excel" || loading}
          >
            {exporting === "excel" ? <Loader2 className="h-4 w-4 animate-spin" /> : <FileSpreadsheet className="h-4 w-4" />} Excel
          </Button>
          <Button variant="outline" className="gap-2" onClick={() => window.print()}>
            <Printer className="h-4 w-4" /> Print
          </Button>
          <Button
            className="gap-2 bg-primary hover:bg-primary/90 text-primary-foreground shadow"
            onClick={fetchReport}
            disabled={loading}
          >
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />} Refresh Data
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card className="border shadow-none bg-card">
        <CardContent className="p-4 flex flex-wrap gap-4 items-center">
          <div className="flex items-center gap-2 bg-muted px-3 py-2 rounded-md border border-border">
            <CalendarIcon className="h-4 w-4 text-muted-foreground" />
            <span className="text-sm font-medium text-foreground">
              Last 30 Days
            </span>
          </div>
          <Select defaultValue="all">
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Supplier" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Supplier: All</SelectItem>
              <SelectItem value="gfm">Golden Flour Mill</SelectItem>
              <SelectItem value="ddl">Daily Dairy Ltd</SelectItem>
            </SelectContent>
          </Select>
          <Select defaultValue="all">
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Status" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">Status: All</SelectItem>
              <SelectItem value="matched">Matched</SelectItem>
              <SelectItem value="partial">Partial</SelectItem>
              <SelectItem value="excess">Excess</SelectItem>
            </SelectContent>
          </Select>
          <Select defaultValue="main">
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Warehouse" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="main">Main Warehouse</SelectItem>
              <SelectItem value="north">North Branch</SelectItem>
            </SelectContent>
          </Select>
          <div className="flex-1" />
          <Button
            variant="secondary"
            className="bg-primary/10 text-primary hover:bg-primary/20 hover:text-primary border border-primary/20"
          >
            Apply Filters
          </Button>
        </CardContent>
      </Card>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-4">
        {METRICS.map((metric, index) => (
          <Card
            key={index}
            className="border shadow-none hover:shadow-md transition-shadow bg-card"
          >
            <CardContent className="p-4">
              <div className="flex justify-between items-start">
                <div>
                  <p className="text-sm font-medium text-muted-foreground">
                    {metric.title}
                  </p>
                  <h3 className="text-2xl font-bold mt-1 text-foreground">
                    {metric.value}
                  </h3>
                  <p
                    className={`text-xs mt-1 ${metric.change?.includes("+") ? "text-green-600 dark:text-green-400" : metric.change?.includes("-") ? "text-red-600 dark:text-red-400" : "text-muted-foreground"}`}
                  >
                    {metric.change || metric.subtext}
                  </p>
                </div>
                <div
                  className={`p-2 rounded-full ${
                    index === 0
                      ? "bg-blue-50 dark:bg-blue-900/20"
                      : index === 1
                        ? "bg-green-50 dark:bg-green-900/20"
                        : index === 2
                          ? "bg-amber-50 dark:bg-amber-900/20"
                          : index === 3
                            ? "bg-red-50 dark:bg-red-900/20"
                            : "bg-muted"
                  }`}
                >
                  <metric.icon
                    className={`h-4 w-4 ${metric.color || "text-muted-foreground"}`}
                  />
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* PO Table */}
        <Card className="lg:col-span-2 border shadow-none bg-card">
          <CardHeader className="pb-2">
            <div className="flex justify-between items-center">
              <CardTitle className="text-lg font-semibold text-foreground">
                PO List
              </CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-0">
            <Table>
              <TableHeader className="bg-muted/50">
                <TableRow>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">
                    PO No
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">
                    Date
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">
                    Supplier
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">
                    Item
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider text-right">
                    PO Qty
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider text-right">
                    GRN Qty
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider text-right">
                    Variance
                  </TableHead>
                  <TableHead className="font-semibold text-xs text-muted-foreground uppercase tracking-wider">
                    Status
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {PO_DATA.map((po) => (
                  <TableRow key={po.id} className="hover:bg-muted/50">
                    <TableCell className="font-medium text-primary">
                      <div>{po.id.replace("PO-", "#PO-")}</div>
                      <div className="text-xs text-muted-foreground">
                        {po.id.split("-").slice(0, 2).join("-")}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground text-sm">
                      {po.date}
                    </TableCell>
                    <TableCell>
                      <div className="font-medium text-foreground text-sm">
                        {po.supplier}
                      </div>
                    </TableCell>
                    <TableCell>
                      <div className="font-medium text-foreground text-sm">
                        {po.item}
                      </div>
                      <div className="text-xs text-primary/80">
                        {po.itemName}
                      </div>
                    </TableCell>
                    <TableCell className="text-right font-medium text-foreground">
                      {po.poQty.toFixed(2)}
                    </TableCell>
                    <TableCell className="text-right font-medium text-foreground">
                      {po.grnQty.toFixed(2)}
                    </TableCell>
                    <TableCell
                      className={`text-right font-bold ${po.variance < 0 ? "text-red-500 dark:text-red-400" : po.variance > 0 ? "text-green-600 dark:text-green-400" : "text-muted-foreground"}`}
                    >
                      {po.variance > 0 ? "+" : ""}
                      {po.variance.toFixed(2)}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={`
                          ${po.status === "MATCHED" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 hover:bg-green-200 dark:hover:bg-green-900/50" : ""}
                          ${po.status === "PARTIAL" ? "bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400 hover:bg-red-200 dark:hover:bg-red-900/50" : ""}
                          ${po.status === "EXCESS" ? "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 hover:bg-green-200 dark:hover:bg-green-900/50" : ""}
                        `}
                      >
                        {po.status}
                      </Badge>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
            <div className="p-4 border-t flex justify-between items-center text-sm text-muted-foreground">
              <span>Page 1 of 12</span>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" className="h-8">
                  Prev
                </Button>
                <Button
                  variant="default"
                  size="sm"
                  className="h-8 bg-primary hover:bg-primary/90"
                >
                  1
                </Button>
                <Button variant="outline" size="sm" className="h-8">
                  2
                </Button>
                <Button variant="outline" size="sm" className="h-8">
                  Next
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <div className="space-y-6">
          {/* Variance Distribution Chart */}
          <Card className="border shadow-none bg-card">
            <CardHeader>
              <CardTitle className="text-lg font-semibold text-foreground">
                Variance Distribution
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="h-[250px] relative">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={VARIANCE_DATA}
                      innerRadius={60}
                      outerRadius={80}
                      paddingAngle={5}
                      dataKey="value"
                      stroke="none"
                    >
                      {VARIANCE_DATA.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip
                      contentStyle={{
                        backgroundColor: "var(--popover)",
                        borderColor: "var(--border)",
                        borderRadius: "var(--radius)",
                        color: "var(--popover-foreground)",
                      }}
                      itemStyle={{ color: "var(--popover-foreground)" }}
                    />
                  </PieChart>
                </ResponsiveContainer>
                <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
                  <span className="text-3xl font-bold text-foreground">
                    74%
                  </span>
                  <span className="text-xs text-muted-foreground font-medium">
                    ACCURACY
                  </span>
                </div>
              </div>
              <div className="mt-4 space-y-2">
                {VARIANCE_DATA.map((item) => (
                  <div
                    key={item.name}
                    className="flex justify-between items-center text-sm"
                  >
                    <div className="flex items-center gap-2">
                      <div
                        className="w-3 h-3 rounded-full"
                        style={{ backgroundColor: item.color }}
                      />
                      <span className="text-muted-foreground">{item.name}</span>
                    </div>
                    <span className="font-bold text-foreground">
                      {item.value}%
                    </span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Action Items */}
          <Card className="border shadow-none bg-card">
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-lg font-semibold text-foreground">
                Action Items
              </CardTitle>
              <Badge
                variant="destructive"
                className="rounded-full h-6 w-6 flex items-center justify-center p-0"
              >
                3
              </Badge>
            </CardHeader>
            <CardContent className="space-y-3">
              {ACTION_ITEMS.map((item) => (
                <div
                  key={item.id}
                  className={`p-3 rounded-lg border-l-4 ${
                    item.type === "critical"
                      ? "bg-red-50 dark:bg-red-900/20 border-red-500"
                      : item.type === "warning"
                        ? "bg-amber-50 dark:bg-amber-900/20 border-amber-500"
                        : "bg-primary/10 border-primary/50"
                  }`}
                >
                  <div className="flex justify-between items-start mb-1">
                    <span
                      className={`text-xs font-bold ${
                        item.type === "critical"
                          ? "text-red-700 dark:text-red-400"
                          : item.type === "warning"
                            ? "text-amber-700 dark:text-amber-400"
                            : "text-primary dark:text-primary"
                      }`}
                    >
                      {item.title}
                    </span>
                  </div>
                  <p className="text-sm text-foreground/80 mb-2 leading-snug">
                    {item.description}
                  </p>
                  <button
                    className={`text-xs font-semibold hover:underline ${
                      item.type === "critical"
                        ? "text-red-600 dark:text-red-400"
                        : item.type === "warning"
                          ? "text-amber-600 dark:text-amber-400"
                          : "text-primary"
                    }`}
                  >
                    {item.action}
                  </button>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
};

export default POvsGRNPage;
