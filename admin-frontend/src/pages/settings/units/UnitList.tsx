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
import { Plus, RefreshCw, Ruler, Trash2, Pencil } from "lucide-react";
import {
  unitService,
  CreateUnitRequest,
  UpdateUnitRequest,
} from "@/services/unit.service";
import { Unit } from "@/types/product.types";
import { toast } from "sonner";

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
import { UnitFormDialog } from "./UnitFormDialog";

export default function UnitList() {
  const [units, setUnits] = useState<Unit[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [editingUnit, setEditingUnit] = useState<Unit | null>(null);
  // Removed handleDeleteUnit and AlertDialog as per user request (soft delete/inactive via edit)
  // const [deletingUnitId, setDeletingUnitId] = useState<string | null>(null);

  useEffect(() => {
    fetchUnits();
  }, []);

  const fetchUnits = async () => {
    setIsLoading(true);
    try {
      const data = await unitService.getUnits();
      setUnits(data || []);
    } catch (error) {
      console.error("Failed to fetch units", error);
      toast.error("Failed to load units");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateUnit = async (data: CreateUnitRequest) => {
    try {
      await unitService.createUnit(data);
      toast.success("Unit created successfully");
      fetchUnits();
      setIsCreateDialogOpen(false);
    } catch (error: any) {
      console.error("Failed to create unit", error);
      toast.error(
        error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Failed to create unit",
      );
    }
  };

  const handleUpdateUnit = async (data: UpdateUnitRequest) => {
    if (!editingUnit) return;
    try {
      await unitService.updateUnit(editingUnit.id, data);
      toast.success("Unit updated successfully");
      fetchUnits();
      setEditingUnit(null);
    } catch (error: any) {
      console.error("Failed to update unit", error);
      toast.error(
        error.response?.data?.error?.message ||
          error.response?.data?.message ||
          "Failed to update unit",
      );
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">
            Units of Measure
          </h2>
          <p className="text-muted-foreground">
            Manage units used for products and inventory.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={() => fetchUnits()}>
            <RefreshCw className="h-4 w-4" />
          </Button>
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Unit
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>All Units</CardTitle>
          <CardDescription>A list of all registered units.</CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Symbol</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Base Unit</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="w-[100px]"></TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      Loading...
                    </TableCell>
                  </TableRow>
                ) : units.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="h-24 text-center">
                      No units found. Add one to get started.
                    </TableCell>
                  </TableRow>
                ) : (
                  units.map((unit) => (
                    <TableRow key={unit.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <Ruler className="h-4 w-4 text-muted-foreground" />
                          {unit.name}
                        </div>
                      </TableCell>
                      <TableCell>{unit.symbol}</TableCell>
                      <TableCell className="capitalize">
                        {unit.unit_type}
                      </TableCell>
                      <TableCell>
                        {unit.is_base_unit ? (
                          <Badge variant="default">Yes</Badge>
                        ) : (
                          <Badge variant="secondary">No</Badge>
                        )}
                      </TableCell>
                      <TableCell>
                        {unit.is_active ? (
                          <Badge className="bg-green-600">Active</Badge>
                        ) : (
                          <Badge variant="destructive">Inactive</Badge>
                        )}
                      </TableCell>
                      {/* Conversion Factor removed as not part of Unit */}
                      <TableCell></TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => setEditingUnit(unit)}
                          >
                            <Pencil className="h-4 w-4" />
                          </Button>
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

      <UnitFormDialog
        open={isCreateDialogOpen || !!editingUnit}
        onOpenChange={(open) => {
          if (!open) {
            setIsCreateDialogOpen(false);
            setEditingUnit(null);
          }
        }}
        onSubmit={editingUnit ? handleUpdateUnit : handleCreateUnit}
        initialData={editingUnit || undefined}
      />
      {/* AlertDialog removed */}
    </div>
  );
}
