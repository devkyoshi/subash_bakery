export type Metric = {
  id: string;
  label: string;
  value: string;
  deltaLabel: string;
  deltaVariant: "up" | "down";
  tone: "success" | "info" | "warning" | "brand";
};

export type ReportItem = {
  id: string;
  title: string;
  subtitle: string;
  tone: "info" | "success" | "brand";
};

export type OrderItem = {
  id: string;
  customer: string;
  amount: string;
  status: "Pending" | "Processing" | "Shipped";
  flag: "red" | "orange" | "green";
};

export type ActivityItem = {
  id: string;
  title: string;
  description: string;
  time: string;
  tone: "info" | "success" | "brand" | "warning";
};

export type SalesPoint = {
  month: string;
  invoiced: number;
  cashed: number;
  cashedPct: number;
};
