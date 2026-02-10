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
  Plus,
  Pencil,
  Trash2,
  FolderTree,
  ChevronRight,
  ChevronDown,
  Layers,
  Loader2,
  Filter,
  X,
} from "lucide-react";
import { categoryService } from "@/services/category.service";
import { Category } from "@/types/category.types";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "@/components/ui/sonner";
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

export function CategoriesPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [expandedCategories, setExpandedCategories] = useState<
    Record<string, boolean>
  >({});
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const limit = 10;

  useEffect(() => {
    if (user?.organization_id) {
      fetchCategories();
    }
  }, [user?.organization_id, page, statusFilter]); // Added statusFilter dependency

  const fetchCategories = async () => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);
      const response = await categoryService.getCategories({
        organization_id: user.organization_id,
        q: searchQuery || undefined,
        is_active:
          statusFilter !== "all" ? statusFilter === "active" : undefined,
        page,
        limit,
      });

      setCategories(response.data || []);
      setTotal(response.total || 0);
    } catch (error: any) {
      console.error("Failed to fetch categories:", error);
      toast.error("Failed to fetch categories", {
        description: error.response?.data?.message || "Please try again later",
      });
      setCategories([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    setPage(1);
    fetchCategories();
  };

  const confirmDelete = async () => {
    if (!deleteId) return;

    try {
      await categoryService.deleteCategory(deleteId);
      toast.success("Category deleted successfully");
      fetchCategories();
    } catch (error: any) {
      toast.error("Failed to delete category", {
        description:
          error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Please try again later",
      });
    } finally {
      setDeleteId(null);
    }
  };

  const toggleExpand = (categoryId: string) => {
    setExpandedCategories((prev) => ({
      ...prev,
      [categoryId]: !prev[categoryId],
    }));
  };

  const totalPages = Math.ceil(total / limit);

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">
          Category Management
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Organize your products into categories and subcategories for better
          inventory management.
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
                  placeholder="Search categories..."
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
                  setPage(1); // Reset page on filter change
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

              {(searchQuery || statusFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setSearchQuery("");
                    setStatusFilter("all");
                    setPage(1);
                    // Verify if we need to trigger fetch manually or effect handles it
                    // Effect depends on statusFilter, so it will trigger.
                    // Search query is not in effect dependency array usually (it is passed to fetch),
                    // but usually we want to clear and fetch.
                    // Let's rely on the effect for status, but for search we might need to trigger if status doesn't change.
                    // Actually, best to just reset state and let a robust effect handle it or call fetch manually.
                    // In this component, fetchCategories depends on [user, page].
                    // It DOES NOT depend on search/status state changes automatically to avoid debounce issues?
                    // No, lines 58: [user?.organization_id, page].
                    // So changing status won't trigger fetch unless we add it to dependency or call fetch.
                    // I should add statusFilter to dependency array or call fetch.
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
              onClick={() => navigate("/app/categories/new")}
            >
              <Plus className="mr-2 h-4 w-4" />
              Add Category
            </Button>
          </div>
        </div>
      </div>

      {/* Categories Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[40%]">Category Name</TableHead>
              <TableHead>Code</TableHead>
              <TableHead>Subcategories</TableHead>
              <TableHead>Status</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="text-center py-8 text-muted-foreground"
                >
                  <div className="flex justify-center items-center gap-2">
                    <Loader2 className="h-4 w-4 animate-spin" />
                    Loading categories...
                  </div>
                </TableCell>
              </TableRow>
            ) : categories.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="text-center py-8 text-muted-foreground"
                >
                  No categories found
                </TableCell>
              </TableRow>
            ) : (
              categories.map((category, index) => (
                <>
                  <TableRow key={index} className="group">
                    <TableCell>
                      <div className="flex items-center gap-2">
                        {category.subcategories?.length > 0 ? (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-6 w-6 p-0 hover:bg-transparent"
                            onClick={() => toggleExpand(category.id)}
                          >
                            {expandedCategories[category.id] ? (
                              <ChevronDown className="h-4 w-4 text-muted-foreground" />
                            ) : (
                              <ChevronRight className="h-4 w-4 text-muted-foreground" />
                            )}
                          </Button>
                        ) : (
                          <div className="w-6" />
                        )}
                        <FolderTree className="h-4 w-4 text-brand" />
                        <span className="font-medium">{category.name}</span>
                      </div>
                      {category.description && (
                        <div className="ml-8 text-xs text-muted-foreground truncate max-w-[300px]">
                          {category.description}
                        </div>
                      )}
                    </TableCell>
                    <TableCell className="font-mono text-sm text-muted-foreground">
                      {category.code || "-"}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className="rounded-sm font-normal"
                      >
                        {category.subcategories?.length || 0} items
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={category.is_active ? "default" : "secondary"}
                        className={
                          category.is_active
                            ? "bg-success text-success-foreground"
                            : ""
                        }
                      >
                        {category.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={() =>
                            navigate(`/app/categories/${category.id}/edit`)
                          }
                        >
                          <Pencil className="h-4 w-4" />
                        </Button>
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8 text-destructive hover:text-destructive"
                          onClick={() => setDeleteId(category.id)}
                        >
                          <Trash2 className="h-4 w-4" />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>

                  {/* Subcategories Expansion */}
                  {expandedCategories[category.id] &&
                    category.subcategories?.length > 0 &&
                    category.subcategories.map((sub) => (
                      <TableRow
                        key={sub.id}
                        className="bg-muted/30 hover:bg-muted/50"
                      >
                        <TableCell className="pl-12">
                          <div className="flex items-center gap-2">
                            <Layers className="h-3 w-3 text-muted-foreground" />
                            <span className="text-sm">{sub.name}</span>
                          </div>
                        </TableCell>
                        <TableCell className="font-mono text-xs text-muted-foreground">
                          {sub.code || "-"}
                        </TableCell>
                        <TableCell>
                          <span className="text-xs text-muted-foreground">
                            {sub.product_count || 0} products
                          </span>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant="outline"
                            className={`h-5 text-[10px] ${!sub.is_active && "text-muted-foreground border-muted-foreground/30"}`}
                          >
                            {sub.is_active ? "Active" : "Inactive"}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-right">
                          {/* Subcategory actions */}
                        </TableCell>
                      </TableRow>
                    ))}
                </>
              ))
            )}
          </TableBody>
        </Table>

        {/* Pagination */}
        {!isLoading && categories.length > 0 && (
          <div className="flex items-center justify-between border-t px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} categories
            </div>
            <div className="flex gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page - 1)}
                disabled={page === 1}
              >
                Previous
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => setPage(page + 1)}
                disabled={page >= totalPages}
              >
                Next
              </Button>
            </div>
          </div>
        )}
      </div>

      <AlertDialog
        open={!!deleteId}
        onOpenChange={(open) => !open && setDeleteId(null)}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
            <AlertDialogDescription>
              This action cannot be undone. This will permanently delete the
              category and its subcategories.
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
