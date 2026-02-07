import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { useFormContext, useFieldArray } from "react-hook-form";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { unitService } from "@/services/unit.service";
import { locationService } from "@/services/location.service";
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
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Plus, Trash2, MapPin } from "lucide-react";
import { ProductFormValues } from "./formSchema";
import { AddLocationPriceDialog } from "../dialogs/AddLocationPriceDialog";

// ... existing imports

export function PricingStep() {
  const { user } = useAuth();
  const { control, watch, setValue } = useFormContext<ProductFormValues>();
  const { fields, append, remove } = useFieldArray({
    control,
    name: "location_prices",
  });

  const [units, setUnits] = useState<any[]>([]);
  const [availableLocations, setAvailableLocations] = useState<any[]>([]);
  const [isLocationDialogOpen, setIsLocationDialogOpen] = useState(false);

  // Watch base unit to set default for new locations
  const baseUnitId = watch("base_unit_id");

  useEffect(() => {
    fetchUnits();
    fetchLocations();
  }, []);

  const fetchUnits = async () => {
    try {
      const data = await unitService.getUnits();
      setUnits(data || []);
      // Set default base unit if needed and not set
      // if (data.length > 0 && !watch("base_unit_id")) setValue("base_unit_id", data[0].id);
    } catch (error) {
      console.error("Failed to fetch units", error);
    }
  };

  const fetchLocations = async () => {
    try {
      if (user?.organization_id) {
        // Fetch all locations for the organization to ensure full list availability
        const data = await locationService.getOrganizationLocations(
          user.organization_id,
        );
        setAvailableLocations(data || []);
      }
    } catch (error) {
      console.error("Failed to fetch locations", error);
    }
  };

  const handleAddLocation = (data: any) => {
    // Check if already added
    if (fields.some((f) => f.location_id === data.location_id)) {
      return;
    }

    append({
      location_id: data.location_id,
      location_name: data.location_name,
      cost_price: data.cost_price,
      purchase_unit_id: data.purchase_unit_id,
      selling_price: data.selling_price,
      selling_unit_id: data.selling_unit_id,
      mrp: data.selling_price * 1.2, // Auto-calc MRP if not provided, or add field
      initial_stock: data.initial_stock,
      is_active: true,
      currency: "LKR",
    });
    setIsLocationDialogOpen(false);
  };

  // Filter out already selected locations
  const unselectedLocations = availableLocations.filter(
    (loc) => !fields.some((f) => f.location_id === loc.id),
  );

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Unit Settings */}
      <Card>
        <CardHeader>
          <CardTitle>Units & Inventory</CardTitle>
          <CardDescription>
            Configure how this product is measured.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              control={control}
              name="base_unit_id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Base Unit</FormLabel>
                  <Select
                    onValueChange={field.onChange}
                    value={field.value || ""}
                  >
                    <FormControl>
                      <SelectTrigger>
                        <SelectValue placeholder="Select base unit" />
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
                  <FormDescription>
                    The primary unit of measure for this product.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={control}
              name="track_inventory"
              render={({ field }) => (
                <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4 shadow-sm mt-8">
                  <div className="space-y-0.5">
                    <FormLabel>Track Inventory</FormLabel>
                    <FormDescription>
                      Monitor stock levels for this product.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Checkbox
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </div>
        </CardContent>
      </Card>

      {/* Location Pricing */}
      <Card>
        <CardHeader className="flex flex-row items-center justify-between">
          <div>
            <CardTitle>Location Pricing</CardTitle>
            <CardDescription>
              Set prices and stock for each location.
            </CardDescription>
          </div>
          <Button
            type="button"
            size="sm"
            variant="outline"
            onClick={() => setIsLocationDialogOpen(true)}
          >
            <Plus className="mr-2 h-4 w-4" />
            Add Location
          </Button>

          <AddLocationPriceDialog
            open={isLocationDialogOpen}
            onOpenChange={setIsLocationDialogOpen}
            onAdd={handleAddLocation}
            availableLocations={unselectedLocations}
            units={units}
          />
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[180px]">Location</TableHead>
                  <TableHead>Cost Price</TableHead>
                  <TableHead>Selling Price</TableHead>
                  <TableHead>MRP</TableHead>
                  <TableHead>Initial Stock</TableHead>
                  <TableHead className="w-[50px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {fields.length === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={6}
                      className="text-center py-8 text-muted-foreground"
                    >
                      No locations added. Click "Add Location" to start pricing.
                    </TableCell>
                  </TableRow>
                ) : (
                  fields.map((field, index) => {
                    const locationName =
                      availableLocations.find((l) => l.id === field.location_id)
                        ?.name || "Unknown";

                    return (
                      <TableRow key={field.id}>
                        <TableCell className="font-medium">
                          <div className="flex items-center gap-2">
                            <MapPin className="h-4 w-4 text-muted-foreground" />
                            {locationName}
                          </div>
                        </TableCell>
                        <TableCell>
                          <FormField
                            control={control}
                            name={`location_prices.${index}.cost_price`}
                            render={({ field }) => (
                              <FormItem>
                                <FormControl>
                                  <Input
                                    type="number"
                                    step="0.01"
                                    min="0"
                                    {...field}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </TableCell>
                        <TableCell>
                          <FormField
                            control={control}
                            name={`location_prices.${index}.selling_price`}
                            render={({ field }) => (
                              <FormItem>
                                <FormControl>
                                  <Input
                                    type="number"
                                    step="0.01"
                                    min="0"
                                    {...field}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </TableCell>
                        <TableCell>
                          <FormField
                            control={control}
                            name={`location_prices.${index}.mrp`}
                            render={({ field }) => (
                              <FormItem>
                                <FormControl>
                                  <Input
                                    type="number"
                                    step="0.01"
                                    min="0"
                                    {...field}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </TableCell>
                        <TableCell>
                          <FormField
                            control={control}
                            name={`location_prices.${index}.initial_stock`}
                            render={({ field }) => (
                              <FormItem>
                                <FormControl>
                                  <Input
                                    type="number"
                                    step="1"
                                    min="0"
                                    {...field}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => remove(index)}
                          >
                            <Trash2 className="h-4 w-4 text-destructive" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    );
                  })
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
