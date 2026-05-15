import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Loader2, Plus, Pencil, Trash2 } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
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
import { useAuth } from "@/contexts/AuthContext";
import { roleService } from "@/services/role.service";
import axiosInstance from "@/lib/axios";

interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
  permissions: string[];
}

export function UserRoleList() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [deleteRole, setDeleteRole] = useState<Role | null>(null);

  const fetchRoles = async () => {
    if (!user?.organization_id) return;

    try {
      const fetchedRoles = await roleService.getRoles(user.organization_id);
      setRoles(fetchedRoles);
    } catch (error) {
      toast.error("Failed to load roles");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRoles();
  }, []);

  const handleDelete = async () => {
    if (!deleteRole) return;

    try {
      await axiosInstance.delete(`/roles/${deleteRole.id}`);
      toast.success("Role deleted successfully");
      fetchRoles();
    } catch (error) {
      toast.error("Failed to delete role");
    } finally {
      setDeleteRole(null);
    }
  };

  if (loading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            User Management
          </h2>
          <p className="text-sm text-muted-foreground">
            Manage all users and roles in one place. Control access, assign
            roles, and monitor activity across your platform.
          </p>
        </div>
        <Button onClick={() => navigate("/app/users/roles/new")}>
          <Plus className="mr-2 h-4 w-4" />
          Create Role
        </Button>
      </div>

      <div className="flex space-x-4 border-b">
        <Button
          variant={location.pathname === "/app/users" ? "secondary" : "ghost"}
          className="rounded-none border-b-2 border-transparent px-4 py-2 hover:bg-transparent hover:text-foreground data-[state=active]:border-primary"
          onClick={() => navigate("/app/users/all")}
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

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Description</TableHead>
              <TableHead>Type</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {roles.map((role) => (
              <TableRow key={role.id}>
                <TableCell>
                  <div className="font-medium">{role.display_name}</div>
                  <div className="text-xs text-muted-foreground">
                    {role.name}
                  </div>
                </TableCell>
                <TableCell>{role.description}</TableCell>
                <TableCell>
                  {role.is_system ? (
                    <Badge variant="secondary">System</Badge>
                  ) : (
                    <Badge
                      variant="outline"
                      className="border-green-500 text-green-500"
                    >
                      Custom
                    </Badge>
                  )}
                </TableCell>
                <TableCell className="text-right">
                  <div className="flex justify-end gap-2">
                    <Button
                      variant="ghost"
                      size="icon"
                      onClick={() =>
                        navigate(`/app/users/roles/${role.id}/edit`)
                      }
                      disabled={role.is_system}
                    >
                      <Pencil className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="text-destructive hover:text-destructive"
                      onClick={() => setDeleteRole(role)}
                      disabled={role.is_system}
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      <AlertDialog open={!!deleteRole} onOpenChange={() => setDeleteRole(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the
              role "{deleteRole?.display_name}" and remove it from our servers.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
              onClick={handleDelete}
            >
              Delete is
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
