import { useState, useEffect } from "react";
import { useFormContext } from "react-hook-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Plus } from "lucide-react";
import { ProductFormValues } from "./formSchema";
import { categoryService } from "@/services/category.service";
import { brandService } from "@/services/brand.service";
import { useAuth } from "@/contexts/AuthContext";
import { CreateCategoryDialog } from "../dialogs/CreateCategoryDialog";
import { CreateBrandDialog } from "../dialogs/CreateBrandDialog";

export function DetailsStep() {
  const { control, watch, setValue } = useFormContext<ProductFormValues>();
  const { user } = useAuth();

  const [categories, setCategories] = useState<any[]>([]);
  const [brands, setBrands] = useState<any[]>([]);
  const [subcategories, setSubcategories] = useState<any[]>([]);

  const [isCategoryDialogOpen, setIsCategoryDialogOpen] = useState(false);
  const [isBrandDialogOpen, setIsBrandDialogOpen] = useState(false);

  const selectedCategoryId = watch("category_id");

  useEffect(() => {
    if (user?.organization_id) {
      fetchCategories();
      fetchBrands();
    }
  }, [user?.organization_id]);

  useEffect(() => {
    if (selectedCategoryId) {
      const category = categories.find((c) => c.id === selectedCategoryId);
      if (category && category.subcategories) {
        setSubcategories(category.subcategories);
      } else {
        setSubcategories([]);
      }
      // Reset subcategory if it doesn't belong to new category
      // setValue("subcategory_id", ""); // Optional: might annoy user if just browsing
    } else {
      setSubcategories([]);
    }
  }, [selectedCategoryId, categories]);

  const fetchCategories = async () => {
    if (!user?.organization_id) return;
    try {
      const response = await categoryService.getCategories({
        organization_id: user.organization_id,
        limit: 100,
      });
      setCategories(response.data || []);
    } catch (error) {
      console.error("Failed to fetch categories", error);
    }
  };

  const fetchBrands = async () => {
    if (!user?.organization_id) return;
    try {
      const response = await brandService.getBrands({
        organization_id: user.organization_id,
        limit: 100,
      });
      setBrands(response.brands || []);
    } catch (error) {
      console.error("Failed to fetch brands", error);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <Card>
        <CardHeader>
          <CardTitle>Categorization & Details</CardTitle>
        </CardHeader>
        <CardContent className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Category Field */}
            <div className="flex gap-2 items-end">
              <FormField
                control={control}
                name="category_id"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Category</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      value={field.value || ""}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select category" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {categories.map((category) => (
                          <SelectItem key={category.id} value={category.id}>
                            {category.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="button"
                variant="outline"
                size="icon"
                className=""
                onClick={() => setIsCategoryDialogOpen(true)}
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>

            {/* Subcategory Field */}
            <FormField
              control={control}
              name="subcategory_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Subcategory (Optional)</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    value={field.value || ""}
                    disabled={!selectedCategoryId || subcategories.length === 0}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue
                          placeholder={
                            !selectedCategoryId
                              ? "Select category first"
                              : subcategories.length === 0
                                ? "No subcategories"
                                : "Select subcategory"
                          }
                        />
                      </SelectTrigger>
                    </FormControl>
                    <SelectContent>
                      {subcategories.map((sub) => (
                        <SelectItem key={sub.id} value={sub.id}>
                          {sub.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <FormMessage />
                </FormItem>
              )}
            />

            {/* Brand Field */}
            <div className="flex gap-2 items-end justify-center">
              <FormField
                control={control}
                name="brand_id"
                render={({ field }) => (
                  <FormItem className="flex-1">
                    <FormLabel>Brand / Manufacturer</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      value={field.value || ""}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select brand" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {brands.map((brand) => (
                          <SelectItem key={brand.id} value={brand.id}>
                            {brand.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button
                type="button"
                variant="outline"
                size="icon"
                className=""
                onClick={() => setIsBrandDialogOpen(true)}
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Dialogs */}
      <CreateCategoryDialog
        open={isCategoryDialogOpen}
        onOpenChange={setIsCategoryDialogOpen}
        onSuccess={(cat) => {
          fetchCategories();
          setValue("category_id", cat.id);
        }}
      />
      <CreateBrandDialog
        open={isBrandDialogOpen}
        onOpenChange={setIsBrandDialogOpen}
        onSuccess={(brand) => {
          fetchBrands();
          setValue("brand_id", brand.id);
        }}
      />
    </div>
  );
}
