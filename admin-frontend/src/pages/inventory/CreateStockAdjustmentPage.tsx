import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { inventoryService } from "@/services/inventory.service";
import { locationService } from "@/services/location.service";
import { productService } from "@/services/product.service";
import { Location } from "@/types/product.types";
import { Product, ProductStatus } from "@/types/product.types";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
} from "@/components/ui/command";
import { Calendar } from "@/components/ui/calendar";
import { format } from "date-fns";
import {
  ArrowLeft,
  Calendar as CalendarIcon,
  ChevronsUpDown,
  Check,
  Plus,
  Trash2,
  Loader2,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useToast } from "@/components/ui/use-toast";

interface AdjustmentItem {
  productId: string;
  productName: string;
  expectedQty: number;
  actualQty: number;
  uom: string;
  unitCost: number;
}

export function CreateStockAdjustmentPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();
  const [loading, setLoading] = useState(false);
  const [locations, setLocations] = useState<Location[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [productSearchOpen, setProductSearchOpen] = useState(false);

  // Form State
  const [selectedLocation, setSelectedLocation] = useState("");
  const [reason, setReason] = useState("");
  const [date, setDate] = useState<Date>(new Date());
  const [notes, setNotes] = useState("");
  const [items, setItems] = useState<AdjustmentItem[]>([]);

  useEffect(() => {
    fetchLocations();
    fetchProducts();
  }, [user?.organization_id]);

  const fetchLocations = async () => {
    if (!user?.organization_id) return;
    try {
      const data = await locationService.getOrganizationLocations(
        user.organization_id,
      );
      setLocations(data);
    } catch (error) {
      console.error("Failed to fetch locations", error);
    }
  };

  const fetchProducts = async () => {
    if (!user?.organization_id) return;
    try {
      // Fetch active products
      const response = await productService.getProducts({
        organization_id: user.organization_id,
        status: ProductStatus.ACTIVE,
        track_inventory: true,
        limit: 100, // TODO: Implement better search/pagination for products
      });
      setProducts(response.data);
    } catch (error) {
      console.error("Failed to fetch products", error);
    }
  };

  const addItem = async (product: Product) => {
    if (!selectedLocation) {
      toast({
        title: "Validation Error",
        description: "Please select a location first",
        variant: "destructive",
      });
      return;
    }

    if (items.some((item) => item.productId === product.id)) {
      toast({
        title: "Validation Error",
        description: "Product already added to adjustment",
        variant: "destructive",
      });
      return;
    }

    try {
      // Fetch current stock for this product at selected location
      const stockResponse = await inventoryService.getStockLevels({
        organization_id: user?.organization_id,
        product_id: product.id,
        location_id: selectedLocation,
      });

      // stockResponse.data.data is PaginatedData<StockLevel>, accessing .data gets StockLevel[]
      const levels = stockResponse.data.data.data;
      const currentStock =
        levels && levels.length > 0 ? levels[0].quantity_on_hand : 0;

      const newItem: AdjustmentItem = {
        productId: product.id,
        productName: product.name,
        expectedQty: currentStock,
        actualQty: currentStock, // Default to current, user changes it
        uom: "unit",
        unitCost: 0, // Ideally fetch cost price, assumed 0 or optional for now
      };

      setItems([...items, newItem]);
      setProductSearchOpen(false);
    } catch (error) {
      console.error("Failed to fetch stock level", error);
      toast({
        title: "Error",
        description: "Failed to fetch current stock for product",
        variant: "destructive",
      });
    }
  };

  const updateItemQty = (index: number, newQty: number) => {
    const newItems = [...items];
    newItems[index].actualQty = newQty;
    setItems(newItems);
  };

  const removeItem = (index: number) => {
    const newItems = [...items];
    newItems.splice(index, 1);
    setItems(newItems);
  };

  const handleSubmit = async () => {
    if (!user?.organization_id) return;
    if (!selectedLocation || !reason || items.length === 0) {
      toast({
        title: "Validation Error",
        description:
          "Please fill in all required fields and add at least one item",
        variant: "destructive",
      });
      return;
    }

    setLoading(true);
    try {
      const payload = {
        location_id: selectedLocation,
        adjustment_date: date.toISOString(),
        reason,
        notes,
        items: items.map((item) => ({
          product_id: item.productId,
          expected_qty: item.expectedQty,
          actual_qty: item.actualQty,
          uom: item.uom,
          unit_cost: item.unitCost,
        })),
      };

      await inventoryService.createStockAdjustment(
        user.organization_id,
        payload,
      );

      toast({
        title: "Success",
        description: "Stock adjustment created successfully",
      });
      navigate("/app/inventory/adjustments");
    } catch (error) {
      console.error("Failed to create adjustment", error);
      toast({
        title: "Error",
        description: "Failed to create adjustment",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/inventory/adjustments")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            New Stock Adjustment
          </h2>
          <p className="text-muted-foreground">
            Create a manual correction for stock levels
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <Card className="lg:col-span-2">
          <CardHeader>
            <CardTitle>Adjustment Items</CardTitle>
            <CardDescription>
              Add products and specify their actual counts.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-end">
              <Popover
                open={productSearchOpen}
                onOpenChange={setProductSearchOpen}
              >
                <PopoverTrigger asChild>
                  <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={productSearchOpen}
                    className="w-[250px] justify-between"
                    disabled={!selectedLocation}
                  >
                    Select Product...
                    <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-[250px] p-0">
                  <Command>
                    <CommandInput placeholder="Search product..." />
                    <CommandEmpty>No product found.</CommandEmpty>
                    <CommandGroup>
                      {products.map((product) => (
                        <CommandItem
                          key={product.id}
                          value={product.name}
                          onSelect={() => addItem(product)}
                        >
                          <Check
                            className={cn(
                              "mr-2 h-4 w-4",
                              items.some((i) => i.productId === product.id)
                                ? "opacity-100"
                                : "opacity-0",
                            )}
                          />
                          {product.name}
                        </CommandItem>
                      ))}
                    </CommandGroup>
                  </Command>
                </PopoverContent>
              </Popover>
            </div>

            <div className="rounded-md border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product</TableHead>
                    <TableHead className="w-[120px]">Expected</TableHead>
                    <TableHead className="w-[120px]">Actual</TableHead>
                    <TableHead className="w-[100px] text-right">Diff</TableHead>
                    <TableHead className="w-[50px]"></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {items.length === 0 ? (
                    <TableRow>
                      <TableCell colSpan={5} className="h-24 text-center">
                        No items added. Select products above.
                      </TableCell>
                    </TableRow>
                  ) : (
                    items.map((item, index) => (
                      <TableRow key={item.productId}>
                        <TableCell className="font-medium">
                          {item.productName}
                        </TableCell>
                        <TableCell>
                          {item.expectedQty} {item.uom}
                        </TableCell>
                        <TableCell>
                          <Input
                            type="number"
                            min="0"
                            value={item.actualQty}
                            onChange={(e) =>
                              updateItemQty(
                                index,
                                parseFloat(e.target.value) || 0,
                              )
                            }
                            className="h-8 w-24"
                          />
                        </TableCell>
                        <TableCell className="text-right">
                          <span
                            className={
                              item.actualQty - item.expectedQty > 0
                                ? "text-green-600"
                                : item.actualQty - item.expectedQty < 0
                                  ? "text-red-600"
                                  : ""
                            }
                          >
                            {item.actualQty - item.expectedQty > 0 ? "+" : ""}
                            {(item.actualQty - item.expectedQty).toFixed(2)}
                          </span>
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-destructive"
                            onClick={() => removeItem(index)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))
                  )}
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Details</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="location">Location</Label>
              <Select
                value={selectedLocation}
                onValueChange={(val) => {
                  if (items.length > 0 && val !== selectedLocation) {
                    if (
                      confirm(
                        "Changing location will clear added items. Continue?",
                      )
                    ) {
                      setItems([]);
                      setSelectedLocation(val);
                    }
                  } else {
                    setSelectedLocation(val);
                  }
                }}
              >
                <SelectTrigger id="location">
                  <SelectValue placeholder="Select location" />
                </SelectTrigger>
                <SelectContent>
                  {locations.map((loc) => (
                    <SelectItem key={loc.id} value={loc.id}>
                      {loc.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="date">Adjustment Date</Label>
              <Popover>
                <PopoverTrigger asChild>
                  <Button
                    variant={"outline"}
                    className={cn(
                      "w-full justify-start text-left font-normal",
                      !date && "text-muted-foreground",
                    )}
                  >
                    <CalendarIcon className="mr-2 h-4 w-4" />
                    {date ? format(date, "PPP") : <span>Pick a date</span>}
                  </Button>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0">
                  <Calendar
                    mode="single"
                    selected={date}
                    onSelect={(d) => d && setDate(d)}
                    initialFocus
                  />
                </PopoverContent>
              </Popover>
            </div>

            <div className="space-y-2">
              <Label htmlFor="reason">Reason</Label>
              <Select value={reason} onValueChange={setReason}>
                <SelectTrigger id="reason">
                  <SelectValue placeholder="Select reason" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Physical Count">Physical Count</SelectItem>
                  <SelectItem value="Damage">Damage / Expiry</SelectItem>
                  <SelectItem value="Theft">Theft / Loss</SelectItem>
                  <SelectItem value="Return">Customer Return</SelectItem>
                  <SelectItem value="Other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="notes">Notes</Label>
              <Textarea
                id="notes"
                placeholder="Optional notes..."
                value={notes}
                onChange={(e) => setNotes(e.target.value)}
                className="min-h-[100px]"
              />
            </div>
          </CardContent>
          <CardFooter>
            <Button
              className="w-full"
              onClick={handleSubmit}
              disabled={loading}
            >
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Create Adjustment
            </Button>
          </CardFooter>
        </Card>
      </div>
    </div>
  );
}
