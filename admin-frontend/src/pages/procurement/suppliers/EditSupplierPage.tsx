import { useState, useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { procurementService } from "@/services/procurement.service";
import { CreateSupplierRequest } from "@/types/procurement.types";
import { SupplierForm } from "./SupplierForm";
import { useToast } from "@/components/ui/use-toast";

export function EditSupplierPage() {
  const navigate = useNavigate();
  const { id } = useParams<{ id: string }>();
  const { toast } = useToast();
  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);
  const [initialValues, setInitialValues] = useState<
    Partial<CreateSupplierRequest>
  >({});

  useEffect(() => {
    if (id) {
      fetchSupplier(id);
    }
  }, [id]);

  const fetchSupplier = async (supplierId: string) => {
    try {
      const supplier = await procurementService.getSupplier(supplierId);
      setInitialValues({
        company_name: supplier.company_name,
        contact_person: supplier.contact_person,
        email: supplier.email,
        phone: supplier.phone,
        mobile: supplier.mobile,
        website: supplier.website,
        address: supplier.address,
        tax_id: supplier.tax_id,
        payment_terms: supplier.payment_terms,
        credit_limit: supplier.credit_limit,
        bank_name: supplier.bank_name,
        account_number: supplier.account_number,
        swift_code: supplier.swift_code,
        notes: supplier.notes,
      });
    } catch (error) {
      console.error("Failed to fetch supplier", error);
      toast({
        title: "Error",
        description: "Failed to load supplier details.",
        variant: "destructive",
      });
      navigate("/app/procurement/suppliers");
    } finally {
      setFetching(false);
    }
  };

  const handleSubmit = async (values: CreateSupplierRequest) => {
    if (!id) return;

    try {
      setLoading(true);
      await procurementService.updateSupplier(id, values);
      toast({
        title: "Success",
        description: "Supplier updated successfully",
      });
      navigate("/app/procurement/suppliers");
    } catch (error) {
      console.error("Failed to update supplier", error);
      toast({
        title: "Error",
        description: "Failed to update supplier. Please try again.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  if (fetching) {
    return <div>Loading...</div>;
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500 max-w-5xl mx-auto pb-10">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Edit Supplier</h2>
        <p className="text-muted-foreground">Update supplier information.</p>
      </div>

      <SupplierForm
        initialValues={initialValues}
        onSubmit={handleSubmit}
        isLoading={loading}
        onCancel={() => navigate("/app/procurement/suppliers")}
      />
    </div>
  );
}
