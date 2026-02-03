import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { ArrowLeft, Loader2, Save } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import axiosInstance from "@/lib/axios";

interface Permission {
  id: string;
  name: string;
  display_name: string;
  description: string;
  category: string;
}

const formSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  description: z.string().optional(),
  permission_ids: z.array(z.string()),
});

type FormValues = z.infer<typeof formSchema>;

export function UserRoleFormPage() {
  const navigate = useNavigate();
  const { id } = useParams();
  const isEditing = !!id;
  const [loading, setLoading] = useState(false);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [groupedPermissions, setGroupedPermissions] = useState<
    Record<string, Permission[]>
  >({});

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
      permission_ids: [],
    },
  });

  useEffect(() => {
    fetchPermissions();
    if (isEditing) {
      fetchRole();
    }
  }, [id]);

  useEffect(() => {
    const grouped = permissions.reduce(
      (acc, permission) => {
        const category = permission.category || "Other";
        if (!acc[category]) {
          acc[category] = [];
        }
        acc[category].push(permission);
        return acc;
      },
      {} as Record<string, Permission[]>,
    );
    setGroupedPermissions(grouped);
  }, [permissions]);

  const fetchPermissions = async () => {
    try {
      const response = await axiosInstance.get("/permissions");
      setPermissions(response.data.data || []);
    } catch (error) {
      toast.error("Failed to load permissions");
    }
  };

  const fetchRole = async () => {
    setLoading(true);
    try {
      const response = await axiosInstance.get(`/roles/${id}`);
      const role = response.data.data;

      form.reset({
        name: role.name,
        description: role.description,
        permission_ids: role.permissions,
      });
    } catch (error) {
      toast.error("Failed to load role details");
      navigate("/app/users/roles");
    } finally {
      setLoading(false);
    }
  };

  const onSubmit = async (values: FormValues) => {
    setLoading(true);
    try {
      if (isEditing) {
        await axiosInstance.put(`/roles/${id}`, values);
      } else {
        await axiosInstance.post("/roles", values);
      }

      toast.success(`Role ${isEditing ? "updated" : "created"} successfully`);
      navigate("/app/users/roles");
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Failed to save role");
    } finally {
      setLoading(false);
    }
  };

  if (loading && isEditing) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6 container mx-auto py-6 max-w-4xl">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/users/roles")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            {isEditing ? "Edit Role" : "Create New Role"}
          </h2>
          <p className="text-sm text-muted-foreground">
            Configure role details and permissions
          </p>
        </div>
      </div>

      <Card>
        <CardContent className="p-6">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
              <div className="grid gap-4 grid-cols-1">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Role Name</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Inventory Manager"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Description</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Describe the role's responsibilities"
                          className="resize-none"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="space-y-4">
                <div className="flex items-center gap-2">
                  <Label className="text-base font-semibold">Permissions</Label>
                  <span className="text-xs text-muted-foreground">
                    Select the permissions for this role
                  </span>
                </div>
                <Separator />

                <div className="space-y-6">
                  {Object.entries(groupedPermissions).map(
                    ([category, perms]) => (
                      <div key={category} className="space-y-3">
                        <h4 className="font-medium text-sm text-muted-foreground uppercase tracking-wider">
                          {category}
                        </h4>
                        <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                          {perms.map((permission) => (
                            <FormField
                              key={permission.id}
                              control={form.control}
                              name="permission_ids"
                              render={({ field }) => (
                                <FormItem
                                  key={permission.id}
                                  className="flex flex-row items-start space-x-3 space-y-0 rounded-md border p-4 shadow-sm"
                                >
                                  <FormControl>
                                    <Checkbox
                                      checked={field.value?.includes(
                                        permission.id,
                                      )}
                                      onCheckedChange={(checked) => {
                                        return checked
                                          ? field.onChange([
                                              ...field.value,
                                              permission.id,
                                            ])
                                          : field.onChange(
                                              field.value?.filter(
                                                (value) =>
                                                  value !== permission.id,
                                              ),
                                            );
                                      }}
                                    />
                                  </FormControl>
                                  <div className="space-y-1 leading-none">
                                    <FormLabel className="font-medium">
                                      {permission.display_name}
                                    </FormLabel>
                                    <FormDescription className="text-xs">
                                      {permission.description}
                                    </FormDescription>
                                  </div>
                                </FormItem>
                              )}
                            />
                          ))}
                        </div>
                      </div>
                    ),
                  )}
                </div>
              </div>

              <div className="flex justify-end gap-4">
                <Button
                  variant="outline"
                  type="button"
                  onClick={() => navigate("/app/users/roles")}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={loading}>
                  {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  <Save className="mr-2 h-4 w-4" />
                  {isEditing ? "Save Changes" : "Create Role"}
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
