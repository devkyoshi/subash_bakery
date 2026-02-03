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
import {
  CreateUnitChartRequest,
  UpdateUnitChartRequest,
  UnitChart,
} from "@/services/unit.service";
import { Unit } from "@/types/product.types";

const unitChartSchema = z
  .object({
    from_unit_id: z.string().min(1, "From Unit is required"),
    to_unit_id: z.string().min(1, "To Unit is required"),
    conversion_rate: z.coerce.number().min(0.000001, "Must be positive"),
    is_active: z.boolean().default(true),
  })
  .refine((data) => data.from_unit_id !== data.to_unit_id, {
    message: "From Unit and To Unit cannot be the same",
    path: ["to_unit_id"],
  });

type UnitChartValues = z.infer<typeof unitChartSchema>;

interface UnitChartFormDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSubmit: (data: any) => Promise<void>;
  initialData?: UnitChart;
  units: Unit[];
}

export function UnitChartFormDialog({
  open,
  onOpenChange,
  onSubmit,
  initialData,
  units,
}: UnitChartFormDialogProps) {
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<UnitChartValues>({
    resolver: zodResolver(unitChartSchema),
    defaultValues: {
      from_unit_id: "",
      to_unit_id: "",
      conversion_rate: 1,
      is_active: true,
    },
  });

  useEffect(() => {
    if (open) {
      if (initialData) {
        form.reset({
          from_unit_id: initialData.from_unit_id,
          to_unit_id: initialData.to_unit_id,
          conversion_rate: initialData.conversion_rate, // Was conversion_factor
          is_active: initialData.is_active,
        });
      } else {
        form.reset({
          from_unit_id: "",
          to_unit_id: "",
          conversion_rate: 1,
          is_active: true,
        });
      }
    }
  }, [open, initialData, form]);

  const handleSubmit = async (values: UnitChartValues) => {
    try {
      setIsLoading(true);
      await onSubmit(values);
    } catch (error) {
      console.error(error);
    } finally {
      setIsLoading(false);
    }
  };

  // Get selected 'from' unit to filter 'to' units
  const selectedFromUnitId = form.watch("from_unit_id");
  const selectedFromUnit = units.find((u) => u.id === selectedFromUnitId);

  // Filter units for the "To Unit" dropdown
  const filteredToUnits = units.filter((unit) => {
    // Must be same type as selected From unit
    if (selectedFromUnit && unit.unit_type !== selectedFromUnit.unit_type) {
      return false;
    }
    // Cannot be the same unit as From unit
    if (unit.id === selectedFromUnitId) {
      return false;
    }
    return true;
  });

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>
            {initialData ? "Edit Conversion Rule" : "Create Conversion Rule"}
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
                name="from_unit_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>From Unit</FormLabel>
                    <Select
                      onValueChange={(val) => {
                        field.onChange(val);
                        // Reset To Unit if it becomes invalid (e.g. type mismatch or same unit)
                        // Actually, just clearing it is safer to force re-selection
                        if (form.getValues("to_unit_id")) {
                          const currentTo = units.find(
                            (u) => u.id === form.getValues("to_unit_id"),
                          );
                          const newFrom = units.find((u) => u.id === val);
                          if (
                            currentTo &&
                            newFrom &&
                            currentTo.unit_type !== newFrom.unit_type
                          ) {
                            form.setValue("to_unit_id", "");
                          } else if (val === form.getValues("to_unit_id")) {
                            form.setValue("to_unit_id", "");
                          }
                        }
                      }}
                      defaultValue={field.value}
                      value={field.value}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue placeholder="Select unit" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {units.map((unit) => (
                          <SelectItem key={unit.id} value={unit.id}>
                            {unit.name} ({unit.symbol})
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
                name="to_unit_id"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>To Unit</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                      value={field.value}
                      disabled={!selectedFromUnitId}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue
                            placeholder={
                              selectedFromUnitId
                                ? "Select unit"
                                : "Select From Unit first"
                            }
                          />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {filteredToUnits.map((unit) => (
                          <SelectItem key={unit.id} value={unit.id}>
                            {unit.name} ({unit.symbol})
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <FormField
              control={form.control}
              name="conversion_rate"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Conversion Factor</FormLabel>
                  <FormControl>
                    <Input type="number" step="0.0001" {...field} />
                  </FormControl>
                  <FormDescription>
                    How many "To Units" make one "From Unit"?
                    <br />
                    Example: 1 Box = 10 Pieces (Factor: 10)
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="is_active"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3">
                  <div className="space-y-0.5">
                    <FormLabel className="text-base">Active</FormLabel>
                    <FormDescription>
                      Enable this conversion rule
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
                    ? "Update Rule"
                    : "Create Rule"}
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
