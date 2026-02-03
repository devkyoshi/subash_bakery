import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { CreateUnitRequest, UpdateUnitRequest } from "@/services/unit.service";
import { Unit } from "@/types/product.types";

const unitSchema = z.object({
  code: z
    .string()
    .min(1, "Code is required")
    .transform((val) => val.toUpperCase()), // Ensure uppercase
  name: z.string().min(1, "Name is required"),
  symbol: z.string().min(1, "Symbol is required"),
  unit_type: z.string().min(1, "Type is required"),
  is_base_unit: z.boolean().default(false),
  is_active: z.boolean().default(true),
});

type UnitValues = z.infer<typeof unitSchema>;

interface UnitFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: any) => Promise<void>;
  initialData?: Unit;
}

export function UnitFormDialog({
  open,
  onOpenChange,
  onSubmit,
  initialData,
}: UnitFormDialogProps) {
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<UnitValues>({
    resolver: zodResolver(unitSchema),
    defaultValues: {
      code: "",
      name: "",
      symbol: "",
      unit_type: "quantity",
      is_base_unit: false,
      is_active: true,
    },
  });

  useEffect(() => {
    if (open) {
      if (initialData) {
        form.reset({
          code: initialData.code || "", // Code is required, but might be missing on legacy data or frontend model if not refreshed.
          name: initialData.name,
          symbol: initialData.symbol,
          unit_type: initialData.unit_type,
          is_base_unit: initialData.is_base_unit,
          is_active: initialData.is_active,
        });
      } else {
        form.reset({
          code: "",
          name: "",
          symbol: "",
          unit_type: "quantity",
          is_base_unit: false,
          is_active: true,
        });
      }
    }
  }, [open, initialData, form]);

  const handleSubmit = async (values: UnitValues) => {
    try {
      setIsLoading(true);
      await onSubmit(values);
      // Form reset handled by parent closing/opening or useEffect
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>
            {initialData ? "Edit Unit" : "Create New Unit"}
          </DialogTitle>
        </DialogHeader>

        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(handleSubmit)}
            className="space-y-4"
          >
            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="code"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Code</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="PCS, KG"
                        {...field}
                        disabled={!!initialData} // Backend doesn't support updating Code? Handler has no Code update.
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Kilogram, Box, Piece" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="symbol"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Symbol</FormLabel>
                    <FormControl>
                      <Input placeholder="kg, box, pc" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="unit_type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Type</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select type" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        <SelectItem value="quantity">Quantity</SelectItem>
                        <SelectItem value="weight">Weight</SelectItem>
                        <SelectItem value="volume">Volume</SelectItem>
                        <SelectItem value="length">Length</SelectItem>
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="is_base_unit"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">Base Unit</FormLabel>
                    <FormDescription>
                      Is this a standard reference unit?
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

            <FormField
              control={form.control}
              name="is_active"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">Active Status</FormLabel>
                    <FormDescription>
                      Enable or disable this unit
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

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading
                  ? "Saving..."
                  : initialData
                    ? "Update Unit"
                    : "Create Unit"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
