import { getDashboardData } from "@/services/dashboard.service";
import { StatCard } from "@/components/dashboard/StatCard";
import { ReportsCard } from "@/components/dashboard/ReportsCard";
import { OrderCard } from "@/components/dashboard/OrderCard";
import { SalesHistoryCard } from "@/components/dashboard/SalesHistoryCard";
import { ActivityCard } from "@/components/dashboard/ActivityCard";
import { CalendarCard } from "@/components/dashboard/CalendarCard";

export function DashboardPage() {
  const data = getDashboardData();

  return (
    <div className="space-y-6">
      {/* Stats Grid - Responsive: 1 col mobile, 2 cols tablet, 4 cols desktop */}
      <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4 lg:gap-6">
        {data.metrics.map((m) => (
          <StatCard key={m.id} metric={m} />
        ))}
      </section>

      {/* Reports & Orders - Responsive: 1 col mobile, 2 cols tablet+ */}
      <section className="grid grid-cols-1 gap-4 lg:grid-cols-2 lg:gap-6">
        <ReportsCard items={data.reports} />
        <OrderCard items={data.orders} />
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
