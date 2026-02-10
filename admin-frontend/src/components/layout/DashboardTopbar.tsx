import { Bell, Search, Menu } from "lucide-react";
import { ThemeToggle } from "@/components/common/ThemeToggle";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useLocation } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { UserMenu } from "@/components/common/UserMenu";
import { useEffect, useState } from "react";
import {
  requestForToken,
  onMessageListener,
} from "@/services/notification.service";
import { toast } from "sonner";

const TITLE_BY_PREFIX: Array<{ prefix: string; title: string }> = [
  { prefix: "/app/dashboard", title: "Dashboard" },
  { prefix: "/app/reports", title: "Reports" },
  { prefix: "/app/orders", title: "Orders" },
  { prefix: "/app/inventory", title: "Inventory" },
  { prefix: "/app/calendar", title: "Calendar" },
  { prefix: "/app/team", title: "Team" },
  { prefix: "/app/users", title: "Users" },
  { prefix: "/app/transactions", title: "Transactions" },
  { prefix: "/app/analytics", title: "Analytics" },
  { prefix: "/app/settings", title: "Settings" },
];

interface DashboardTopbarProps {
  onMenuClick: () => void;
}

export function DashboardTopbar({ onMenuClick }: DashboardTopbarProps) {
  const location = useLocation();
  const { user } = useAuth();
  const [hasUnread, setHasUnread] = useState(false);

  useEffect(() => {
    // Request permission and token on mount
    requestForToken();

    // Listen for foreground messages
    const unsubscribe = onMessageListener((payload: any) => {
      console.log("Received foreground message:", payload);
      setHasUnread(true);
      toast(payload.notification?.title || "New Notification", {
        description: payload.notification?.body,
        action: {
          label: "View",
          onClick: () => console.log("Notification clicked"),
        },
      });
    });

    return () => {
      unsubscribe();
    };
  }, []);

  const title =
    TITLE_BY_PREFIX.find((t) => location.pathname.startsWith(t.prefix))
      ?.title ?? "Dashboard";

  return (
    <header className="sticky top-0 z-20 border-b border-border bg-elevated/90 text-elevated-foreground backdrop-blur">
      <div className="mx-auto flex h-16 w-full max-w-[1180px] items-center gap-4 px-4 md:gap-6 md:px-6 lg:px-8">
        {/* Mobile menu button */}
        <Button
          variant="ghost"
          size="icon"
          className="lg:hidden"
          onClick={onMenuClick}
          aria-label="Open menu"
        >
          <Menu className="h-5 w-5" />
        </Button>

        <div className="min-w-0 flex-1 md:min-w-[240px] md:flex-none">
          <h1 className="truncate text-lg font-semibold tracking-tight md:text-xl">
            {title}
          </h1>
          <p className="hidden text-xs text-muted-foreground md:block">
            Welcome back, {user?.first_name || "Admin"}
          </p>
        </div>

        <div className="hidden flex-1 md:block">
          <div className="relative mx-auto w-full max-w-[420px]">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              placeholder="Search..."
              className="h-10 rounded-xl bg-muted/30 pl-10"
            />
          </div>
        </div>

        <div className="flex items-center gap-2">
          <ThemeToggle />

          <Button
            variant="ghost"
            size="icon"
            onClick={() => setHasUnread(false)}
          >
            <Bell className="h-4 w-4" />
            {hasUnread && (
              <span className="absolute right-2 top-2 h-2 w-2 rounded-full bg-destructive" />
            )}
          </Button>

          {/* User Menu */}
          <UserMenu variant="topbar" />
        </div>
      </div>
    </header>
  );
}
