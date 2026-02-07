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
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Building2,
  Plus,
  RefreshCw,
  ChevronRight,
  ChevronDown,
  MapPin,
  Warehouse,
  Pencil,
  Search,
  Filter,
  X,
} from "lucide-react";
import { Input } from "@/components/ui/input";
import {
  organizationService,
  CreateCompanyRequest,
  Company,
} from "@/services/organization.service";
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "sonner";
import { CompanyFormDialog } from "./CompanyFormDialog";
import { useNavigate } from "react-router-dom";

export default function CompanyList() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // Search state
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<string>("all");

  // Dialog states
  const [isCompanyDialogOpen, setIsCompanyDialogOpen] = useState(false);
  const [editingCompany, setEditingCompany] = useState<Company | null>(null);

  useEffect(() => {
    fetchData();
  }, [user?.organization_id, statusFilter]); // trigger on status change

  const fetchData = async () => {
    if (!user?.organization_id) return;

    setIsLoading(true);
    try {
      const response = await organizationService.getCompanies(
        user.organization_id,
        {
          limit: 100,
          q: searchQuery || undefined,
          is_active:
            statusFilter !== "all" ? statusFilter === "active" : undefined,
        },
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

  const handleSearch = () => {
    fetchData();
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
      </div>

      {/* Toolbar */}
      <div className="rounded-lg border border-border bg-elevated p-6 shadow-none">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div className="flex items-center gap-2">
              <div className="relative w-full sm:w-64">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search companies..."
                  className="h-10 pl-10"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  onKeyDown={(e) => e.key === "Enter" && handleSearch()}
                />
              </div>
              <Button variant="secondary" onClick={handleSearch}>
                Search
              </Button>
            </div>

            <div className="flex items-center gap-2">
              <Select
                value={statusFilter}
                onValueChange={(val) => {
                  setStatusFilter(val);
                }}
              >
                <SelectTrigger className="h-10 w-[140px]">
                  <Filter className="mr-2 h-4 w-4" />
                  <SelectValue placeholder="Status" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All Status</SelectItem>
                  <SelectItem value="active">Active</SelectItem>
                  <SelectItem value="inactive">Inactive</SelectItem>
                </SelectContent>
              </Select>

              {(searchQuery || statusFilter !== "all") && (
                <Button
                  variant="ghost"
                  onClick={() => {
                    setSearchQuery("");
                    setStatusFilter("all");
                    // fetch triggered by effect on status change or handled manually if we want immediate clear
                  }}
                >
                  <X className="mr-2 h-4 w-4" />
                  Clear
                </Button>
              )}
            </div>
          </div>

          <div className="flex items-center gap-2">
            <Button variant="outline" size="icon" onClick={() => fetchData()}>
              <RefreshCw className="h-4 w-4" />
            </Button>
            <Button onClick={openAddCompany}>
              <Plus className="mr-2 h-4 w-4" />
              Add Company
            </Button>
          </div>
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
