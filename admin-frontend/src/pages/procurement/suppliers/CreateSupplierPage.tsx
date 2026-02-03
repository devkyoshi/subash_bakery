import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "@/contexts/AuthContext";
import { procurementService } from "@/services/procurement.service";
import { CreateSupplierRequest } from "@/types/procurement.types";
import { SupplierForm } from "./SupplierForm";
import { useToast } from "@/components/ui/use-toast";

export function CreateSupplierPage() {
  const navigate = useNavigate();
  const { user } = useAuth();
  const { toast } = useToast();
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (values: CreateSupplierRequest) => {
    if (!user?.organization_id) return;

    try {
      setLoading(true);
      await procurementService.createSupplier(user.organization_id, values);
      toast({
        title: "Success",
        description: "Supplier created successfully",
      });
      navigate("/app/procurement/suppliers");
    } catch (error) {
      console.error("Failed to create supplier", error);
      toast({
        title: "Error",
        description: "Failed to create supplier. Please try again.",
        variant: "destructive",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="space-y-6 animate-in fade-in duration-500 max-w-5xl mx-auto pb-10">
      <div>
        <h2 className="text-3xl font-bold tracking-tight">Add Supplier</h2>
        <p className="text-muted-foreground">
          Register a new vendor in the system.
        </p>
      </div>

      <SupplierForm
        onSubmit={handleSubmit}
        isLoading={loading}
        onCancel={() => navigate("/app/procurement/suppliers")}
      />
    </div>
  );
}
