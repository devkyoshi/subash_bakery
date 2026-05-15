import { useEffect, useState } from "react";
import {
  getDashboardData,
  type DashboardData,
} from "@/services/dashboard.service";
import { StatCard } from "@/components/dashboard/StatCard";
import { ReportsCard } from "@/components/dashboard/ReportsCard";
import { OrderCard } from "@/components/dashboard/OrderCard";
import { SalesHistoryCard } from "@/components/dashboard/SalesHistoryCard";
import { ActivityCard } from "@/components/dashboard/ActivityCard";
import { CalendarCard } from "@/components/dashboard/CalendarCard";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { InventoryAlerts } from "@/components/dashboard/InventoryAlerts";

import { useNavigate } from "react-router-dom";

export function DashboardPage() {
  const navigate = useNavigate();
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadData = async () => {
      try {
        const dashboardData = await getDashboardData();
        setData(dashboardData);
      } catch (error) {
        console.error("Failed to load dashboard data", error);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 lg:gap-6">
          {[1, 2, 3, 4].map((i) => (
            <Card key={i}>
              <CardContent className="p-6">
                <Skeleton className="h-4 w-[100px] mb-4" />
                <Skeleton className="h-8 w-[60px]" />
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    );
  }

  // Handlers
  const handleApprove = async (id: string, mongoId: string) => {
    try {
      const { procurementService } =
        await import("@/services/procurement.service");
      await procurementService.approvePurchaseOrder(mongoId);
      // Reload data
      const dashboardData = await getDashboardData();
      setData(dashboardData);
    } catch (error) {
      console.error("Failed to approve order", error);
    }
  };

  const handleReject = async (id: string, mongoId: string) => {
    try {
      const { procurementService } =
        await import("@/services/procurement.service");
      const { POStatus } = await import("@/types/procurement.types");
      await procurementService.updatePOStatus(mongoId, POStatus.Cancelled);
      // Reload data
      const dashboardData = await getDashboardData();
      setData(dashboardData);
    } catch (error) {
      console.error("Failed to reject order", error);
    }
  };

  if (!data) return <div>Failed to load data</div>;

  return (
    <div className="space-y-6">
      {/* Stats Grid - Responsive: 1 col mobile, 2 cols tablet, 4 cols desktop */}
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 lg:gap-6">
        {data.metrics.map((m) => (
          <StatCard key={m.id} metric={m} />
        ))}
      </section>

      {/* Reports & Orders & Alerts - Responsive: 1 col mobile, 3 cols desktop */}
      <section className="grid grid-cols-1 gap-4 lg:grid-cols-3 lg:gap-6">
        <InventoryAlerts items={data.inventoryAlerts || []} />
        <ReportsCard items={data.reports} />
        {/* Reuse OrderCard for Pending Approvals */}
        <OrderCard
          items={data.orders}
          title="Pending Approvals"
          onNewOrder={() => navigate("/app/procurement/orders/new")}
          onApprove={(id) => {
            // Find the mongoId from the order item
            const order = data.orders.find((o) => o.id === id);
            if (order && order.mongoId) {
              handleApprove(id, order.mongoId);
            }
          }}
          onReject={(id) => {
            const order = data.orders.find((o) => o.id === id);
            if (order && order.mongoId) {
              handleReject(id, order.mongoId);
            }
          }}
        />
      </section>

      {/* Sales History - Full width */}
      <section>
        <SalesHistoryCard data={data.sales} />
      </section>

      {/* Activity & Calendar - Responsive: 1 col mobile, 2 cols tablet+ */}
      <section className="grid grid-cols-1 gap-4 lg:grid-cols-2 lg:gap-6">
        <ActivityCard items={data.activities} />
        <CalendarCard />
      </section>
    </div>
  );
}
