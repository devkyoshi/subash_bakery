import { useState, useEffect } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Plus, RefreshCw, ArrowRightLeft, Pencil } from "lucide-react";
import {
  unitService,
  CreateUnitChartRequest,
  UpdateUnitChartRequest,
  UnitChart,
} from "@/services/unit.service";
import { Unit } from "@/types/product.types";
import { toast } from "sonner";
import { UnitChartFormDialog } from "./UnitChartFormDialog";

// Interface matching the actual API response
interface UnitConversionResponse {
  to_unit_id: string;
  to_unit_name: string;
  to_unit_code: string;
  conversion_rate: number;
}

interface UnitResponse {
  id: string; // from_unit_id
  chart_id: string;
  name: string;
  code: string;
  unit_type: string;
  is_base_unit: boolean;
  is_active: boolean; // Added is_active
  conversion: UnitConversionResponse;
}

export default function UnitChartList() {
  const [charts, setCharts] = useState<UnitResponse[]>([]);
  const [units, setUnits] = useState<Unit[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingChart, setEditingChart] = useState<UnitChart | null>(null);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    setIsLoading(true);
    try {
      const [chartsData, unitsData] = await Promise.all([
        unitService.getUnitCharts(),
        unitService.getUnits(),
      ]);
      // The API returns UnitResponse[], but we might need to cast or handle it
      setCharts((chartsData as unknown as UnitResponse[]) || []);
      setUnits(unitsData || []);
    } catch (error) {
      console.error("Failed to fetch data", error);
      toast.error("Failed to load unit charts");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateChart = async (data: CreateUnitChartRequest) => {
    try {
      await unitService.createUnitChart(data);
      toast.success("Conversion rule created successfully");
      fetchData(); // Refresh to get updated list
      setIsCreateDialogOpen(false);
    } catch (error: any) {
      console.error("Failed to create conversion rule", error);
      toast.error(
        error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Failed to create conversion rule",
      );
    }
  };

  const handleUpdateChart = async (data: UpdateUnitChartRequest) => {
    if (!editingChart) return;
    try {
      await unitService.updateUnitChart(editingChart.id, data);
      toast.success("Conversion rule updated successfully");
      fetchData();
      setEditingChart(null);
    } catch (error: any) {
      console.error("Failed to update conversion rule", error);
      toast.error(
        error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Failed to update conversion rule",
      );
    }
  };

  // Helper to reconstruct a UnitChart object for editing
  // Since the API returns flattened data, we need to map it back to what the form expects
  const openEditDialog = (chartItem: UnitResponse) => {
    const unitChart: UnitChart = {
      id: chartItem.chart_id,
      from_unit_id: chartItem.id,
      to_unit_id: chartItem.conversion.to_unit_id,
      conversion_rate: chartItem.conversion.conversion_rate, // Changed from factor
      organization_id: "", // Not needed for update
      is_active: chartItem.is_active, // Now correctly using API value
      created_at: 0,
      updated_at: 0,
    };
    setEditingChart(unitChart);
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            Unit Conversions
          </h2>
          <p className="text-muted-foreground">
            Manage conversion rules between different units.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={() => fetchData()}>
            <RefreshCw className="h-4 w-4" />
          </Button>
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Conversion
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Conversion Rules</CardTitle>
          <CardDescription>
            Defines how one unit converts to another.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>From Unit</TableHead>
                  <TableHead>To Unit</TableHead>
                  <TableHead>Conversion Factor</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-24 text-center">
                      Loading...
                    </TableCell>
                  </TableRow>
                ) : charts.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-24 text-center">
                      No conversion rules found. Add one to get started.
                    </TableCell>
                  </TableRow>
                ) : (
                  charts.map((item, index) => (
                    <TableRow key={index}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <ArrowRightLeft className="h-4 w-4 text-muted-foreground" />
                          {item.name} ({item.code})
                        </div>
                      </TableCell>
                      <TableCell>
                        {item.conversion.to_unit_name} (
                        {item.conversion.to_unit_code})
                      </TableCell>
                      <TableCell>
                        1 {item.name} = {item.conversion.conversion_rate}{" "}
                        {item.conversion.to_unit_name}
                      </TableCell>
                      <TableCell>
                        {item.is_active ? (
                          <Badge className="bg-green-600">Active</Badge>
                        ) : (
                          <Badge variant="destructive">Inactive</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => openEditDialog(item)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
                          {/* Delete button removed */}
                        </div>
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <UnitChartFormDialog
        open={isCreateDialogOpen || !!editingChart}
        onOpenChange={(open) => {
          if (!open) {
            setIsCreateDialogOpen(false);
            setEditingChart(null);
          }
        }}
        onSubmit={editingChart ? handleUpdateChart : handleCreateChart}
        initialData={editingChart || undefined}
        units={units}
      />
    </div>
  );
}
