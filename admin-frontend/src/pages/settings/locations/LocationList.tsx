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
import { MapPin, Plus, RefreshCw, Warehouse } from "lucide-react";
import {
  locationService,
  CreateLocationRequest,
} from "@/services/location.service";
import { useAuth } from "@/contexts/AuthContext";
import { Location } from "@/types/product.types";
import { toast } from "sonner";
import { LocationFormDialog } from "./LocationFormDialog";

interface LocationListProps {
  companyId?: string;
  locations?: Location[];
}

export default function LocationList({
  companyId: propsCompanyId,
  locations: preFetchedLocations,
}: LocationListProps) {
  const { user } = useAuth();
  const [locations, setLocations] = useState<Location[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
  const [companyId, setCompanyId] = useState<string | null>(
    propsCompanyId || null,
  );

  useEffect(() => {
    // If propsCompanyId changes, update internal state
    if (propsCompanyId) {
      setCompanyId(propsCompanyId);
    }
  }, [propsCompanyId]);

  useEffect(() => {
    initData();
  }, [user, companyId, preFetchedLocations]); // Depends on companyId now

  const initData = async () => {
    if (!user) return;

    // If pre-fetched locations are provided, use them
    if (preFetchedLocations) {
      setLocations(preFetchedLocations);
      setIsLoading(false);
      return;
    }

    // If companyId is provided via props or state, fetch locations for IT
    if (companyId) {
      setIsLoading(true);
      try {
        const response = await locationService.getLocations(companyId, {
          limit: 100,
        });
        setLocations(response.data?.data || []);
      } catch (error) {
        console.error("Failed to fetch locations", error);
        toast.error("Failed to load locations");
      } finally {
        setIsLoading(false);
      }
      return;
    }

    // Otherwise (if no prop), utilize the 'my access' logic -> auto select first company
    if (!propsCompanyId) {
      // Only do auto-select logic if NOT in detail view
      setIsLoading(true);
      try {
        const accessData = await locationService.getUserAccess();
        if (accessData.companies && accessData.companies.length > 0) {
          setCompanyId(accessData.companies[0].id);
          // Logic will re-run due to companyId dependency change, so we don't need to fetch here
          // But actually, the useEffect depends on companyId, so verify loop
        } else {
          setLocations(accessData.locations || []); // fallback
          setIsLoading(false);
        }
      } catch (error) {
        console.error("Failed to fetch locations", error);
        setIsLoading(false);
      }
    }
  };

  const handleCreateLocation = async (data: CreateLocationRequest) => {
    if (!companyId) {
      toast.error("No company selected");
      return;
    }

    try {
      await locationService.createLocation(companyId, data);
      toast.success("Location created successfully");
      initData(); // Refresh list
      setIsCreateDialogOpen(false);
    } catch (error: any) {
      console.error("Failed to create location", error);
      toast.error(error.response?.data?.message || "Failed to create location");
    }
  };

  const fetchLocations = initData; // Alias for refresh button

  return (
    <div className="space-y-6 animate-in fade-in duration-500 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Locations</h2>
          <p className="text-muted-foreground">
            Manage your warehouses, stores, and other locations.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={fetchLocations}>
            <RefreshCw className="h-4 w-4" />
          </Button>
          <Button onClick={() => setIsCreateDialogOpen(true)}>
            <Plus className="mr-2 h-4 w-4" />
            Add Location
          </Button>
        </div>
      </div>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Code</TableHead>
              <TableHead>Type</TableHead>
              <TableHead>City</TableHead>
              <TableHead>Status</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell colSpan={5} className="h-24 text-center">
                  Loading...
                </TableCell>
              </TableRow>
            ) : locations.length === 0 ? (
              <TableRow>
                <TableCell colSpan={5} className="h-24 text-center">
                  No locations found. Add one to get started.
                </TableCell>
              </TableRow>
            ) : (
              locations.map((loc) => (
                <TableRow key={loc.id}>
                  <TableCell className="font-medium">
                    <div className="flex items-center gap-2">
                      <Warehouse className="h-4 w-4 text-muted-foreground" />
                      {loc.name}
                    </div>
                  </TableCell>
                  <TableCell>{loc.code}</TableCell>
                  <TableCell className="capitalize">{loc.type}</TableCell>
                  <TableCell>{loc.address?.city || "-"}</TableCell>
                  <TableCell>
                    <Badge variant={loc.is_active ? "default" : "secondary"}>
                      {loc.is_active ? "Active" : "Inactive"}
                    </Badge>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      <LocationFormDialog
        open={isCreateDialogOpen}
        onOpenChange={setIsCreateDialogOpen}
        onSubmit={handleCreateLocation}
      />
    </div>
  );
}
