import { NavLink } from "@/components/NavLink";
import { SIDEBAR_NAV, SidebarNavItem } from "@/routes/sidebarNav";
import { cn } from "@/lib/utils";
import {
  ShieldCheck,
  X,
  ChevronDown,
  ChevronRight,
  PanelLeftClose,
  PanelLeftOpen,
  Store,
} from "lucide-react";
import { UserMenu } from "@/components/common/UserMenu";
import { useState } from "react";
import { useLocation } from "react-router-dom";
import { Button } from "@/components/ui/button";

interface DashboardSidebarProps {
  isOpen: boolean;
  onClose: () => void;
  isCollapsed: boolean;
  toggleCollapse: () => void;
}

export function DashboardSidebar({
  isOpen,
  onClose,
  isCollapsed,
  toggleCollapse,
}: DashboardSidebarProps) {
  const location = useLocation();
  const [expandedItems, setExpandedItems] = useState<string[]>(["Settings"]);

  const toggleExpand = (label: string) => {
    if (isCollapsed) {
      toggleCollapse(); // Expand sidebar if user clicks a parent item while collapsed
      setExpandedItems((prev) => [...prev, label]);
      return;
    }
    setExpandedItems((prev) =>
      prev.includes(label)
        ? prev.filter((item) => item !== label)
        : [...prev, label],
    );
  };

  const renderNavItem = (item: SidebarNavItem, depth = 0) => {
    const Icon = item.icon;
    const isExpanded = expandedItems.includes(item.label);
    const hasSubmenu = item.submenu && item.submenu.length > 0;
    const isActive =
      location.pathname === item.to ||
      (hasSubmenu &&
        item.submenu!.some((sub) => location.pathname.startsWith(sub.to)));

    // When collapsed, only show top-level items or items that perform a navigation
    // We hide submenus in collapsed mode for simplicity as per plan
    if (isCollapsed && depth > 0) return null;

    return (
      <li key={item.to}>
        {hasSubmenu ? (
          <div>
            <button
              onClick={() => toggleExpand(item.label)}
              className={cn(
                "flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-colors",
                isActive
                  ? "text-foreground"
                  : "text-muted-foreground hover:bg-muted/60 hover:text-foreground",
                isCollapsed && "justify-center px-2",
              )}
              style={
                !isCollapsed ? { paddingLeft: `${12 + depth * 12}px` } : {}
              }
              title={isCollapsed ? item.label : undefined}
            >
              <Icon className="h-4 w-4 shrink-0" />
              {!isCollapsed && (
                <>
                  <span className="flex-1 text-left">{item.label}</span>
                  {isExpanded ? (
                    <ChevronDown className="h-4 w-4" />
                  ) : (
                    <ChevronRight className="h-4 w-4" />
                  )}
                </>
              )}
            </button>
            {/* Recursively show subnav only if NOT collapsed and expanded */}
            {!isCollapsed && isExpanded && (
              <ul className="mt-1 space-y-1">
                {item.submenu!.map((subItem) =>
                  renderNavItem(subItem, depth + 1),
                )}
              </ul>
            )}
          </div>
        ) : (
          <NavLink
            to={item.to}
            onClick={onClose}
            className={cn(
              "flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-muted-foreground transition-colors",
              "hover:bg-muted/60 hover:text-foreground",
              isCollapsed && "justify-center px-2",
            )}
            activeClassName="bg-brand text-brand-foreground shadow-sm"
            style={!isCollapsed ? { paddingLeft: `${12 + depth * 12}px` } : {}}
            title={isCollapsed ? item.label : undefined}
          >
            <Icon className="h-4 w-4 shrink-0" />
            {!isCollapsed && <span className="flex-1">{item.label}</span>}
            {!isCollapsed && item.badge && (
              <span
                className={cn(
                  "rounded-full px-2 py-0.5 text-xs",
                  "bg-brand text-brand-foreground",
                )}
              >
                {item.badge}
              </span>
            )}
          </NavLink>
        )}
      </li>
    );
  };

  return (
    <>
      {/* Mobile overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/50 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <aside
        className={cn(
          "fixed left-0 top-0 z-50 h-screen border-r border-border bg-elevated text-elevated-foreground transition-all duration-300 lg:z-30 lg:translate-x-0",
          isOpen ? "translate-x-0" : "-translate-x-full",
          isCollapsed ? "w-[80px]" : "w-[280px]",
        )}
      >
        <div className="flex h-full flex-col">
          {/* Header with close button on mobile */}
          {/* Header */}
          <div
            className={cn(
              "flex items-center transition-all duration-300 py-5",
              isCollapsed ? "justify-center px-2" : "justify-between px-6",
            )}
          >
            <div
              className={cn(
                "flex items-center gap-3 overflow-hidden transition-all",
                isCollapsed && "cursor-pointer",
              )}
              onClick={isCollapsed ? toggleCollapse : undefined}
            >
              <div className="grid h-10 w-10 shrink-0 place-items-center rounded-xl bg-[#c25939] text-white">
                <Store className="h-6 w-6" />
              </div>
              {!isCollapsed && (
                <div className="flex flex-col">
                  <span className="font-bold tracking-tight text-foreground text-sm">
                    Subash Bakery
                  </span>
                  <span className="text-xs text-muted-foreground whitespace-nowrap">
                    Management System
                  </span>
                </div>
              )}
            </div>

            {/* Toggle Button - Visible in Expanded Mode for Desktop */}
            {!isCollapsed && (
              <Button
                variant="ghost"
                size="icon"
                onClick={toggleCollapse}
                className="hidden lg:flex h-8 w-8 text-muted-foreground hover:bg-muted hover:text-foreground shrink-0 ml-1"
              >
                <PanelLeftClose className="h-4 w-4" />
              </Button>
            )}

            {/* Close button - only visible on mobile */}
            <button
              onClick={onClose}
              className="grid h-9 w-9 place-items-center rounded-xl hover:bg-muted/60 lg:hidden"
              aria-label="Close menu"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          {/* Collapsed Toggle Hint (optional, or just handle via logo click) - 
              actually let's add a toggle button in the condensed view if needed, 
              but usually logo click is fine. 
              Let's add a small toggle button below logo in collapsed mode to be explicit? 
              No, minimal is better. But I'll make the logo click very obvious or add a small chevron.
          */}
          {isCollapsed && (
            <div className="hidden lg:flex justify-center mb-4">
              <Button
                variant="ghost"
                size="icon"
                onClick={toggleCollapse}
                className="h-6 w-6"
              >
                <PanelLeftOpen className="h-4 w-4" />
              </Button>
            </div>
          )}

          <nav className="flex-1 overflow-y-auto overflow-x-hidden px-3 mt-2">
            <ul className="space-y-1.5">
              {SIDEBAR_NAV.map((item) => renderNavItem(item))}
            </ul>
          </nav>

          {/* User Menu at bottom */}
          <div className="mt-auto border-t border-border p-4">
            <UserMenu variant="sidebar" isCollapsed={isCollapsed} />
          </div>
        </div>
      </aside>
    </>
  );
}
