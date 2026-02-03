import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { useAuth } from "@/contexts/AuthContext";
import { procurementService } from "@/services/procurement.service";
import { locationService } from "@/services/location.service";
import { PurchaseOrder } from "@/types/procurement.types";
import { Location } from "@/types/product.types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Form,
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
import {
  CalendarIcon,
  ArrowLeft,
  Save,
  AlertCircle,
  CheckCircle2,
  Circle,
  PackageCheck,
  ClipboardList,
  FileCheck,
  AlertTriangle,
  ThumbsUp,
  XCircle,
} from "lucide-react";
import { format } from "date-fns";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import { cn } from "@/lib/utils";
import { useToast } from "@/components/ui/use-toast";
import { Textarea } from "@/components/ui/textarea";

// Schema
const createGRNSchema = z.object({
  purchase_order_id: z.string().min(1, "Purchase Order ID is required"),
  location_id: z.string().min(1, "Location is required"),
  receipt_date: z.date({
    required_error: "Receipt date is required",
  }),
  notes: z.string().max(500, "Notes must not exceed 500 characters").optional(),
  items: z
    .array(
      z.object({
        product_id: z.string(),
        sku: z.string().optional(),
        description: z.string().optional(),
        ordered_quantity: z.number(), // for reference
        received_quantity: z.coerce.number().min(0, "Must be positive"),
        batch_number: z.string().optional(),
        expiry_date: z.string().optional(), // simplified as string for now, or Date
        condition: z.enum(["good", "partial", "damaged"]).default("good"),
      }),
    )
    .min(1, "At least one item is required"),
});

type CreateGRNValues = z.infer<typeof createGRNSchema>;

// Stepper Component
const Stepper = ({ currentStep }: { currentStep: number }) => {
  const steps = [
    { title: "Select PO", step: 1, icon: ClipboardList },
    { title: "Verify Items", step: 2, icon: PackageCheck },
    { title: "Confirm", step: 3, icon: FileCheck },
  ];

  return (
    <div className="flex items-center space-x-4 mb-6">
      {steps.map((s, idx) => {
        const Icon = s.icon;
        return (
          <div key={s.step} className="flex items-center">
            <div
              className={cn(
                "flex items-center",
                s.step < currentStep
                  ? "text-green-600 font-medium"
                  : s.step === currentStep
                    ? "text-primary font-medium"
                    : "text-muted-foreground",
              )}
            >
              <div className="flex items-center justify-center mr-2">
                {s.step < currentStep || (s.step === 1 && currentStep > 1) ? (
                  <CheckCircle2 className="h-6 w-6 text-green-600" />
                ) : s.step === currentStep ? (
                  <div className="h-8 w-8 rounded-full border-2 border-primary flex items-center justify-center bg-background z-10">
                    <Icon className="h-4 w-4 text-primary" />
                  </div>
                ) : (
                  <div className="h-8 w-8 rounded-full border-2 border-muted-foreground/30 flex items-center justify-center">
                    <Icon className="h-4 w-4" />
                  </div>
                )}
              </div>
              <span className={cn(s.step === currentStep && "font-bold")}>
                {s.title}
              </span>
            </div>
            {idx < steps.length - 1 && (
              <div className="h-[2px] w-12 bg-border mx-4 hidden sm:block" />
            )}
          </div>
        );
      })}
    </div>
  );
};

export function CreateGRNPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const poId = searchParams.get("po_id");
  const { user } = useAuth();
  const { toast } = useToast();

  const [po, setPo] = useState<PurchaseOrder | null>(null);
  const [locations, setLocations] = useState<Location[]>([]);
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [currentStep, setCurrentStep] = useState(2); // Start at Step 2

  const form = useForm<CreateGRNValues>({
    resolver: zodResolver(createGRNSchema),
    defaultValues: {
      receipt_date: new Date(),
      items: [],
      purchase_order_id: poId || "",
    },
  });

  const { fields, replace } = useFieldArray({
    control: form.control,
    name: "items",
  });

  // Watch items for summary
  const watchItems = form.watch("items");

  const summary = watchItems.reduce(
    (acc, item) => {
      const qty = Number(item.received_quantity) || 0;
      acc.totalQty += qty;
      acc.totalItems += 1;
      if (item.condition) {
        acc.byCondition[item.condition] =
          (acc.byCondition[item.condition] || 0) + qty;
      }
      return acc;
    },
    { totalQty: 0, totalItems: 0, byCondition: {} as Record<string, number> },
  );

  useEffect(() => {
    if (poId) {
      fetchPO(poId);
    }
    if (user?.organization_id) {
      fetchLocations(user.organization_id);
    }
  }, [poId, user?.organization_id]);

  const fetchLocations = async (orgId: string) => {
    try {
      const locations = await locationService.getOrganizationLocations(orgId);
      setLocations(locations || []);
    } catch (error) {
      console.error("Failed to fetch locations", error);
    }
  };

  const fetchPO = async (id: string) => {
    try {
      setLoading(true);
      const data = await procurementService.getPurchaseOrder(id);
      setPo(data);
      form.setValue("purchase_order_id", data.id);
      if (data.delivery_location_id) {
        form.setValue("location_id", data.delivery_location_id);
      }

      const grnItems = data.items.map((item) => ({
        product_id: item.product_id,
        sku: item.sku || "",
        description: item.description || "",
        ordered_quantity: item.quantity,
        received_quantity: item.quantity,
        batch_number: "",
        expiry_date: "",
        condition: "good" as "good" | "partial" | "damaged",
      }));

      replace(grnItems);
    } catch (error) {
      console.error("Failed to fetch PO", error);
      toast({
        title: "Error",
        description: "Failed to load Purchase Order details",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleNext = async () => {
    if (currentStep === 2) {
      const valid = await form.trigger("items");
      if (valid) {
        setCurrentStep(3);
      }
    }
  };

  const handleBack = () => {
    if (currentStep === 3) {
      setCurrentStep(2);
    }
  };

  const onSubmit = async (values: CreateGRNValues) => {
    if (!user?.organization_id) return;
    try {
      setSubmitting(true);

      const payload = {
        purchase_order_id: values.purchase_order_id,
        location_id: values.location_id,
        receipt_date: values.receipt_date.toISOString(),
        notes: values.notes,
        items: values.items.map((item) => ({
          product_id: item.product_id,
          ordered_quantity: Number(item.ordered_quantity),
          received_quantity: Number(item.received_quantity),
          batch_number: item.batch_number,
          ...(item.expiry_date ? { expiry_date: item.expiry_date } : {}),
          condition: item.condition,
        })),
      };

      await procurementService.createGRN(user.organization_id, payload);

      toast({
        title: "Success",
        description: "Goods Receipt Note created successfully",
      });
      navigate("/app/procurement/grn");
    } catch (error) {
      console.error("Failed to create GRN", error);
      toast({
        title: "Error",
        description: "Failed to create GRN",
        variant: "destructive",
      });
    } finally {
      setSubmitting(false);
    }
  };

  if (!poId) {
    return (
      <div className="p-8 text-center space-y-4">
        <AlertCircle className="h-10 w-10 text-muted-foreground mx-auto" />
        <h3 className="text-lg font-medium">No Purchase Order Selected</h3>
        <p className="text-muted-foreground">
          Please select a Purchase Order to receive goods against.
        </p>
        <Button onClick={() => navigate("/app/procurement/orders")}>
          Go to Orders
        </Button>
      </div>
    );
  }

  if (loading) {
    return <div className="p-8 text-center">Loading PO details...</div>;
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" onClick={() => navigate(-1)}>
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">Receive Goods</h2>
            <p className="text-muted-foreground">Create GRN</p>
          </div>
        </div>
      </div>

      <Stepper currentStep={currentStep} />

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          {/* Step 2: Verify Items */}
          {currentStep === 2 && (
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center justify-between">
                  <p>Verify Items</p>

                  <p className="text-sm text-muted-foreground font-light">
                    PO: {po?.po_number}
                  </p>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-6">
                  {fields.map((field, index) => (
                    <div
                      key={field.id}
                      className="flex flex-col gap-4 border-b pb-4 last:border-0 last:pb-0"
                    >
                      <div className="flex justify-between items-start">
                        <div>
                          <p className="font-semibold text-lg">
                            {field.description || "Product"}
                          </p>
                          <div className="flex gap-4 text-sm text-muted-foreground">
                            <span>SKU: {field.sku || "N/A"}</span>
                          </div>
                        </div>
                      </div>

                      <div className="flex flex-col md:flex-row justify-between gap-4 md:gap-0 items-center">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4 items-center w-full md:w-auto">
                          <FormField
                            name={`items.${index}.ordered_quantity`}
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Ordered Qty</FormLabel>
                                <FormControl>
                                  <Input disabled value={field.value} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                          <FormField
                            control={form.control}
                            name={`items.${index}.received_quantity`}
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Received Qty</FormLabel>
                                <FormControl>
                                  <Input type="number" min="0" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </div>

                        <FormField
                          control={form.control}
                          name={`items.${index}.condition`}
                          render={({ field }) => (
                            <FormItem className="col-span-2 md:ml-4 mt-4 md:mt-0">
                              <FormLabel>Condition</FormLabel>
                              <FormControl>
                                <div className="flex flex-wrap gap-2">
                                  <div
                                    className={cn(
                                      "cursor-pointer flex items-center gap-2 px-3 py-2 rounded-md border text-sm transition-colors",
                                      field.value === "good"
                                        ? "bg-green-100 border-green-300 text-green-800"
                                        : "hover:bg-muted bg-background",
                                    )}
                                    onClick={() => field.onChange("good")}
                                  >
                                    <ThumbsUp className="h-4 w-4" />
                                    Good
                                  </div>
                                  <div
                                    className={cn(
                                      "cursor-pointer flex items-center gap-2 px-3 py-2 rounded-md border text-sm transition-colors",
                                      field.value === "partial"
                                        ? "bg-yellow-100 border-yellow-300 text-yellow-800"
                                        : "hover:bg-muted bg-background",
                                    )}
                                    onClick={() => field.onChange("partial")}
                                  >
                                    <AlertTriangle className="h-4 w-4" />
                                    Partial
                                  </div>
                                  <div
                                    className={cn(
                                      "cursor-pointer flex items-center gap-2 px-3 py-2 rounded-md border text-sm transition-colors",
                                      field.value === "damaged"
                                        ? "bg-red-100 border-red-300 text-red-800"
                                        : "hover:bg-muted bg-background",
                                    )}
                                    onClick={() => field.onChange("damaged")}
                                  >
                                    <XCircle className="h-4 w-4" />
                                    Damaged
                                  </div>
                                </div>
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        {/* <FormField
                          control={form.control}
                          name={`items.${index}.batch_number`}
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Batch # (Opt)</FormLabel>
                              <FormControl>
                                <Input {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        /> */}
                      </div>
                    </div>
                  ))}
                </div>

                <div className="flex justify-end mt-6">
                  <Button type="button" onClick={handleNext}>
                    Next: Confirm
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Step 3: Confirm */}
          {currentStep === 3 && (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="md:col-span-2 space-y-6">
                <Card>
                  <CardHeader>
                    <CardTitle>Receipt Details</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <FormField
                      control={form.control}
                      name="location_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Receiving Location</FormLabel>
                          <Select
                            onValueChange={field.onChange}
                            defaultValue={field.value}
                            value={field.value}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Select location" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              {locations.map((loc) => (
                                <SelectItem key={loc.id} value={loc.id}>
                                  {loc.name}
                                </SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="receipt_date"
                      render={({ field }) => (
                        <FormItem className="flex flex-col">
                          <FormLabel>Receipt Date</FormLabel>
                          <Popover>
                            <PopoverTrigger asChild>
                              <FormControl>
                                <Button
                                  variant={"outline"}
                                  className={cn(
                                    "w-full pl-3 text-left font-normal",
                                    !field.value && "text-muted-foreground",
                                  )}
                                >
                                  {field.value ? (
                                    format(field.value, "PPP")
                                  ) : (
                                    <span>Pick a date</span>
                                  )}
                                  <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                </Button>
                              </FormControl>
                            </PopoverTrigger>
                            <PopoverContent
                              className="w-auto p-0"
                              align="start"
                            >
                              <Calendar
                                mode="single"
                                selected={field.value}
                                onSelect={field.onChange}
                                disabled={
                                  (date) => date > new Date() // Cannot be in future ideally
                                }
                                initialFocus
                              />
                            </PopoverContent>
                          </Popover>
                          <FormMessage />
                        </FormItem>
                      )}
                    />

                    <FormField
                      control={form.control}
                      name="notes"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Notes</FormLabel>
                          <FormControl>
                            <Textarea
                              placeholder="Receipt notes (max 500 chars)..."
                              className="resize-none"
                              rows={4}
                              maxLength={500}
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                          <div className="text-xs text-muted-foreground text-right mt-1">
                            {field.value?.length || 0}/500
                          </div>
                        </FormItem>
                      )}
                    />
                  </CardContent>
                </Card>
              </div>

              <div className="space-y-6">
                <Card className="bg-muted/50">
                  <CardHeader>
                    <CardTitle>Summary</CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-4">
                    <div className="flex justify-between items-center text-sm">
                      <span className="text-muted-foreground">Total Items</span>
                      <span className="font-medium">{summary.totalItems}</span>
                    </div>
                    <div className="flex justify-between items-center text-sm">
                      <span className="text-muted-foreground">
                        Total Received Qty
                      </span>
                      <span className="font-medium">{summary.totalQty}</span>
                    </div>

                    <div className="pt-2 border-t space-y-2">
                      <p className="text-xs font-semibold text-muted-foreground mb-2">
                        By Condition
                      </p>
                      <div className="flex justify-between items-center text-sm">
                        <span className="flex items-center gap-2">
                          <span className="h-2 w-2 rounded-full bg-green-500"></span>
                          Good
                        </span>
                        <span>{summary.byCondition["good"] || 0}</span>
                      </div>
                      <div className="flex justify-between items-center text-sm">
                        <span className="flex items-center gap-2">
                          <span className="h-2 w-2 rounded-full bg-yellow-500"></span>
                          Partial
                        </span>
                        <span>{summary.byCondition["partial"] || 0}</span>
                      </div>
                      <div className="flex justify-between items-center text-sm">
                        <span className="flex items-center gap-2">
                          <span className="h-2 w-2 rounded-full bg-red-500"></span>
                          Damaged
                        </span>
                        <span>{summary.byCondition["damaged"] || 0}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                <div className="flex flex-col gap-4">
                  <Button type="button" variant="outline" onClick={handleBack}>
                    Back to Items
                  </Button>
                  <Button
                    type="submit"
                    className="w-full"
                    disabled={submitting}
                  >
                    {submitting ? (
                      "Creating..."
                    ) : (
                      <>
                        <Save className="mr-2 h-4 w-4" /> Confirm & Create GRN
                      </>
                    )}
                  </Button>
                </div>
              </div>
            </div>
          )}
        </form>
      </Form>
    </div>
  );
}
