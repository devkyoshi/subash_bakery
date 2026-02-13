import React from "react";
import {
  Download,
  FileSpreadsheet,
  Printer,
  RefreshCw,
  Calendar as CalendarIcon,
  Filter,
  MoreVertical,
  ArrowUpRight,
  ArrowDownRight,
  CheckCircle2,
  AlertCircle,
  Clock,
  FileText,
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
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";

// Mock Data
const VARIANCE_DATA = [
  { name: "Matched", value: 74, color: "#f97316" }, // Orange-500
  { name: "Shortage", value: 18, color: "#eab308" }, // Yellow-500
  { name: "Excess", value: 8, color: "#3b82f6" }, // Blue-500
  { name: "Pending", value: 8, color: "#ef4444" }, // Red-500
];

const PO_DATA = [
  {
    id: "PO-2023-0854",
    date: "2023-10-24",
    supplier: "Golden Flour Mill",
    item: "FLR-001",
    itemName: "Premium All-Purpose Flour",
    poQty: 500.0,
    grnQty: 450.0,
    variance: -50.0,
    status: "PARTIAL",
  },
  {
    id: "PO-2023-0856",
    date: "2023-10-23",
    supplier: "Daily Dairy Ltd",
    item: "BTR-042",
    itemName: "Unsalted Butter 25kg",
    poQty: 20.0,
    grnQty: 20.0,
    variance: 0.0,
    status: "MATCHED",
  },
  {
    id: "PO-2023-0858",
    date: "2023-10-22",
    supplier: "Sugar Co Supply",
    item: "SGR-012",
    itemName: "Fine Granulated Sugar",
    poQty: 100.0,
    grnQty: 105.0,
    variance: 5.0,
    status: "EXCESS",
  },
  {
    id: "PO-2023-0860",
    date: "2023-10-21",
    supplier: "Packaging Pros",
    item: "BOX-X01",
    itemName: "Corrugated Cake Box Lrg",
    poQty: 1000.0,
    grnQty: 980.0,
    variance: -20.0,
    status: "PARTIAL",
  },
];

const ACTION_ITEMS = [
  {
    id: 1,
    type: "critical",
    title: "High Discrepancy",
    description:
      "Golden Flour Mill: -50.00 Variance. Follow up regarding shortage on #PO-0854.",
    action: "Mark as Resolved",
  },
  {
    id: 2,
    type: "warning",
    title: "Review Required",
    description:
      "Sugar Co Supply: +5.00 Variance. Review excess delivery for PO #0858.",
    action: "Review GRN",
  },
  {
    id: 3,
    type: "info",
    title: "System Suggestion",
    description:
      "Packaging Pros: -20.00 Variance. Update safety stock for BOX-X01 due to partials.",
    action: "Auto-Adjust",
  },
];

const METRICS = [
  {
    title: "Total POs",
    value: "128",
    change: "+5.2%",
    icon: FileText,
  },
  {
    title: "Completed",
    value: "85",
    subtext: "66.4% of total",
    icon: CheckCircle2,
    color: "text-green-600 dark:text-green-400",
  },
  {
    title: "Partial",
    value: "32",
    change: "+12%",
    icon: AlertCircle,
    color: "text-amber-600 dark:text-amber-400",
  },
  {
    title: "Pending",
    value: "11",
    change: "-8%",
    icon: Clock,
    color: "text-red-600 dark:text-red-400",
  },
  {
    title: "Total Variance",
    value: "$4,250.00",
    change: "1.8% of Total Spend",
    icon: ArrowUpRight,
    color: "text-green-600 dark:text-green-400",
  },
];

export const POvsGRNPage: React.FC = () => {
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
          >
            <Download className="h-4 w-4" /> PDF
          </Button>
          <Button
            variant="outline"
            className="gap-2 text-green-600 dark:text-green-400 border-green-200 dark:border-green-800 hover:bg-green-50 dark:hover:bg-green-900/20"
          >
            <FileSpreadsheet className="h-4 w-4" /> Excel
          </Button>
          <Button variant="outline" className="gap-2">
            <Printer className="h-4 w-4" /> Print
          </Button>
          <Button className="gap-2 bg-primary hover:bg-primary/90 text-primary-foreground shadow">
            <RefreshCw className="h-4 w-4" /> Refresh Data
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
