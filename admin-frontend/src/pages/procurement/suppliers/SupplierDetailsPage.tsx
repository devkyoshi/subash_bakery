import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { procurementService } from "@/services/procurement.service";
import { Supplier } from "@/types/procurement.types";
import { formatCurrency } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  ArrowLeft,
  Building2,
  Phone,
  Mail,
  MapPin,
  CreditCard,
  FileText,
  Clock,
  Globe,
  Star,
} from "lucide-react";

export function SupplierDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [supplier, setSupplier] = useState<Supplier | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) {
      fetchSupplier(id);
    }
  }, [id]);

  const fetchSupplier = async (supplierId: string) => {
    try {
      setLoading(true);
      const data = await procurementService.getSupplier(supplierId);
      setSupplier(data);
    } catch (error) {
      console.error("Failed to fetch supplier details", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="p-8 text-center">Loading supplier details...</div>;
  }

  if (!supplier) {
    return (
      <div className="p-8 text-center">
        <h3 className="text-lg font-medium text-destructive">
          Supplier not found
        </h3>
        <Button
          variant="outline"
          onClick={() => navigate("/app/procurement/suppliers")}
          className="mt-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to Suppliers
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/app/procurement/suppliers")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">
              {supplier.company_name}
            </h2>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant="outline">{supplier.supplier_code}</Badge>
              <Badge
                variant={
                  supplier.status === "active"
                    ? "default"
                    : supplier.status === "inactive"
                      ? "secondary"
                      : "destructive"
                }
              >
                {supplier.status?.toUpperCase()}
              </Badge>
            </div>
          </div>
        </div>
        <div className="flex gap-2">
          <Button
            variant="outline"
            onClick={() => navigate(`/app/procurement/suppliers/${id}/edit`)}
          >
            Edit Supplier
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Main Info */}
        <Card className="md:col-span-2">
          <CardHeader>
            <CardTitle>Company Information</CardTitle>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <Building2 className="h-4 w-4" /> Contact Person
                </span>
                <p className="text-base">{supplier.contact_person}</p>
              </div>
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <Mail className="h-4 w-4" /> Email
                </span>
                <p className="text-base">{supplier.email || "-"}</p>
              </div>
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <Phone className="h-4 w-4" /> Phone
                </span>
                <p className="text-base">{supplier.phone || "-"}</p>
              </div>
              <div className="space-y-1">
                <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                  <Globe className="h-4 w-4" /> Website
                </span>
                <p className="text-base">
                  {supplier.website ? (
                    <a
                      href={supplier.website}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="text-primary hover:underline"
                    >
                      {supplier.website}
                    </a>
                  ) : (
                    "-"
                  )}
                </p>
              </div>
            </div>

            {supplier.address && (
              <>
                <div className="h-px bg-border my-4" />
                <div className="space-y-2">
                  <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                    <MapPin className="h-4 w-4" /> Address
                  </span>
                  <div className="text-base">
                    <p>{supplier.address.street}</p>
                    <p>
                      {supplier.address.city}, {supplier.address.state}{" "}
                      {supplier.address.postal_code}
                    </p>
                    <p>{supplier.address.country}</p>
                  </div>
                </div>
              </>
            )}

            {supplier.notes && (
              <>
                <div className="h-px bg-border my-4" />
                <div className="space-y-2">
                  <span className="text-sm font-medium text-muted-foreground flex items-center gap-2">
                    <FileText className="h-4 w-4" /> Notes
                  </span>
                  <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                    {supplier.notes}
                  </p>
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {/* Side Info */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Financial Details</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground flex items-center gap-2">
                  <CreditCard className="h-4 w-4" /> Currency
                </span>
                <span className="font-medium">
                  {supplier.currency || "LKR"}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground flex items-center gap-2">
                  <Clock className="h-4 w-4" /> Payment Terms
                </span>
                <span className="font-medium">
                  {supplier.payment_terms} Days
                </span>
              </div>
              <div className="space-y-1 pt-2">
                <span className="text-sm text-muted-foreground flex items-center justify-between">
                  <span>Credit Limit</span>
                  <span className="font-medium">
                    {formatCurrency(supplier.credit_limit || 0)}
                  </span>
                </span>
              </div>
              <div className="flex justify-between items-center pt-2">
                <span className="text-sm text-muted-foreground flex items-center gap-2">
                  <FileText className="h-4 w-4" /> Tax ID
                </span>
                <span className="font-medium">{supplier.tax_id || "-"}</span>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Performance</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground flex items-center gap-2">
                  <Star className="h-4 w-4" /> Rating
                </span>
                <span className="font-medium">
                  {supplier.rating ? `${supplier.rating} / 5` : "N/A"}
                </span>
              </div>
              {/* Add more metrics if available */}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
