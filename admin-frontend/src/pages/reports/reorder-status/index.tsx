import React from "react";
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

// Mock Data
const REORDER_ITEMS = [
  {
    id: "FLR-001",
    name: "All-purpose Flour",
    unit: "KG",
    priority: "CRITICAL",
    currentStock: 45,
    minLevel: 100,
    remainingDays: 3,
    pending: "200 KG",
    sugQty: 500,
    leadTime: "5 Days",
  },
  {
    id: "SUG-042",
    name: "Granulated Sugar",
    unit: "KG",
    priority: "WARNING",
    currentStock: 120,
    minLevel: 150,
    remainingDays: 8,
    pending: "—",
    sugQty: 300,
    leadTime: "3 Days",
  },
  {
    id: "YST-099",
    name: "Dry Yeast (Active)",
    unit: "PKT",
    priority: "NORMAL",
    currentStock: 15,
    minLevel: 10,
    remainingDays: 25,
    pending: "—",
    sugQty: 20,
    leadTime: "7 Days",
  },
];

const CONSUMPTION_DATA = [
  {
    category: "Dry Goods",
    avgDaily: "142.5 KG",
    trend: "+ 12%",
    trendDir: "up",
    forecast: "4,275 KG",
  },
  {
    category: "Dairy",
    avgDaily: "88.2 L",
    trend: "-0.5%",
    trendDir: "neutral",
    forecast: "2,646 L",
  },
  {
    category: "Fermentation",
    avgDaily: "15.8 PKT",
    trend: "↓ 4%",
    trendDir: "down",
    forecast: "474 PKT",
  },
];

const AUTO_REQUISITIONS = [
  { id: "REQ-2023-089", vendor: "Global Grains Co.", items: "3 Items" },
  { id: "REQ-2023-090", vendor: "Dairy Fresh Ltd.", items: "5 Items" },
  { id: "REQ-2023-091", vendor: "Sweetener Hub", items: "1 Item" },
];

export const ReorderStatusPage: React.FC = () => {
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
            <Select defaultValue="main">
              <SelectTrigger className="w-[200px] h-9">
                <SelectValue placeholder="Location" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="main">Main Central Warehouse</SelectItem>
                <SelectItem value="north">North Branch</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Category
            </label>
            <Select defaultValue="all">
              <SelectTrigger className="w-[180px] h-9">
                <SelectValue placeholder="Category" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Categories</SelectItem>
                <SelectItem value="dry">Dry Goods</SelectItem>
                <SelectItem value="dairy">Dairy</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-1">
            <label className="text-xs font-medium text-muted-foreground">
              Urgency Level
            </label>
            <Select defaultValue="all">
              <SelectTrigger className="w-[180px] h-9">
                <SelectValue placeholder="Urgency" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Statuses</SelectItem>
                <SelectItem value="critical">Critical</SelectItem>
                <SelectItem value="warning">Warning</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="flex-1" />

          <div className="flex items-center gap-3 pb-1">
            <div className="flex items-center gap-2">
              <Switch
                id="pending-orders"
                className="data-[state=checked]:bg-primary"
              />
              <label
                htmlFor="pending-orders"
                className="text-sm font-medium text-foreground cursor-pointer"
              >
                Include pending orders
              </label>
            </div>
            <Button className="h-9 bg-primary hover:bg-primary/90 text-primary-foreground shadow">
              Apply Filters
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Critical Stock */}
        <Card className="border shadow-none bg-card">
          <CardContent className="p-4 flex justify-between items-start">
            <div>
              <p className="text-sm font-medium text-muted-foreground">
                Critical stock
              </p>
              <div className="flex items-baseline gap-2 mt-1">
                <span className="text-3xl font-bold text-foreground">12</span>
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
                <span className="text-3xl font-bold text-foreground">28</span>
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
                <span className="text-3xl font-bold text-foreground">145</span>
                <span className="text-sm text-muted-foreground">Items</span>
              </div>
            </div>
            <div className="h-5 w-5 rounded-full bg-green-100 flex items-center justify-center">
              <CheckCircle2 className="h-4 w-4 text-green-600" />
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Reorder Requirements Table */}
      <Card className="border shadow-none bg-card">
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-semibold text-foreground">
            Reorder Requirements
          </CardTitle>
        </CardHeader>
        <CardContent>
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
              {REORDER_ITEMS.map((item) => (
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
              {CONSUMPTION_DATA.map((row, i) => (
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
        </CardContent>
      </Card>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        {/* Recommended Actions */}
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

        {/* Auto-Generated Requisitions */}
        <Card className="border shadow-none bg-card h-full flex flex-col">
          <CardHeader className="py-4 border-b flex flex-row items-center justify-between space-y-0">
            <div className="flex items-center gap-2">
              <FileText className="h-4 w-4 text-green-600 dark:text-green-500" />
              <CardTitle className="text-base font-semibold text-foreground">
                Auto-Generated Requisitions
              </CardTitle>
            </div>
            <span className="text-xs text-primary font-medium cursor-pointer hover:underline">
              View All
            </span>
          </CardHeader>
          <CardContent className="p-0 flex-1 flex flex-col">
            <Table>
              <TableHeader className="bg-transparent">
                <TableRow className="border-b-0 hover:bg-transparent">
                  <TableHead className="text-xs font-semibold text-muted-foreground h-9">
                    Req ID
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground h-9">
                    Vendor
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground h-9">
                    Items
                  </TableHead>
                  <TableHead className="text-xs font-semibold text-muted-foreground text-right h-9">
                    Action
                  </TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {AUTO_REQUISITIONS.map((req, i) => (
                  <TableRow key={i} className="border-b hover:bg-muted/50">
                    <TableCell className="text-sm font-medium text-foreground py-3">
                      {req.id}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground py-3">
                      {req.vendor}
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground py-3">
                      {req.items}
                    </TableCell>
                    <TableCell className="text-right py-3">
                      <span className="text-xs font-bold text-primary hover:underline cursor-pointer">
                        Review
                      </span>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            <div className="p-4 mt-auto">
              <Button className="w-full bg-primary hover:bg-primary/90 text-primary-foreground shadow">
                Approve All Pending Requisitions
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default ReorderStatusPage;
