import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  ArrowLeft,
  Pencil,
  Package,
  Barcode,
  Tag,
  Building2,
  Ruler,
  Scale,
  Archive,
  AlertTriangle,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { productService } from "@/services/product.service";
import { Product } from "@/types/product.types";
import { toast } from "sonner";
import { Skeleton } from "@/components/ui/skeleton";

export function ProductDetailsPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (id) {
      fetchProduct(id);
    }
  }, [id]);

  const fetchProduct = async (productId: string) => {
    try {
      setIsLoading(true);
      const data = await productService.getProduct(productId);
      setProduct(data);
    } catch (error) {
      console.error("Failed to fetch product:", error);
      toast.error("Failed to load product details");
    } finally {
      setIsLoading(false);
    }
  };

  const formatCurrency = (amount: number, currency: string = "USD") => {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: currency,
      minimumFractionDigits: 0,
    }).format(amount);
  };

  if (isLoading) {
    return <ProductDetailsSkeleton />;
  }

  if (!product) {
    return (
      <div className="flex flex-col items-center justify-center py-12">
        <h2 className="text-2xl font-semibold">Product not found</h2>
        <Button
          variant="outline"
          className="mt-4"
          onClick={() => navigate("/app/products")}
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          Back to Products
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/app/products")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h1 className="text-2xl font-semibold tracking-tight">
              {product.name}
            </h1>
            <div className="flex items-center gap-2 text-muted-foreground">
              <span className="font-mono text-sm">{product.sku}</span>
              {product.barcode && (
                <>
                  <span>•</span>
                  <div className="flex items-center gap-1">
                    <Barcode className="h-3 w-3" />
                    <span className="text-sm">{product.barcode}</span>
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
        <div className="flex gap-2">
          <Button onClick={() => navigate(`/app/products/${product.id}/edit`)}>
            <Pencil className="mr-2 h-4 w-4" />
            Edit Product
          </Button>
        </div>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {/* Main Info */}
        <div className="md:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Basic Information</CardTitle>
            </CardHeader>
            <CardContent className="grid gap-6 sm:grid-cols-2">
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground">
                  Type
                </span>
                <div className="flex items-center gap-2">
                  <Tag className="h-4 w-4 text-muted-foreground" />
                  <span className="capitalize">
                    {product.type.replace("_", " ")}
                  </span>
                </div>
              </div>
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground">
                  Status
                </span>
                <Badge
                  variant={
                    product.status === "active" ? "default" : "secondary"
                  }
                >
                  {product.status.toUpperCase()}
                </Badge>
              </div>

              {product.category && (
                <div className="space-y-1">
                  <span className="text-sm font-medium text-muted-foreground">
                    Category
                  </span>
                  <div className="flex items-center gap-2">
                    <Package className="h-4 w-4 text-muted-foreground" />
                    <span>{product.category.name}</span>
                  </div>
                </div>
              )}

              {product.brand && (
                <div className="space-y-1">
                  <span className="text-sm font-medium text-muted-foreground">
                    Brand
                  </span>
                  <div className="flex items-center gap-2">
                    <Building2 className="h-4 w-4 text-muted-foreground" />
                    <span>{product.brand.name}</span>
                  </div>
                </div>
              )}

              <div className="space-y-1 sm:col-span-2">
                <span className="text-sm font-medium text-muted-foreground">
                  Description
                </span>
                <p className="text-sm">
                  {product.description || "No description provided."}
                </p>
              </div>
            </CardContent>
          </Card>

          {/* Location & Stock details */}
          <Card>
            <CardHeader>
              <CardTitle>Location Pricing & Stock</CardTitle>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Location</TableHead>
                    <TableHead>Stock Status</TableHead>
                    <TableHead className="text-right">Available</TableHead>
                    <TableHead className="text-right">Cost Price</TableHead>
                    <TableHead className="text-right">Selling Price</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {product.location_prices.map((lp) => (
                    <TableRow key={lp.location_id}>
                      <TableCell className="font-medium">
                        {lp.location_name}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant={
                            lp.available_stock > 0 ? "default" : "destructive"
                          }
                          className={
                            lp.available_stock > 0
                              ? "bg-success hover:bg-success/90"
                              : ""
                          }
                        >
                          {lp.available_stock > 0 ? "In Stock" : "Out of Stock"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right font-mono">
                        {lp.available_stock}{" "}
                        {lp.selling_unit?.symbol || "units"}
                      </TableCell>
                      <TableCell className="text-right font-mono text-muted-foreground">
                        {formatCurrency(lp.cost_price, lp.currency)}
                      </TableCell>
                      <TableCell className="text-right font-mono font-medium">
                        {formatCurrency(lp.selling_price, lp.currency)}
                      </TableCell>
                    </TableRow>
                  ))}
                  {product.location_prices.length === 0 && (
                    <TableRow>
                      <TableCell
                        colSpan={5}
                        className="text-center text-muted-foreground py-6"
                      >
                        No location prices configured
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          {/* Stock Summary */}
          <Card>
            <CardHeader>
              <CardTitle>Total Inventory</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">
                  Total Stock
                </span>
                <span className="text-lg font-bold">{product.total_stock}</span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Available</span>
                <span className="text-lg font-bold text-success">
                  {product.available_stock}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <span className="text-sm text-muted-foreground">Allocated</span>
                <span className="text-lg font-bold text-orange-500">
                  {product.allocated_stock}
                </span>
              </div>

              <div className="pt-4 border-t">
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <AlertTriangle className="h-4 w-4 text-yellow-500" />
                    <span className="text-sm font-medium">Reorder Level</span>
                  </div>
                  <span>{product.reorder_level}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Dimensions */}
          <Card>
            <CardHeader>
              <CardTitle>Dimensions & Weight</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center gap-3">
                <Scale className="h-4 w-4 text-muted-foreground" />
                <div className="flex-1">
                  <span className="text-sm text-muted-foreground block">
                    Weight
                  </span>
                  <span className="font-medium">
                    {product.weight > 0
                      ? `${product.weight} ${product.weight_unit}`
                      : "-"}
                  </span>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <Ruler className="h-4 w-4 text-muted-foreground" />
                <div className="flex-1">
                  <span className="text-sm text-muted-foreground block">
                    Dimensions (LxWxH)
                  </span>
                  <span className="font-medium">
                    {product.length > 0
                      ? `${product.length} x ${product.width} x ${product.height} ${product.dimension_unit}`
                      : "-"}
                  </span>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <Archive className="h-4 w-4 text-muted-foreground" />
                <div className="flex-1">
                  <span className="text-sm text-muted-foreground block">
                    Volume
                  </span>
                  <span className="font-medium">
                    {product.volume > 0
                      ? `${product.volume} ${product.volume_unit}`
                      : "-"}
                  </span>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}

function ProductDetailsSkeleton() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="space-y-2">
          <Skeleton className="h-8 w-64" />
          <Skeleton className="h-4 w-32" />
        </div>
        <Skeleton className="h-10 w-32" />
      </div>
      <div className="grid gap-6 md:grid-cols-3">
        <div className="md:col-span-2 space-y-6">
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[300px]" />
        </div>
        <div className="space-y-6">
          <Skeleton className="h-[200px]" />
          <Skeleton className="h-[200px]" />
        </div>
      </div>
    </div>
  );
}
