import React, { useEffect, useState, useCallback } from "react";
import {
  Download,
  Mail,
  Filter,
  ArrowUpRight,
  ArrowDownRight,
  CheckCircle2,
  AlertCircle,
  AlertTriangle,
  Calendar as CalendarIcon,
  TrendingUp,
  TrendingDown,
  Info,
  Lightbulb,
  FileText,
  Clock,
  ArrowRight,
  Search,
  Loader2,
  Package,
} from "lucide-react";
import { format } from "date-fns";
import {
  Card,
  CardContent,
  CardDescription,
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
import { Switch } from "@/components/ui/switch";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useToast } from "@/hooks/use-toast";
import reportService from "@/services/report.service";
import { categoryService } from "@/services/category.service";
import { locationService } from "@/services/location.service";
import type {
  ReorderItem,
  ConsumptionRow,
  ReorderMetrics,
  ReorderStatusFilters,
} from "@/types/report.types";
import type { Category } from "@/types/category.types";

interface LocationOption {
  id: string;
  name: string;
}

export const ReorderStatusPage: React.FC = () => {
  const { toast } = useToast();
  const orgId = localStorage.getItem("organizationId") || "";

  // Data state
  const [items, setItems] = useState<ReorderItem[]>([]);
  const [metrics, setMetrics] = useState<ReorderMetrics | null>(null);
  const [consumptionData, setConsumptionData] = useState<ConsumptionRow[]>([]);
  const [loading, setLoading] = useState(false);

  // Pagination
  const [page, setPage] = useState(1);
  const [limit] = useState(20);
  const [totalPages, setTotalPages] = useState(1);
  const [totalItems, setTotalItems] = useState(0);

  // Filters
  const [search, setSearch] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("all");
  const [locationFilter, setLocationFilter] = useState("all");
  const [priorityFilter, setPriorityFilter] = useState("all");
  const [includePending, setIncludePending] = useState(false);

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
            locResponse.map((l: any) => ({ id: l.id || l._id, name: l.name })),
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
      const params: ReorderStatusFilters & { page?: number; limit?: number } = {
        page,
        limit,
      };

      if (categoryFilter !== "all") params.category_id = categoryFilter;
      if (locationFilter !== "all") params.location_id = locationFilter;
      if (priorityFilter !== "all") params.priority = priorityFilter;
      if (search.trim()) params.search = search.trim();
      if (includePending) params.include_pending = true;

      const response = await reportService.getReorderStatusReport(orgId, params);
      const reportData = response.data.data;

      setItems(reportData.items || []);
      setMetrics(reportData.metrics);
      setConsumptionData(reportData.consumption_data || []);
      setTotalItems(Number(response.data.pagination.total));
      setTotalPages(response.data.pagination.total_pages);
    } catch (err: any) {
      toast({
        title: "Error",
        description:
          err?.response?.data?.error?.message ||
          "Failed to load reorder status report",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  }, [orgId, page, limit, categoryFilter, locationFilter, priorityFilter, search, includePending, toast]);

  useEffect(() => {
    fetchReport();
  }, [fetchReport]);

  const handleApplyFilters = () => {
    setPage(1);
    fetchReport();
  };

  return (
    <div className="space-y-6 pt-6 pb-12 w-full max-w-[1600px] mx-auto px-4 sm:px-6">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 ">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground border-l-4 border-primary pl-3">
            Reorder Status Report Detailed
          </h1>
          <p className="text-muted-foreground mt-1 pl-4 flex items-center gap-2 text-sm">
            <CalendarIcon className="h-4 w-4" />
            Report Generated {format(new Date(), "dd MMM yyyy, hh:mm a")}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" className="gap-2 bg-transparent">
            <Filter className="h-4 w-4 text-muted-foreground" />
          </Button>
          <Button variant="outline" className="gap-2">
            <Mail className="h-4 w-4" /> Email Report
          </Button>
          <Button className="gap-2 bg-primary hover:bg-primary/90 text-primary-foreground shadow">
            <Download className="h-4 w-4" /> Export CSV
          </Button>
        </div>
      </div>

      {/* Filters */}
      <Card className="border shadow-none bg-card">
        <CardContent className="p-4 flex flex-wrap gap-4 items-end">
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
              Urgency Level
            </label>
            <Select value={priorityFilter} onValueChange={setPriorityFilter}>
              <SelectTrigger className="w-[180px] h-9">
                <SelectValue placeholder="All Statuses" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="CRITICAL">Critical</SelectItem>
                <SelectItem value="WARNING">Warning</SelectItem>
                <SelectItem value="NORMAL">Normal</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex-1" />

          <div className="flex items-center gap-3 pb-1">
            <div className="flex items-center gap-2">
              <Switch
                id="pending-orders"
                className="data-[state=checked]:bg-primary"
                checked={includePending}
                onCheckedChange={setIncludePending}
              />
              <label
                htmlFor="pending-orders"
                className="text-sm font-medium text-foreground cursor-pointer"
              >
                Include pending orders
              </label>
            </div>
            <Button
              className="h-9 bg-primary hover:bg-primary/90 text-primary-foreground shadow"
              onClick={handleApplyFilters}
            >
              Apply Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Metrics Cards */}
      {metrics && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {/* Critical Stock */}
          <Card className="border shadow-none bg-card">
            <CardContent className="p-4 flex justify-between items-start">
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Critical stock
                </p>
                <div className="flex items-baseline gap-2 mt-1">
                  <span className="text-3xl font-bold text-foreground">{metrics.critical_count}</span>
                  <span className="text-sm text-muted-foreground">Items</span>
                </div>
              </div>
              <div className="h-5 w-5 rounded-full bg-red-100 flex items-center justify-center">
                <AlertCircle className="h-4 w-4 text-red-600" />
              </div>
            </CardContent>
          </Card>

          {/* Warning Level */}
          <Card className="border shadow-none bg-card">
            <CardContent className="p-4 flex justify-between items-start">
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Warning level
                </p>
                <div className="flex items-baseline gap-2 mt-1">
                  <span className="text-3xl font-bold text-foreground">{metrics.warning_count}</span>
                  <span className="text-sm text-muted-foreground">Items</span>
                </div>
              </div>
              <div className="h-5 w-5 rounded-full bg-amber-100 flex items-center justify-center">
                <AlertTriangle className="h-3 w-3 text-amber-600" />
              </div>
            </CardContent>
          </Card>

          {/* Normal Stock */}
          <Card className="border shadow-none bg-card">
            <CardContent className="p-4 flex justify-between items-start">
              <div>
                <p className="text-sm font-medium text-muted-foreground">
                  Normal stock
                </p>
                <div className="flex items-baseline gap-2 mt-1">
                  <span className="text-3xl font-bold text-foreground">{metrics.normal_count}</span>
                  <span className="text-sm text-muted-foreground">Items</span>
                </div>
              </div>
              <div className="h-5 w-5 rounded-full bg-green-100 flex items-center justify-center">
                <CheckCircle2 className="h-4 w-4 text-green-600" />
              </div>
            </CardContent>
          </Card>
        </div>
      )}

      {/* Reorder Requirements Table */}
      <Card className="border shadow-none bg-card">
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-semibold text-foreground">
            Reorder Requirements
          </CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <div className="flex items-center justify-center py-20">
              <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
              <Package className="h-10 w-10 mb-2 opacity-40" />
              <p className="text-sm">No reorder items found</p>
            </div>
          ) : (
            <>
              <Table>
                <TableHeader className="bg-muted/50">
                  <TableRow>
                    <TableHead className="w-[100px] text-xs uppercase font-semibold text-muted-foreground">
                      Priority
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Item
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Current / Reorder
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Remaining
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Pending
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Sug. Qty
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground">
                      Lead Time
                    </TableHead>
                    <TableHead className="text-xs uppercase font-semibold text-muted-foreground text-right">
                      Action
                    </TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {items.map((item) => (
                    <TableRow key={item.id} className="hover:bg-muted/50">
                      <TableCell>
                        <Badge
                          variant="outline"
                          className={`
                            border-0 px-2 py-0.5 text-xs font-semibold
                            ${item.priority === "CRITICAL" ? "bg-red-50 text-red-700 dark:bg-red-900/30 dark:text-red-400" : ""}
                            ${item.priority === "WARNING" ? "bg-amber-50 text-amber-700 dark:bg-amber-900/30 dark:text-amber-400" : ""}
                            ${item.priority === "NORMAL" ? "bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-400" : ""}
                          `}
                        >
                          {item.priority === "CRITICAL" && (
                            <AlertCircle className="h-3 w-3 mr-1" />
                          )}
                          {item.priority === "WARNING" && (
                            <AlertTriangle className="h-3 w-3 mr-1" />
                          )}
                          {item.priority === "NORMAL" && (
                            <CheckCircle2 className="h-3 w-3 mr-1" />
                          )}
                          {item.priority}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <div className="font-medium text-sm text-foreground">
                          {item.name}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {item.id} • {item.unit}
                        </div>
                      </TableCell>
                      <TableCell>
                        <div
                          className={`font-bold text-sm ${item.priority === "CRITICAL" ? "text-red-600 dark:text-red-400" : item.priority === "WARNING" ? "text-amber-600 dark:text-amber-400" : "text-green-600 dark:text-green-400"}`}
                        >
                          {item.currentStock}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Level: {item.minLevel}
                        </div>
                      </TableCell>
                      <TableCell
                        className={`text-sm font-medium ${item.remainingDays <= 3 ? "text-red-600 dark:text-red-400" : item.remainingDays <= 7 ? "text-amber-600 dark:text-amber-400" : "text-green-600 dark:text-green-400"}`}
                      >
                        {item.remainingDays} Days{item.remainingDays > 20 && "+"}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {item.pending}
                      </TableCell>
                      <TableCell className="text-sm font-medium text-primary">
                        {item.sugQty}
                      </TableCell>
                      <TableCell className="text-sm text-muted-foreground">
                        {item.leadTime}
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="secondary"
                          size="sm"
                          className="h-8 bg-primary/10 text-primary hover:bg-primary/20 hover:text-primary"
                        >
                          + Create Requisition
                        </Button>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between px-4 py-3 border-t mt-2">
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
            </>
          )}
        </CardContent>
      </Card>

      {/* Consumption Analysis */}
      <Card className="border shadow-none bg-card">
        <CardHeader className="py-4 border-b">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4 text-primary" />
              <CardTitle className="text-base font-semibold text-foreground">
                Consumption Analysis
              </CardTitle>
            </div>
            <span className="text-xs text-muted-foreground">
              Rolling 30-Day Window
            </span>
          </div>
        </CardHeader>
        <CardContent className="p-0">
          {consumptionData.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-muted-foreground">
              <p className="text-sm">No consumption data available</p>
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow className="border-b-0">
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Category
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Avg. Daily Consumption
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground">
                    Trend (30D)
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground text-right">
                    Forecasted Monthly
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {consumptionData.map((row, i) => (
                  <TableRow
                    key={i}
                    className="border-b border-gray-100 dark:border-gray-800 hover:bg-transparent"
                  >
                    <TableCell className="font-medium text-sm text-foreground py-3">
                      {row.category}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground py-3">
                      {row.avgDaily}
                    </TableCell>
                    <TableCell className="py-3">
                      <span
                        className={`text-xs font-medium flex items-center gap-1 ${
                          row.trendDir === "up"
                            ? "text-green-600 dark:text-green-400"
                            : row.trendDir === "down"
                              ? "text-red-600 dark:text-red-400"
                              : "text-muted-foreground"
                        }`}
                      >
                        {row.trendDir === "up" && (
                          <ArrowUpRight className="h-3 w-3" />
                        )}
                        {row.trendDir === "down" && (
                          <ArrowDownRight className="h-3 w-3" />
                        )}
                        {row.trendDir === "neutral" && (
                          <ArrowRight className="h-3 w-3" />
                        )}
                        {row.trend}
                      </span>
                    </TableCell>
                    <TableCell className="text-sm font-bold text-foreground text-right py-3">
                      {row.forecast}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      {/* Recommended Actions - Mock data kept as requested */}
      <div className="grid grid-cols-1 lg:grid-cols-1 gap-4">
        <Card className="border shadow-none bg-card h-full">
          <CardHeader className="py-4 border-b flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              <Lightbulb className="h-4 w-4 text-amber-500" />
              <CardTitle className="text-base font-semibold text-foreground">
                Recommended Actions
              </CardTitle>
            </div>
            <Badge
              variant="secondary"
              className="bg-amber-100 text-amber-700 hover:bg-amber-200 dark:bg-amber-900/30 dark:text-amber-400"
            >
              AI Suggested
            </Badge>
          </CardHeader>
          <CardContent className="p-4 space-y-3">
            <div className="p-3 bg-red-50 dark:bg-red-900/10 rounded-lg flex gap-3 items-start">
              <div className="mt-0.5 p-1.5 bg-red-100 dark:bg-red-900/30 rounded text-red-600 dark:text-red-400">
                <Clock className="h-4 w-4" />
              </div>
              <div>
                <h4 className="text-sm font-semibold text-foreground">
                  Process Orders &lt; 24h
                </h4>
                <p className="text-xs text-muted-foreground mt-1 leading-snug">
                  There are 5 critical items where lead time equals days
                  remaining. Convert requisitions to POs within the next 24
                  hours to prevent stockouts.
                </p>
              </div>
            </div>

            <div className="p-3 bg-amber-50 dark:bg-amber-900/10 rounded-lg flex gap-3 items-start">
              <div className="mt-0.5 p-1.5 bg-amber-100 dark:bg-amber-900/30 rounded text-amber-600 dark:text-amber-400">
                <AlertCircle className="h-4 w-4" />
              </div>
              <div>
                <h4 className="text-sm font-semibold text-foreground">
                  Review Warning Items
                </h4>
                <p className="text-xs text-muted-foreground mt-1 leading-snug">
                  28 items have entered the warning zone. Review these items for
                  bulk discount eligibility before finalizing this week's
                  procurement plan.
                </p>
              </div>
            </div>

            <div className="p-3 bg-primary/5 rounded-lg flex gap-3 items-start">
              <div className="mt-0.5 p-1.5 bg-primary/10 rounded text-primary">
                <TrendingUp className="h-4 w-4" />
              </div>
              <div>
                <h4 className="text-sm font-semibold text-foreground">
                  Long Lead Time Items List
                </h4>
                <p className="text-xs text-muted-foreground mt-1 leading-snug">
                  Identify items with lead times &gt; 7 days (e.g., Specialty
                  Grains). Adjust reorder buffers to accommodate supply chain
                  volatility.
                </p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ReorderStatusPage;
