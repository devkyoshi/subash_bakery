import { useState, useEffect } from "react";
import { useForm, FormProvider } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";
import { DetailsStep } from "./DetailsStep";
import { PricingStep } from "./PricingStep";
import { useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Check, ChevronRight, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useAuth } from "@/contexts/AuthContext";
import { productService } from "@/services/product.service";
import { productSchema, ProductFormValues } from "./formSchema";
import { BasicInfoStep } from "./BasicInfoStep";
import { ImageUploadStep } from "./ImageUploadStep";

const steps = [
  { id: "basic", title: "Basic Info" },
  { id: "details", title: "Details" },
  { id: "images", title: "Images" },
  { id: "pricing", title: "Pricing & Inventory" },
];

export function ProductFormLayout() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const [currentStep, setCurrentStep] = useState("basic");
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(false);

  const methods = useForm<ProductFormValues>({
    resolver: zodResolver(productSchema),
    defaultValues: {
      type: "finished_goods",
      status: "active",
      track_inventory: true,
      location_prices: [],
    },
    mode: "onChange",
  });

  useEffect(() => {
    if (id && user?.organization_id) {
      fetchProduct(id);
    }
  }, [id, user?.organization_id]);

  const fetchProduct = async (productId: string) => {
    try {
      setIsFetching(true);
      const product = await productService.getProduct(productId);

      // Reset form with fetched data
      methods.reset({
        name: product.name,
        description: product.description || "",
        sku: product.sku,
        barcode: product.barcode || "",
        type: product.type,
        status: product.status,
        category_id: product.category_id || "",
        subcategory_id: product.subcategory_id || "",
        brand_id: product.brand_id || "",
        base_unit_id: product.base_unit_id || "",

        track_inventory: product.track_inventory,
        reorder_level: product.reorder_level || 0,

        // Dimensions & Weight
        weight: product.weight || 0,
        weight_unit: product.weight_unit || "kg",
        length: product.length || 0,
        width: product.width || 0,
        height: product.height || 0,
        dimension_unit: product.dimension_unit || "cm",
        volume: product.volume || 0,
        volume_unit: product.volume_unit || "m3",

        tags: product.tags || [],
        images: product.images || [],

        location_prices: product.location_prices.map((lp) => ({
          location_id: lp.location_id,
          location_name: lp.location_name,
          cost_price: lp.cost_price || 0,
          selling_price: lp.selling_price || 0,
          currency: lp.currency || "LKR",
          mrp: lp.mrp || 0,
          initial_stock: lp.initial_stock || 0,
          purchase_unit_id: lp.purchase_unit_id,
          selling_unit_id: lp.selling_unit_id,
        })),
      });
    } catch (error) {
      console.error("Failed to fetch product:", error);
      toast.error("Failed to load product details");
      navigate("/app/products");
    } finally {
      setIsFetching(false);
    }
  };

  const onSubmit = async (values: ProductFormValues) => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);

      const payload = {
        ...values,
        organization_id: user.organization_id,
      };

      if (id) {
        await productService.updateProduct(id, payload as any);
        toast.success("Product updated successfully");
      } else {
        await productService.createProduct(payload as any);
        toast.success("Product created successfully");
      }

      navigate("/app/products");
    } catch (error: any) {
      console.error("Failed to save product:", error);
      toast.error(
        id ? "Failed to update product" : "Failed to create product",
        {
          description:
            error.response?.data?.message || "Please check your inputs",
        },
      );
    } finally {
      setIsLoading(false);
    }
  };

  const onError = (errors: any) => {
    console.error("Form validation errors:", errors);
    toast.error("Please check the form for errors", {
      description: "Some required fields are missing or invalid.",
    });
  };

  const nextStep = async () => {
    const fields = methods.watch();
    let valid = false;

    if (currentStep === "basic") {
      valid = await methods.trigger(["name", "sku", "barcode"]);
      if (valid) setCurrentStep("details");
    } else if (currentStep === "details") {
      setCurrentStep("images");
    } else if (currentStep === "images") {
      setCurrentStep("pricing");
    }
  };

  const prevStep = () => {
    if (currentStep === "details") setCurrentStep("basic");
    else if (currentStep === "images") setCurrentStep("details");
    else if (currentStep === "pricing") setCurrentStep("images");
  };

  if (isFetching) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-5xl mx-auto pb-20">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/products")}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            {id ? "Edit Product" : "Create Product"}
          </h2>
          <p className="text-sm text-muted-foreground">
            {id
              ? "Modify product details"
              : "Add a new product to your inventory"}
          </p>
        </div>
      </div>

      {/* Progress Steps */}
      <div className="flex justify-between items-center max-w-2xl mx-auto mb-8 pb-4">
        {steps.map((step, index) => {
          const isActive = step.id === currentStep;
          const isCompleted =
            steps.findIndex((s) => s.id === currentStep) > index;

          return (
            <div
              key={step.id}
              className="flex flex-col items-center gap-2 relative"
            >
              <div
                className={cn(
                  "w-10 h-10 rounded-full flex items-center justify-center border-2 transition-colors",
                  isActive
                    ? "border-brand bg-brand text-brand-foreground"
                    : isCompleted
                      ? "border-brand bg-brand text-brand-foreground"
                      : "border-muted-foreground/30 bg-background text-muted-foreground",
                )}
              >
                {isCompleted ? (
                  <Check className="h-5 w-5" />
                ) : (
                  <span>{index + 1}</span>
                )}
              </div>
              <span
                className={cn(
                  "text-xs font-medium absolute -bottom-6 w-32 text-center",
                  isActive ? "text-foreground" : "text-muted-foreground",
                )}
              >
                {step.title}
              </span>
            </div>
          );
        })}
      </div>

      <FormProvider {...methods}>
        <form
          onSubmit={methods.handleSubmit(onSubmit, onError)}
          className="space-y-8 mt-8"
        >
          <Tabs
            value={currentStep}
            onValueChange={setCurrentStep}
            className="w-full"
          >
            <TabsContent value="basic" className="mt-0">
              <BasicInfoStep />
            </TabsContent>

            <TabsContent value="details" className="mt-0">
              <DetailsStep />
            </TabsContent>

            <TabsContent value="pricing" className="mt-0">
              <PricingStep />
            </TabsContent>

            <TabsContent value="images" className="mt-0">
              <ImageUploadStep />
            </TabsContent>
          </Tabs>

          {/* Navigation Buttons */}
          <div className="flex justify-between pt-6 border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 sticky bottom-0 z-10 p-4 rounded-lg border shadow-sm">
            <Button
              type="button"
              variant="outline"
              onClick={prevStep}
              disabled={currentStep === "basic"}
            >
              Back
            </Button>

            {currentStep === "pricing" ? (
              <Button
                key="submit-btn"
                type="submit"
                disabled={isLoading}
                className="bg-brand text-brand-foreground hover:bg-brand/90"
              >
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    {id ? "Updating..." : "Creating..."}
                  </>
                ) : (
                  <>{id ? "Update Product" : "Create Product"}</>
                )}
              </Button>
            ) : (
              <Button key="next-btn" type="button" onClick={nextStep}>
                Next Step
                <ChevronRight className="ml-2 h-4 w-4" />
              </Button>
            )}
          </div>
        </form>
      </FormProvider>
    </div>
  );
}
