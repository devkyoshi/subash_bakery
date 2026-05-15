import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Search,
  Plus,
  Pencil,
  Trash2,
  Monitor,
  Loader2,
  Filter,
  X,
  Smartphone,
  Tablet,
  MonitorSmartphone,
  Laptop,
} from "lucide-react";
import { deviceService } from "@/services/device.service";
import { Device, DEVICE_TYPES } from "@/types/device.types";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "sonner";
import { useNavigate } from "react-router-dom";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";

const deviceTypeIcons: Record<string, React.ReactNode> = {
  pos: <Monitor className="h-4 w-4 text-brand" />,
  tablet: <Tablet className="h-4 w-4 text-brand" />,
  mobile: <Smartphone className="h-4 w-4 text-brand" />,
  desktop: <Laptop className="h-4 w-4 text-brand" />,
  kiosk: <MonitorSmartphone className="h-4 w-4 text-brand" />,
  other: <Monitor className="h-4 w-4 text-brand" />,
};

export function DevicesPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [devices, setDevices] = useState<Device[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [typeFilter, setTypeFilter] = useState<string>("all");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const limit = 10;

  useEffect(() => {
    if (user?.organization_id) {
      fetchDevices();
    }
  }, [user?.organization_id, page, statusFilter, typeFilter]);

  const fetchDevices = async () => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);
      const response = await deviceService.getDevices({
        organization_id: user.organization_id,
        search: searchQuery || undefined,
        is_active:
          statusFilter !== "all" ? statusFilter === "active" : undefined,
        device_type: typeFilter !== "all" ? (typeFilter as any) : undefined,
        page,
        limit,
      });

      setDevices(response.data || []);
      setTotal(response.pagination?.total || 0);
    } catch (error: any) {
      console.error("Failed to fetch devices:", error);
      toast.error("Failed to fetch devices", {
        description: error.response?.data?.message || "Please try again later",
      });
      setDevices([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    setPage(1);
    fetchDevices();
  };

  const confirmDelete = async () => {
    if (!deleteId) return;

    try {
      await deviceService.deleteDevice(deleteId);
      toast.success("Device deleted successfully");
      fetchDevices();
    } catch (error: any) {
      toast.error("Failed to delete device", {
        description:
          error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Please try again later",
      });
    } finally {
      setDeleteId(null);
    }
  };

  const totalPages = Math.ceil(total / limit);

  const getDeviceTypeLabel = (type: string) => {
    return DEVICE_TYPES.find((t) => t.value === type)?.label || type;
  };

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">
          Device Management
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Register and manage organization devices. Registered devices allow
          users to sign up without entering an organization ID.
        </p>
      </div>

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search devices..."
                  className="h-10 pl-10"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSearch()}
                />
              </div>
              <Button variant="secondary" onClick={handleSearch}>
                Search
              </Button>
            </div>

            <div className="flex items-center gap-2">
              <Select
                value={statusFilter}
                onValueChange={(val) => {
                  setStatusFilter(val);
                  setPage(1);
                }}
              >
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>

              <Select
                value={typeFilter}
                onValueChange={(val) => {
                  setTypeFilter(val);
                  setPage(1);
                }}
              >
                <SelectTrigger className="h-10 w-[160px]">
                  <Monitor className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Device Type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Types</SelectItem>
                  {DEVICE_TYPES.map((type) => (
                    <SelectItem key={type.value} value={type.value}>
                      {type.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>

              {(searchQuery ||
                statusFilter !== "all" ||
                typeFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setSearchQuery("");
                    setStatusFilter("all");
                    setTypeFilter("all");
                    setPage(1);
                  }}
                >
                  <X className="mr-2 h-4 w-4" />
                  Clear
                </Button>
              )}
            </div>
          </div>

          <div className="flex gap-2">
            <Button
              className="h-10 bg-brand text-brand-foreground hover:bg-brand/90 px-4"
              onClick={() => navigate("/app/devices/new")}
            >
              <Plus className="mr-2 h-4 w-4" />
              Register Device
            </Button>
          </div>
        </div>
      </div>

      {/* Devices Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[25%]">Device Name</TableHead>
              <TableHead>MAC Address</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center py-8 text-muted-foreground"
                >
                  <div className="flex justify-center items-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Loading devices...
                  </div>
                </TableCell>
              </TableRow>
            ) : devices.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center py-8 text-muted-foreground"
                >
                  No devices found
                </TableCell>
              </TableRow>
            ) : (
              devices.map((device) => (
                <TableRow key={device.id} className="group">
                  <TableCell>
                    <div className="flex items-center gap-2">
                      {deviceTypeIcons[device.device_type] || (
                        <Monitor className="h-4 w-4 text-brand" />
                      )}
                      <span className="font-medium">{device.name}</span>
                    </div>
                    {device.description && (
                      <div className="ml-6 text-xs text-muted-foreground truncate max-w-[200px]">
                        {device.description}
                      </div>
                    )}
                  </TableCell>
                  <TableCell className="font-mono text-sm text-muted-foreground">
                    {device.mac_address}
                  </TableCell>
                  <TableCell>
                    <Badge variant="outline">
                      {getDeviceTypeLabel(device.device_type)}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {device.location || "-"}
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={device.is_active ? "default" : "secondary"}
                      className={
                        device.is_active
                          ? "bg-success text-success-foreground"
                          : ""
                      }
                    >
                      {device.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() =>
                          navigate(`/app/devices/${device.id}/edit`)
                        }
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-destructive hover:text-destructive"
                        onClick={() => setDeleteId(device.id)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="flex items-center justify-between border-t px-4 py-3">
            <p className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} devices
            </p>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage((p) => p - 1)}
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
      </div>

      {/* Delete Confirmation Dialog */}
      <AlertDialog
        open={!!deleteId}
        onOpenChange={() => setDeleteId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Device</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this device? Users will no longer
              be able to register from this device without providing an
              organization ID. This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={confirmDelete}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
