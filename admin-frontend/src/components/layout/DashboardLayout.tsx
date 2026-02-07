import { useState } from "react";
import { Outlet } from "react-router-dom";
import { DashboardSidebar } from "@/components/layout/DashboardSidebar";
import { DashboardTopbar } from "@/components/layout/DashboardTopbar";

export function DashboardLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <div className="min-h-screen bg-app text-app-foreground">
      <DashboardSidebar
        isOpen={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        isCollapsed={isCollapsed}
        toggleCollapse={() => setIsCollapsed(!isCollapsed)}
      />

      <div
        className={`min-h-screen transition-all duration-300 ${
          isCollapsed ? "lg:pl-[80px]" : "lg:pl-[280px]"
        }`}
      >
        <DashboardTopbar onMenuClick={() => setSidebarOpen(true)} />
        <main className="mx-auto w-full max-w-[1180px] px-4 py-6 md:px-6 md:py-8 lg:px-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
