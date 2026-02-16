import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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
  Search,
  Download,
  UserPlus,
  Filter,
  Pencil,
  Trash2,
  Calendar,
  Activity,
} from "lucide-react";
import { userService, User } from "@/services/user.service";
import { roleService } from "@/services/role.service";
import { Role } from "@/types/role.types";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "@/components/ui/sonner";
import { format } from "date-fns";

import { AddUserDialog } from "@/components/users/AddUserDialog";
import { useNavigate, useLocation } from "react-router-dom";

export function UsersPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]); // State for roles
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [roleFilter, setRoleFilter] = useState("all");
  const [statusFilter, setStatusFilter] = useState("all");
  const [dateFilter, setDateFilter] = useState("all");
  const [isAddUserOpen, setIsAddUserOpen] = useState(false);

  useEffect(() => {
    if (user?.organization_id) {
      fetchUsers();
    }
  }, [user?.organization_id, roleFilter, statusFilter, dateFilter]);

  useEffect(() => {
    if (user?.organization_id) {
      fetchUsers();
      fetchRoles(); // Fetch roles
    }
  }, [user?.organization_id, roleFilter, statusFilter, dateFilter]);

  const fetchRoles = async () => {
    if (!user?.organization_id) return;
    try {
      const fetchedRoles = await roleService.getRoles(user.organization_id);
      if (Array.isArray(fetchedRoles)) {
        setRoles(fetchedRoles);
      } else {
        setRoles([]);
      }
    } catch (error) {
      console.error("Failed to fetch roles:", error);
    }
  };

  const fetchUsers = async () => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);
      const response = await userService.getUsers(user.organization_id, {
        search: searchQuery || undefined,
        role: roleFilter !== "all" ? roleFilter : undefined,
        status: statusFilter !== "all" ? statusFilter : undefined,
      });
      setUsers(response.data?.data || []);
    } catch (error: any) {
      console.error("Failed to fetch users:", error);
      toast.error("Failed to fetch users", {
        description: error.response?.data?.message || "Please try again later",
      });
      setUsers([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    fetchUsers();
  };

  const handleDelete = async (userId: string) => {
    if (!confirm("Are you sure you want to delete this user?")) return;

    try {
      await userService.deleteUser(userId);
      toast.success("User deleted successfully");
      fetchUsers();
    } catch (error: any) {
      toast.error("Failed to delete user", {
        description: error.response?.data?.message || "Please try again later",
      });
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return "N/A";
    try {
      return format(new Date(dateString), "MMM dd, yyyy");
    } catch {
      return "Invalid date";
    }
  };

  const getFullName = (user: User) => {
    return `${user.first_name} ${user.last_name}`.trim() || "N/A";
  };

  const getInitials = (user: User) => {
    if (!user.first_name) return "U";
    return user.first_name.substring(0, 2).toUpperCase();
  };

  const getAvatarColor = (userId: string) => {
    const colors = [
      "bg-red-100 text-red-700",
      "bg-green-100 text-green-700",
      "bg-blue-100 text-blue-700",
      "bg-yellow-100 text-yellow-700",
      "bg-purple-100 text-purple-700",
      "bg-pink-100 text-pink-700",
      "bg-indigo-100 text-indigo-700",
      "bg-orange-100 text-orange-700",
    ];
    let hash = 0;
    for (let i = 0; i < userId.length; i++) {
      hash = userId.charCodeAt(i) + ((hash << 5) - hash);
    }
    return colors[Math.abs(hash) % colors.length];
  };

  const getRoleName = (roleId: string) => {
    if (!Array.isArray(roles)) return "Unknown Role";
    const role = roles.find((r) => r.id === roleId);
    return role ? role.display_name : "Unknown Role";
  };

  return (
    <div className="space-y-6">
      <AddUserDialog
        open={isAddUserOpen}
        onOpenChange={setIsAddUserOpen}
        onSuccess={fetchUsers}
      />

      {/* Header Section */}
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">
          User Management
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage all users and roles in one place. Control access, assign roles,
          and monitor activity across your platform.
        </p>
      </div>

      <div className="flex space-x-4 border-b">
        <Button
          variant={location.pathname === "/app/users" ? "secondary" : "ghost"}
          className="rounded-none border-b-2 border-transparent px-4 py-2 hover:bg-transparent hover:text-foreground data-[state=active]:border-primary"
          onClick={() => navigate("/app/users")}
        >
          Users
        </Button>
        <Button
          variant={
            location.pathname.includes("/app/users/roles")
              ? "secondary"
              : "ghost"
          }
          className="rounded-none border-b-2 border-transparent px-4 py-2 hover:bg-transparent hover:text-foreground data-[state=active]:border-primary"
          onClick={() => navigate("/app/users/roles")}
        >
          Roles
        </Button>
      </div>

      {/* Filter Section */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          {/* Left Side - Search and Filters */}
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            {/* Search Bar */}
            <div className="relative w-full sm:w-64">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search"
                className="h-10 pl-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
            </div>

            {/* Filter Dropdowns */}
            <div className="flex gap-2">
              {/* Role Filter */}
              <Select value={roleFilter} onValueChange={setRoleFilter}>
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Role" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Roles</SelectItem>
                  <SelectItem value="admin">Admin</SelectItem>
                  <SelectItem value="manager">Manager</SelectItem>
                  <SelectItem value="staff">Staff</SelectItem>
                  <SelectItem value="viewer">Viewer</SelectItem>
                </SelectContent>
              </Select>

              {/* Status Filter */}
              <Select value={statusFilter} onValueChange={setStatusFilter}>
                <SelectTrigger className="h-10 w-[140px]">
                  <Activity className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>

              {/* Date Filter */}
              <Select value={dateFilter} onValueChange={setDateFilter}>
                <SelectTrigger className="h-10 w-[140px]">
                  <Calendar className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Date" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Time</SelectItem>
                  <SelectItem value="today">Today</SelectItem>
                  <SelectItem value="week">This Week</SelectItem>
                  <SelectItem value="month">This Month</SelectItem>
                  <SelectItem value="year">This Year</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Right Side - Action Buttons */}
          <div className="flex gap-2">
            {/* Export Button - Secondary with outline
            <Button
              variant="outline"
              className="h-10 bg-background hover:bg-muted/50"
            >
              <Download className="mr-2 h-4 w-4" />
              Export
            </Button> */}

            {/* Add New User Button - Primary with brand color */}
            <Button
              className="h-10 bg-brand text-brand-foreground hover:bg-brand/90 px-4"
              onClick={() => setIsAddUserOpen(true)}
            >
              <UserPlus className="mr-2 h-4 w-4" />
              Add New User
            </Button>
          </div>
        </div>
      </div>

      {/* Users Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Full Name</TableHead>
              <TableHead>Email Address</TableHead>
              <TableHead>Role</TableHead>
              <TableHead>Status</TableHead>
              <TableHead>Join Date</TableHead>
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
                  Loading users...
                </TableCell>
              </TableRow>
            ) : users.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center py-8 text-muted-foreground"
                >
                  No users found
                </TableCell>
              </TableRow>
            ) : (
              users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>
                    <div className="flex items-center gap-3">
                      <div
                        className={`grid h-9 w-9 place-items-center rounded-full text-xs font-semibold ${getAvatarColor(user.id)}`}
                      >
                        {getInitials(user)}
                      </div>
                      <span className="font-medium">{getFullName(user)}</span>
                    </div>
                  </TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>
                    <Badge variant="outline">
                      {user.role?.display_name || getRoleName(user.role_id)}
                    </Badge>
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={user.is_active ? "default" : "secondary"}
                      className={
                        user.is_active
                          ? "bg-success text-success-foreground"
                          : "bg-muted text-muted-foreground"
                      }
                    >
                      {user.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell>{formatDate(user.created_at)}</TableCell>

                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => {
                          // TODO: Navigate to edit page
                          toast.info("Edit user functionality coming soon");
                        }}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-destructive hover:text-destructive"
                        onClick={() => handleDelete(user.id)}
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
      </div>
    </div>
  );
}
