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
  Search,
  Plus,
  Pencil,
  Trash2,
  Tag,
  Globe,
  MapPin,
  Loader2,
  AlertTriangle,
} from "lucide-react";
import { brandService } from "@/services/brand.service";
import { Brand } from "@/types/brand.types";
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

export function BrandsPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [brands, setBrands] = useState<Brand[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const limit = 10;

  useEffect(() => {
    if (user?.organization_id) {
      fetchBrands();
    }
  }, [user?.organization_id, page]);

  const fetchBrands = async () => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);
      const response = await brandService.getBrands({
        organization_id: user.organization_id,
        q: searchQuery || undefined,
        page,
        limit,
      });

      setBrands(response.brands || []);
      setTotal(response.pagination?.total || 0);
    } catch (error: any) {
      console.error("Failed to fetch brands:", error);
      toast.error("Failed to fetch brands", {
        description: error.response?.data?.message || "Please try again later",
      });
      setBrands([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    setPage(1);
    fetchBrands();
  };

  const confirmDelete = async () => {
    if (!deleteId) return;

    try {
      await brandService.deleteBrand(deleteId);
      toast.success("Brand deleted successfully");
      fetchBrands();
    } catch (error: any) {
      toast.error("Failed to delete brand", {
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

  return (
    <div className="space-y-6">
      {/* Header Section */}
      <div>
        <h2 className="text-2xl font-semibold tracking-tight">
          Brand Management
        </h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage product brands, manufacturers, and their details.
        </p>
      </div>

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="relative w-full sm:w-64">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                placeholder="Search brands..."
                className="h-10 pl-10"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              />
            </div>
          </div>

          <div className="flex gap-2">
            <Button
              className="h-10 bg-brand text-brand-foreground hover:bg-brand/90 px-4"
              onClick={() => navigate("/app/brands/new")}
            >
              <Plus className="mr-2 h-4 w-4" />
              Add Brand
            </Button>
          </div>
        </div>
      </div>

      {/* Brands Table */}
      <div className="rounded-lg border border-border bg-elevated shadow-none overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-[30%]">Brand Name</TableHead>
              <TableHead>Code</TableHead>
              <TableHead>Location</TableHead>
              <TableHead>Website</TableHead>
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
                    Loading brands...
                  </div>
                </TableCell>
              </TableRow>
            ) : brands.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={6}
                  className="text-center py-8 text-muted-foreground"
                >
                  No brands found
                </TableCell>
              </TableRow>
            ) : (
              brands.map((brand) => (
                <TableRow key={brand.id} className="group">
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Tag className="h-4 w-4 text-brand" />
                      <span className="font-medium">{brand.name}</span>
                    </div>
                    {brand.description && (
                      <div className="ml-6 text-xs text-muted-foreground truncate max-w-[200px]">
                        {brand.description}
                      </div>
                    )}
                  </TableCell>
                  <TableCell className="font-mono text-sm text-muted-foreground">
                    {brand.code || "-"}
                  </TableCell>
                  <TableCell>
                    {brand.country ? (
                      <div className="flex items-center gap-1.5 text-sm">
                        <MapPin className="h-3 w-3 text-muted-foreground" />
                        {brand.country}
                      </div>
                    ) : (
                      "-"
                    )}
                  </TableCell>
                  <TableCell>
                    {brand.website ? (
                      <a
                        href={brand.website}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-1.5 text-sm text-blue-500 hover:underline"
                      >
                        <Globe className="h-3 w-3" />
                        Visit
                      </a>
                    ) : (
                      "-"
                    )}
                  </TableCell>
                  <TableCell>
                    <Badge
                      variant={brand.is_active ? "default" : "secondary"}
                      className={
                        brand.is_active
                          ? "bg-success text-success-foreground"
                          : ""
                      }
                    >
                      {brand.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                  <TableCell className="text-right">
                    <div className="flex justify-end gap-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8"
                        onClick={() => navigate(`/app/brands/${brand.id}/edit`)}
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="h-8 w-8 text-destructive hover:text-destructive"
                        onClick={() => setDeleteId(brand.id)}
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
        {!isLoading && brands.length > 0 && (
          <div className="flex items-center justify-between border-t px-6 py-4">
            <div className="text-sm text-muted-foreground">
              Showing {(page - 1) * limit + 1} to{" "}
              {Math.min(page * limit, total)} of {total} brands
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
              brand and remove it from our servers.
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
