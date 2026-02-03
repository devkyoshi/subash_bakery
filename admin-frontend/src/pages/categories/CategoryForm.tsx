import { useState, useEffect, useRef } from "react";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  ArrowLeft,
  Loader2,
  Plus,
  Trash2,
  Save,
  AlertTriangle,
} from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { categoryService } from "@/services/category.service";
import { useAuth } from "@/contexts/AuthContext";
import { Category } from "@/types/category.types";

const subcategorySchema = z.object({
  id: z.string().optional(),
  name: z.string().min(1, "Name is required"),
  code: z.string().optional(),
  is_active: z.boolean().default(true),
});

const formSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  code: z.string().optional(),
  description: z.string().optional(),
  is_active: z.boolean().default(true),
  subcategories: z.array(subcategorySchema).optional(),
});

type FormValues = z.infer<typeof formSchema>;

export function CategoryFormPage() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { user } = useAuth();
  const isEditMode = !!id;
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(isEditMode);

  // Store original category data to calculate diffs for updates
  const originalCategoryRef = useRef<Category | null>(null);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      code: "",
      description: "",
      is_active: true,
      subcategories: [],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: "subcategories",
  });

  useEffect(() => {
    if (isEditMode && id) {
      fetchCategory(id);
    }
  }, [isEditMode, id]);

  const fetchCategory = async (categoryId: string) => {
    try {
      setIsFetching(true);
      const category = await categoryService.getCategory(categoryId);
      originalCategoryRef.current = category;

      // Transform subcategories to match form shape
      const subcategories =
        category.subcategories?.map((sub) => ({
          id: sub.id,
          name: sub.name,
          code: sub.code || "",
          is_active: sub.is_active,
        })) || [];

      form.reset({
        name: category.name,
        code: category.code || "",
        description: category.description || "",
        is_active: category.is_active,
        subcategories: subcategories,
      });
    } catch (error) {
      console.error("Failed to fetch category:", error);
      toast.error("Failed to load category details");
      navigate("/app/categories");
    } finally {
      setIsFetching(false);
    }
  };

  const onSubmit = async (values: FormValues) => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);

      if (isEditMode && id) {
        // Calculate subcategory changes
        const currentSubcategories = values.subcategories || [];
        const originalSubcategories =
          originalCategoryRef.current?.subcategories || [];

        // 1. Identify Removed Subcategories
        // IDs present in original but NOT in current form values
        const currentIds = new Set(
          currentSubcategories.map((s) => s.id).filter(Boolean),
        );
        const removeSubcategories = originalSubcategories
          .filter((s) => !currentIds.has(s.id))
          .map((s) => s.id);

        // 2. Identify Added Subcategories
        // Items without IDs
        const addSubcategories = currentSubcategories
          .filter((s) => !s.id)
          .map((s) => ({
            name: s.name,
            code: s.code,
            is_active: s.is_active,
          }));

        // 3. Identify Updated Subcategories
        // Items with IDs that have changed
        const updateSubcategories = currentSubcategories
          .filter((s) => s.id)
          .map((s) => {
            const original = originalSubcategories.find(
              (orig) => orig.id === s.id,
            );
            if (!original) return null; // Should not happen

            // Check if modified
            const isModified =
              original.name !== s.name ||
              (original.code || "") !== (s.code || "") ||
              original.is_active !== s.is_active;

            if (isModified && s.id) {
              return {
                id: s.id,
                name: s.name,
                code: s.code,
                is_active: s.is_active,
              };
            }
            return null;
          })
          .filter((s): s is NonNullable<typeof s> => s !== null);

        const updateData: any = {
          name: values.name,
          code: values.code,
          description: values.description,
          is_active: values.is_active,
          add_subcategories:
            addSubcategories.length > 0 ? addSubcategories : undefined,
          update_subcategories:
            updateSubcategories.length > 0 ? updateSubcategories : undefined,
          remove_subcategories:
            removeSubcategories.length > 0 ? removeSubcategories : undefined,
        };

        if (Object.keys(updateData).length === 0) {
          toast.info("No changes detected");
          navigate("/app/categories");
          return;
        }

        // Don't send empty object if only name/etc didn't change but subcategories did
        // Actually updateData contains all fields like name which are always sent or we can send optional
        // In this implementation we sends name/code always. Backend handles it.

        await categoryService.updateCategory(id, updateData);
        toast.success("Category updated successfully");
      } else {
        // Create Logic remains same
        const createData: any = {
          organization_id: user.organization_id,
          name: values.name,
          code: values.code,
          description: values.description,
          is_active: values.is_active,
          subcategories: values.subcategories?.map((sub) => ({
            name: sub.name,
            code: sub.code,
            is_active: sub.is_active,
          })),
        };
        await categoryService.createCategory(createData);
        toast.success("Category created successfully");
      }

      navigate("/app/categories");
    } catch (error: any) {
      console.error("Failed to save category:", error);
      toast.error(
        isEditMode ? "Failed to update category" : "Failed to create category",
        {
          description:
            error.response?.data?.message || "Please try again later",
        },
      );
    } finally {
      setIsLoading(false);
    }
  };

  if (isFetching) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-10">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/categories")}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            {isEditMode ? "Edit Category" : "Add New Category"}
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditMode
              ? "Update category details and manage subcategories"
              : "Create a new category structure for your products"}
          </p>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          {/* Main Category Details */}
          <Card>
            <CardHeader>
              <CardTitle>Category Details</CardTitle>
              <CardDescription>
                Basic information about the category.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Category Name</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. Raw Materials" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="code"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Category Code (Optional)</FormLabel>
                      <FormControl>
                        <Input placeholder="e.g. RAW-001" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Describe what this category contains..."
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="is_active"
                render={({ field }) => (
                  <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                    <div className="space-y-0.5">
                      <FormLabel className="text-base">Active Status</FormLabel>
                      <FormDescription>
                        Inactive categories will be hidden from selection menus.
                      </FormDescription>
                    </div>
                    <FormControl>
                      <Switch
                        checked={field.value}
                        onCheckedChange={field.onChange}
                      />
                    </FormControl>
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>

          {/* Subcategories */}
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <div className="space-y-1">
                <CardTitle>Subcategories</CardTitle>
                <CardDescription>
                  Manage subcategories within this parent category.
                </CardDescription>
              </div>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={() => append({ name: "", code: "", is_active: true })}
              >
                <Plus className="mr-2 h-4 w-4" />
                Add Subcategory
              </Button>
            </CardHeader>
            <CardContent>
              {fields.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-8 text-center border-2 border-dashed rounded-lg bg-muted/20">
                  <div className="bg-muted rounded-full p-2 mb-2">
                    <AlertTriangle className="h-5 w-5 text-muted-foreground" />
                  </div>
                  <p className="text-sm font-medium">No subcategories yet</p>
                  <p className="text-xs text-muted-foreground mt-1 max-w-sm">
                    Add subcategories to better organize your items. For
                    example, "Flour" under "Raw Materials".
                  </p>
                  <Button
                    type="button"
                    variant="link"
                    className="mt-2 text-brand"
                    onClick={() =>
                      append({ name: "", code: "", is_active: true })
                    }
                  >
                    Add your first subcategory
                  </Button>
                </div>
              ) : (
                <div className="space-y-4">
                  {fields.map((field, index) => (
                    <div
                      key={field.id}
                      className="flex gap-4 items-start p-4 border rounded-lg bg-card hover:bg-muted/10 transition-colors"
                    >
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 flex-1">
                        <FormField
                          control={form.control}
                          name={`subcategories.${index}.name`}
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="text-xs">Name</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="Subcategory Name"
                                  {...field}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name={`subcategories.${index}.code`}
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="text-xs">Code</FormLabel>
                              <FormControl>
                                <Input
                                  placeholder="Code (Optional)"
                                  {...field}
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                        <FormField
                          control={form.control}
                          name={`subcategories.${index}.is_active`}
                          render={({ field }) => (
                            <FormItem className="flex flex-row items-center justify-between rounded-md border p-3 h-10 mt-6 md:mt-0 md:h-auto md:p-0 md:border-0">
                              <div className="md:hidden text-sm font-medium">
                                Active
                              </div>
                              <FormControl>
                                <div className="flex items-center gap-2 md:mt-8">
                                  <Switch
                                    checked={field.value}
                                    onCheckedChange={field.onChange}
                                  />
                                  <span className="text-xs text-muted-foreground hidden md:inline-block">
                                    {field.value ? "Active" : "Inactive"}
                                  </span>
                                </div>
                              </FormControl>
                            </FormItem>
                          )}
                        />
                      </div>
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="mt-8 text-muted-foreground hover:text-destructive"
                        onClick={() => remove(index)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>

          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate("/app/categories")}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={isLoading}
              className="bg-brand text-brand-foreground hover:bg-brand/90"
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <Save className="mr-2 h-4 w-4" />
              {isEditMode ? "Update Category" : "Create Category"}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
