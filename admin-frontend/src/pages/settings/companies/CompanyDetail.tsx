import { useState, useEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { ArrowLeft, Building2, Mail, Phone, Globe } from "lucide-react";
import { organizationService, Company } from "@/services/organization.service";
import { locationService } from "@/services/location.service";
import { userService, User } from "@/services/user.service";
import { Location } from "@/types/product.types";
import { toast } from "sonner";
import LocationList from "../locations/LocationList";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { useAuth } from "@/contexts/AuthContext";

export default function CompanyDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [company, setCompany] = useState<Company | null>(null);
  const [locations, setLocations] = useState<Location[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (id) {
      fetchCompany(id);
    }
  }, [id]);

  const { user } = useAuth(); // Ensure useAuth is imported/used

  const fetchCompany = async (companyId: string) => {
    setIsLoading(true);
    try {
      const orgId = user?.organization_id;

      const promises: [Promise<Company>, Promise<any>, Promise<any>] = [
        organizationService.getCompany(companyId),
        locationService.getLocations(companyId, { limit: 100 }),
        orgId
          ? userService.getUsers(orgId, { company_id: companyId, limit: 100 })
          : Promise.resolve({ data: { data: [] } }),
      ];

      const [companyData, locationsRes, usersRes] = await Promise.all(promises);

      setCompany(companyData);
      setLocations(locationsRes.data?.data || []); // Handle PaginatedData
      setUsers(usersRes.data?.data || []); // Handle GetUsersResponse structure
    } catch (error) {
      console.error("Failed to fetch company details", error);
      toast.error("Failed to load company details");
    } finally {
      setIsLoading(false);
    }
  };

  if (isLoading) return <div className="p-8">Loading company details...</div>;
  if (!company) return <div className="p-8">Company not found.</div>;

  return (
    <div className="space-y-6 animate-in fade-in duration-500 p-6">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/companies")}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div>
          <h2 className="text-3xl font-bold tracking-tight">{company.name}</h2>
          <p className="text-muted-foreground flex items-center gap-2">
            <Building2 className="h-4 w-4" />
            {company.code}
          </p>
        </div>
      </div>

      <Tabs defaultValue="info" className="w-full">
        <TabsList>
          <TabsTrigger value="info">Overview</TabsTrigger>
          <TabsTrigger value="locations">
            Locations ({locations.length})
          </TabsTrigger>
          <TabsTrigger value="users">Users ({users.length})</TabsTrigger>
        </TabsList>

        <TabsContent value="info" className="space-y-4 mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Company Information</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <h4 className="text-sm font-medium text-muted-foreground">
                    Address
                  </h4>
                  <p className="mt-1">
                    {company.address?.street}
                    <br />
                    {company.address?.city}, {company.address?.state}
                    <br />
                    {company.address?.country} {company.address?.postal_code}
                  </p>
                </div>
                <div className="space-y-2">
                  <h4 className="text-sm font-medium text-muted-foreground">
                    Contact Details
                  </h4>
                  <div className="flex justify-between border-b py-1">
                    <span>Email</span>
                    <span className="font-medium">{company.email}</span>
                  </div>
                  <div className="flex justify-between border-b py-1">
                    <span>Phone</span>
                    <span className="font-medium">{company.phone || "-"}</span>
                  </div>
                  <div className="flex justify-between border-b py-1">
                    <span>Website</span>
                    <span className="font-medium">
                      {company.website || "-"}
                    </span>
                  </div>

                  <h4 className="text-sm font-medium text-muted-foreground mt-4">
                    Identificators
                  </h4>
                  <div className="flex justify-between border-b py-1">
                    <span>Tax ID</span>
                    <span className="font-medium">{company.tax_id || "-"}</span>
                  </div>
                  <div className="flex justify-between border-b py-1">
                    <span>Currency</span>
                    <span className="font-medium">{company.currency}</span>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="locations" className="mt-4">
          <LocationList companyId={company.id} locations={locations} />
        </TabsContent>

        <TabsContent value="users" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Users</CardTitle>
              <CardDescription>Users assigned to this company.</CardDescription>
            </CardHeader>
            <CardContent>
              {users.length === 0 ? (
                <p className="text-muted-foreground text-sm">
                  No users assigned to this company.
                </p>
              ) : (
                <div className="space-y-4">
                  {users.map((user) => (
                    <div
                      key={user.id}
                      className="flex items-center justify-between border-b pb-4 last:border-0 last:pb-0"
                    >
                      <div className="flex items-center gap-4">
                        <Avatar>
                          <AvatarImage src={user.avatar} />
                          <AvatarFallback>
                            {user.first_name[0]}
                            {user.last_name[0]}
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="font-medium text-sm">
                            {user.first_name} {user.last_name}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {user.email}
                          </p>
                        </div>
                      </div>
                      <Badge variant={user.is_active ? "default" : "secondary"}>
                        {user.is_active ? "Active" : "Inactive"}
                      </Badge>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
