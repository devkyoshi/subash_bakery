import type { ActivityItem, Metric, OrderItem, ReportItem, SalesPoint } from "@/types/dashboard";

export type DashboardData = {
  metrics: Metric[];
  reports: ReportItem[];
  orders: OrderItem[];
  activities: ActivityItem[];
  sales: SalesPoint[];
};

export function getDashboardData(): DashboardData {
  return {
    metrics: [
      {
        id: "revenue",
        label: "Total Revenue",
        value: "$45,231",
        deltaLabel: "+12.5%",
        deltaVariant: "up",
        tone: "success",
      },
      {
        id: "orders",
        label: "Total Orders",
        value: "1,234",
        deltaLabel: "+8.2%",
        deltaVariant: "up",
        tone: "info",
      },
      {
        id: "users",
        label: "Active Users",
        value: "892",
        deltaLabel: "+15.3%",
        deltaVariant: "up",
        tone: "brand",
      },
      {
        id: "tasks",
        label: "Pending Tasks",
        value: "23",
        deltaLabel: "-5.1%",
        deltaVariant: "down",
        tone: "warning",
      },
    ],
    reports: [
      { id: "sales", title: "Sales Report", subtitle: "Last updated: Today", tone: "info" },
      { id: "inventory", title: "Inventory Report", subtitle: "Last updated: Today", tone: "success" },
      { id: "financial", title: "Financial Report", subtitle: "Last updated: Today", tone: "brand" },
    ],
    orders: [
      { id: "#ORD-1234", customer: "John Smith", amount: "$245.00", status: "Pending", flag: "red" },
      { id: "#ORD-1235", customer: "Sarah Johnson", amount: "$189.50", status: "Processing", flag: "orange" },
      { id: "#ORD-1236", customer: "Mike Wilson", amount: "$432.00", status: "Pending", flag: "red" },
      { id: "#ORD-1237", customer: "Emily Brown", amount: "$156.75", status: "Shipped", flag: "green" },
    ],
    activities: [
      { id: "a1", title: "New order received", description: "Order #1234 from John Smith", time: "5 minutes ago", tone: "info" },
      { id: "a2", title: "Inventory updated", description: "Product A restocked with 50 units", time: "1 hour ago", tone: "success" },
      { id: "a3", title: "New user registered", description: "Sarah Johnson joined the team", time: "2 hours ago", tone: "brand" },
      { id: "a4", title: "Report generated", description: "Monthly sales report is ready", time: "3 hours ago", tone: "warning" },
      { id: "a5", title: "Payment received", description: "$1,245.00 from Order #1230", time: "4 hours ago", tone: "success" },
    ],
    sales: [
      { month: "Mar 2022", invoiced: 15000, cashed: 14000, cashedPct: 93.3 },
      { month: "Apr 2022", invoiced: 13000, cashed: 12000, cashedPct: 92.3 },
      { month: "May 2022", invoiced: 12000, cashed: 11500, cashedPct: 95.8 },
      { month: "Jun 2022", invoiced: 13500, cashed: 12800, cashedPct: 94.8 },
      { month: "Jul 2022", invoiced: 11000, cashed: 10500, cashedPct: 95.5 },
      { month: "Aug 2022", invoiced: 9000, cashed: 8600, cashedPct: 95.6 },
      { month: "Sep 2022", invoiced: 12000, cashed: 11800, cashedPct: 98.3 },
      { month: "Oct 2022", invoiced: 11800, cashed: 11650, cashedPct: 98.7 },
      { month: "Nov 2022", invoiced: 12300, cashed: 12000, cashedPct: 97.6 },
      { month: "Dec 2022", invoiced: 14000, cashed: 13500, cashedPct: 96.4 },
      { month: "Jan 2023", invoiced: 15000, cashed: 14200, cashedPct: 94.7 },
      { month: "Feb 2023", invoiced: 12014, cashed: 2554, cashedPct: 21.26 },
      { month: "Mar 2023", invoiced: 10000, cashed: 2000, cashedPct: 20 },
    ],
  };
}
