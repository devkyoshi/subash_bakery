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
import {
  Building2,
  Plus,
  RefreshCw,
  ChevronRight,
  ChevronDown,
  MapPin,
  Warehouse,
  Pencil,
} from "lucide-react";
import {
  organizationService,
  CreateCompanyRequest,
  Company,
} from "@/services/organization.service";
import {
  locationService,
  CreateLocationRequest,
} from "@/services/location.service";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "sonner";
import { CompanyFormDialog } from "./CompanyFormDialog";
import { LocationFormDialog } from "../locations/LocationFormDialog";
import { useNavigate } from "react-router-dom";
import { Location } from "@/types/product.types";

export default function CompanyList() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Dialog states
  const [isCompanyDialogOpen, setIsCompanyDialogOpen] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);

  useEffect(() => {
    fetchData();
  }, [user?.organization_id]);

  const fetchData = async () => {
    if (!user?.organization_id) return;

    setIsLoading(true);
    try {
      const response = await organizationService.getCompanies(
        user.organization_id,
        { limit: 100 },
      );
      // Access nested data.data array safely
      setCompanies(response.data?.data || []);
    } catch (error) {
      console.error("Failed to fetch companies", error);
      toast.error("Failed to load companies");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateOrUpdateCompany = async (data: CreateCompanyRequest) => {
    if (!user?.organization_id) return;

    try {
      if (editingCompany) {
        await organizationService.updateCompany(editingCompany.id, data);
        toast.success("Company updated successfully");
      } else {
        await organizationService.createCompany(user.organization_id, data);
        toast.success("Company created successfully");
      }
      fetchData();
      setIsCompanyDialogOpen(false);
      setEditingCompany(null);
    } catch (error: any) {
      console.error("Failed to save company", error);
      toast.error(error.response?.data?.message || "Failed to save company");
    }
  };

  const openAddCompany = () => {
    setEditingCompany(null);
    setIsCompanyDialogOpen(true);
  };

  const openEditCompany = (company: Company, e: React.MouseEvent) => {
    e.stopPropagation();
    setEditingCompany(company);
    setIsCompanyDialogOpen(true);
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500 p-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-3xl font-bold tracking-tight">Companies</h2>
          <p className="text-muted-foreground">
            Manage your companies and settings.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="icon" onClick={fetchData}>
            <RefreshCw className="h-4 w-4" />
          </Button>
          <Button onClick={openAddCompany}>
            <Plus className="mr-2 h-4 w-4" />
            Add Company
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Organization Companies</CardTitle>
          <CardDescription>
            List of all companies within your organization.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="rounded-md border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[40%]">Name</TableHead>
                  <TableHead>Code</TableHead>
                  <TableHead>Headquarters</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {isLoading ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-24 text-center">
                      Loading...
                    </TableCell>
                  </TableRow>
                ) : companies.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={5} className="h-24 text-center">
                      No companies found. Add one to get started.
                    </TableCell>
                  </TableRow>
                ) : (
                  companies.map((company) => (
                    <TableRow
                      key={company.id}
                      className="cursor-pointer hover:bg-muted/50"
                      onClick={() => navigate(`/app/companies/${company.id}`)}
                    >
                      <TableCell className="font-medium">
                        <div className="flex items-center gap-2">
                          <Building2 className="h-4 w-4 text-primary" />
                          {company.name}
                        </div>
                      </TableCell>
                      <TableCell>{company.code || "-"}</TableCell>
                      <TableCell>{company.address?.city || "-"}</TableCell>
                      <TableCell>
                        <Badge
                          variant={company.is_active ? "default" : "secondary"}
                        >
                          {company.is_active ? "Active" : "Inactive"}
                        </Badge>
                      </TableCell>
                      <TableCell className="text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-8 w-8"
                          onClick={(e) => openEditCompany(company, e)}
                        >
                          <Pencil className="h-4 w-4 text-muted-foreground" />
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

      <CompanyFormDialog
        open={isCompanyDialogOpen}
        onOpenChange={setIsCompanyDialogOpen}
        onSubmit={handleCreateOrUpdateCompany}
        initialData={editingCompany}
      />
    </div>
  );
}
