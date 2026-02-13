import React, { useEffect, useState, useCallback } from "react";
import {
  Calendar as CalendarIcon,
  FileText,
  Search,
  Package,
  AlertTriangle,
  TrendingUp,
  ShieldAlert,
  Archive,
  XCircle,
  FileSpreadsheet,
  Loader2,
} from "lucide-react";
import { format } from "date-fns";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
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
import { Input } from "@/components/ui/input";
import { useToast } from "@/hooks/use-toast";
import reportService from "@/services/report.service";
import { categoryService } from "@/services/category.service";
import { locationService } from "@/services/location.service";
import type {
  StockLevelComparisonItem,
  StockLevelMetrics,
  StockLevelFilters,
} from "@/types/report.types";
import type { Category } from "@/types/category.types";

interface LocationOption {
  id: string;
  name: string;
}

const STATUS_CONFIG: Record<
  string,
  { label: string; color: string; bg: string; icon: React.ElementType }
> = {
  OPTIMAL: {
    label: "Optimal",
    color: "text-green-700 dark:text-green-400",
    bg: "bg-green-100 dark:bg-green-900/30",
    icon: TrendingUp,
  },
  LOW: {
    label: "Low Stock",
    color: "text-amber-700 dark:text-amber-400",
    bg: "bg-amber-100 dark:bg-amber-900/30",
    icon: AlertTriangle,
  },
  CRITICAL: {
    label: "Critical",
    color: "text-red-700 dark:text-red-400",
    bg: "bg-red-100 dark:bg-red-900/30",
    icon: ShieldAlert,
  },
  OVERSTOCK: {
    label: "Overstock",
    color: "text-blue-700 dark:text-blue-400",
    bg: "bg-blue-100 dark:bg-blue-900/30",
    icon: Archive,
  },
  OUT_OF_STOCK: {
    label: "Out of Stock",
    color: "text-gray-700 dark:text-gray-400",
    bg: "bg-gray-100 dark:bg-gray-900/30",
    icon: XCircle,
  },
};

export const StockLevelReportPage: React.FC = () => {
  const { toast } = useToast();
  const orgId = localStorage.getItem("organizationId") || "";

  // Data state
  const [items, setItems] = useState<StockLevelComparisonItem[]>([]);
  const [metrics, setMetrics] = useState<StockLevelMetrics | null>(null);
  const [loading, setLoading] = useState(false);
  const [exporting, setExporting] = useState<"pdf" | "excel" | null>(null);

  // Pagination
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

  // Filters
  const [search, setSearch] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("all");
  const [locationFilter, setLocationFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");

  // Filter options
  const [categories, setCategories] = useState<Category[]>([]);
  const [locations, setLocations] = useState<LocationOption[]>([]);

  // Load filter options
  useEffect(() => {
    if (!orgId) return;

    const loadFilterOptions = async () => {
      try {
        const [catResponse, locResponse] = await Promise.all([
          categoryService.getCategories({
            organization_id: orgId,
            limit: 100,
            page: 1,
          }),
          locationService.getOrganizationLocations(orgId, { limit: 100 }),
        ]);

        if (catResponse?.data) {
          setCategories(catResponse.data);
        }

        if (Array.isArray(locResponse)) {
          setLocations(
            locResponse.map((l: any) => ({ id: l.id || l._id, name: l.name }))
          );
        }
      } catch (err) {
        console.error("Failed to load filter options:", err);
      }
    };

    loadFilterOptions();
  }, [orgId]);

  // Load report data
  const fetchReport = useCallback(async () => {
    if (!orgId) return;
    setLoading(true);

    try {
      const params: StockLevelFilters & { page?: number; limit?: number } = {
        page,
        limit,
      };

      if (categoryFilter !== "all") params.category_id = categoryFilter;
      if (locationFilter !== "all") params.location_id = locationFilter;
      if (statusFilter !== "all") params.stock_status = statusFilter;
      if (search.trim()) params.search = search.trim();

      const response = await reportService.getStockLevelReport(orgId, params);
      const reportData = response.data.data;

      setItems(reportData.items || []);
      setMetrics(reportData.metrics);
      setTotalItems(Number(response.data.pagination.total));
      setTotalPages(response.data.pagination.total_pages);
    } catch (err: any) {
      toast({
        title: "Error",
        description:
          err?.response?.data?.error?.message ||
          "Failed to load stock level report",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  }, [orgId, page, limit, categoryFilter, locationFilter, statusFilter, search, toast]);

  useEffect(() => {
    fetchReport();
  }, [fetchReport]);

  // Export handlers
  const handleExportExcel = async () => {
    if (!orgId) return;
    setExporting("excel");
    try {
      const filters: StockLevelFilters = {};
      if (categoryFilter !== "all") filters.category_id = categoryFilter;
      if (locationFilter !== "all") filters.location_id = locationFilter;
      if (statusFilter !== "all") filters.stock_status = statusFilter;
      if (search.trim()) filters.search = search.trim();

      const blob = await reportService.exportStockLevelExcel(orgId, filters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `Stock_Level_Report_${format(new Date(), "yyyyMMdd_HHmmss")}.xlsx`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch {
      toast({ title: "Export failed", variant: "destructive" });
    } finally {
      setExporting(null);
    }
  };

  const handleExportPDF = async () => {
    if (!orgId) return;
    setExporting("pdf");
    try {
      const filters: StockLevelFilters = {};
      if (categoryFilter !== "all") filters.category_id = categoryFilter;
      if (locationFilter !== "all") filters.location_id = locationFilter;
      if (statusFilter !== "all") filters.stock_status = statusFilter;
      if (search.trim()) filters.search = search.trim();

      const blob = await reportService.exportStockLevelPDF(orgId, filters);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `Stock_Level_Report_${format(new Date(), "yyyyMMdd_HHmmss")}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch {
      toast({ title: "Export failed", variant: "destructive" });
    } finally {
      setExporting(null);
    }
  };

  const handleApplyFilters = () => {
    setPage(1);
    fetchReport();
  };

  const kpiCards = metrics
    ? [
        {
          label: "Total Products",
          value: metrics.total_products,
          icon: Package,
          color: "text-blue-600",
          bg: "bg-blue-50 dark:bg-blue-900/10",
          border: "border-blue-200 dark:border-blue-800",
        },
        {
          label: "Optimal",
          value: metrics.optimal_count,
          icon: TrendingUp,
          color: "text-green-600",
          bg: "bg-green-50 dark:bg-green-900/10",
          border: "border-green-200 dark:border-green-800",
        },
        {
          label: "Low Stock",
          value: metrics.low_stock_count,
          icon: AlertTriangle,
          color: "text-amber-600",
          bg: "bg-amber-50 dark:bg-amber-900/10",
          border: "border-amber-200 dark:border-amber-800",
        },
        {
          label: "Critical",
          value: metrics.critical_count,
          icon: ShieldAlert,
          color: "text-red-600",
          bg: "bg-red-50 dark:bg-red-900/10",
          border: "border-red-200 dark:border-red-800",
        },
        {
          label: "Overstock",
          value: metrics.overstock_count,
          icon: Archive,
          color: "text-blue-500",
          bg: "bg-blue-50 dark:bg-blue-900/10",
          border: "border-blue-200 dark:border-blue-800",
        },
        {
          label: "Total Stock Value",
          value: `${metrics.total_stock_value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`,
          icon: Package,
          color: "text-purple-600",
          bg: "bg-purple-50 dark:bg-purple-900/10",
          border: "border-purple-200 dark:border-purple-800",
        },
      ]
    : [];

  return (
    <div className="space-y-6 pt-6 pb-12 w-full max-w-[1600px] mx-auto px-4 sm:px-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground border-l-4 border-primary pl-3">
            Stock Level Comparison Report
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
            disabled={exporting !== null}
          >
            {exporting === "pdf" ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <FileText className="h-4 w-4" />
            )}
            PDF
          </Button>
          <Button
            variant="outline"
            className="gap-2 text-green-600 border-green-200 hover:bg-green-50"
            onClick={handleExportExcel}
            disabled={exporting !== null}
          >
            {exporting === "excel" ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <FileSpreadsheet className="h-4 w-4" />
            )}
            Excel
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card className="border shadow-none bg-card">
        <CardContent className="p-4 flex flex-wrap gap-4 items-end">
          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Search
            </label>
            <div className="relative">
              <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
              <Input
                placeholder="Product name or SKU..."
                className="pl-9 w-[200px] h-9"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleApplyFilters()}
              />
            </div>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Category
            </label>
            <Select value={categoryFilter} onValueChange={setCategoryFilter}>
              <SelectTrigger className="w-[180px] h-9">
                <SelectValue placeholder="All Categories" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                {categories.map((cat) => (
                  <SelectItem key={cat.id} value={cat.id}>
                    {cat.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Location
            </label>
            <Select value={locationFilter} onValueChange={setLocationFilter}>
              <SelectTrigger className="w-[200px] h-9">
                <SelectValue placeholder="All Locations" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Locations</SelectItem>
                {locations.map((loc) => (
                  <SelectItem key={loc.id} value={loc.id}>
                    {loc.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Stock Status
            </label>
            <Select value={statusFilter} onValueChange={setStatusFilter}>
              <SelectTrigger className="w-[180px] h-9">
                <SelectValue placeholder="All Statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="OPTIMAL">Optimal</SelectItem>
                <SelectItem value="LOW">Low Stock</SelectItem>
                <SelectItem value="CRITICAL">Critical</SelectItem>
                <SelectItem value="OVERSTOCK">Overstock</SelectItem>
                <SelectItem value="OUT_OF_STOCK">Out of Stock</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex-1" />
          <Button
            className="h-9 bg-primary hover:bg-primary/90 text-primary-foreground shadow"
            onClick={handleApplyFilters}
          >
            Apply Filters
          </Button>
        </CardContent>
      </Card>

      {/* KPI Cards */}
      {metrics && (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
          {kpiCards.map((kpi, idx) => (
            <Card
              key={idx}
              className={`border shadow-none ${kpi.bg} ${kpi.border}`}
            >
              <CardContent className="p-4">
                <div className="flex items-center gap-2 mb-2">
                  <kpi.icon className={`h-4 w-4 ${kpi.color}`} />
                  <p className="text-xs font-medium text-muted-foreground">
                    {kpi.label}
                  </p>
                </div>
                <p className={`text-xl font-bold ${kpi.color}`}>{kpi.value}</p>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Stock Levels Table */}
      <Card className="border shadow-none bg-card">
        <CardHeader className="flex flex-row items-center justify-between pb-2">
          <CardTitle className="text-base font-semibold text-foreground">
            Stock Levels
          </CardTitle>
          <Badge variant="secondary" className="bg-muted text-muted-foreground">
            {totalItems} products
          </Badge>
        </CardHeader>
        <CardContent className="p-0">
          {loading ? (
            <div className="flex items-center justify-center py-20">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
              <Package className="h-10 w-10 mb-2 opacity-40" />
              <p className="text-sm">No stock level data found</p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <Table className="min-w-[900px]">
                <TableHeader className="bg-muted/30">
                  <TableRow className="hover:bg-transparent border-b-0">
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 w-[80px]">
                      SKU
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10">
                      Product Name
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10">
                      Category
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10">
                      Location
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 w-[60px]">
                      Unit
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-right w-[80px]">
                      On Hand
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-right w-[80px]">
                      Available
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-right w-[80px]">
                      Allocated
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-right w-[80px]">
                      Reorder Lvl
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-right w-[90px]">
                      Total Value
                    </TableHead>
                    <TableHead className="text-xs font-semibold text-muted-foreground h-10 text-center w-[100px]">
                      Status
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {items.map((item, idx) => {
                    const statusCfg =
                      STATUS_CONFIG[item.stock_status] || STATUS_CONFIG.OPTIMAL;
                    const StatusIcon = statusCfg.icon;

                    return (
                      <TableRow
                        key={`${item.product_id}-${item.location_id}-${idx}`}
                        className="hover:bg-muted/50 border-b border-gray-100 dark:border-gray-800"
                      >
                        <TableCell className="font-bold text-xs text-amber-600 dark:text-amber-500 py-3 whitespace-nowrap">
                          {item.sku || "—"}
                        </TableCell>
                        <TableCell className="font-medium text-sm text-foreground py-3">
                          {item.product_name || "—"}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground py-3">
                          {item.category_name}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground py-3">
                          {item.location_name}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground py-3">
                          {item.unit}
                        </TableCell>
                        <TableCell className="text-sm font-medium text-foreground text-right py-3 tabular-nums">
                          {item.system_qty.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-sm font-medium text-foreground text-right py-3 tabular-nums">
                          {item.available_qty.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground text-right py-3 tabular-nums">
                          {item.allocated_qty.toFixed(2)}
                        </TableCell>
                        <TableCell className="text-sm text-muted-foreground text-right py-3 tabular-nums">
                          {item.reorder_level}
                        </TableCell>
                        <TableCell className="text-sm font-medium text-foreground text-right py-3 tabular-nums">
                          {item.total_value.toLocaleString(undefined, {
                            minimumFractionDigits: 2,
                            maximumFractionDigits: 2,
                          })}
                        </TableCell>
                        <TableCell className="text-center py-3">
                          <Badge
                            variant="secondary"
                            className={`${statusCfg.bg} ${statusCfg.color} font-semibold text-xs px-2 py-0.5 gap-1`}
                          >
                            <StatusIcon className="h-3 w-3" />
                            {statusCfg.label}
                          </Badge>
                        </TableCell>
                      </TableRow>
                    );
                  })}
                </TableBody>
              </Table>
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex items-center justify-between px-4 py-3 border-t">
              <p className="text-xs text-muted-foreground">
                Page {page} of {totalPages} ({totalItems} items)
              </p>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page <= 1}
                  onClick={() => setPage((p) => Math.max(1, p - 1))}
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={page >= totalPages}
                  onClick={() => setPage((p) => p + 1)}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

export default StockLevelReportPage;
