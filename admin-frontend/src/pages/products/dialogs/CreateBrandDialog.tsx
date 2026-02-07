import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
  DialogDescription,
} from "@/components/ui/dialog";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
  FormDescription,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { brandService } from "@/services/brand.service";
import { useAuth } from "@/contexts/AuthContext";
import { Loader2 } from "lucide-react";

const createBrandSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  code: z
    .string()
    .min(2, "Code must be at least 2 characters")
    .regex(
      /^[A-Z0-9-_]+$/,
      "Code must contain only uppercase letters, numbers, hyphens, and underscores",
    ),
});

type CreateBrandForm = z.infer<typeof createBrandSchema>;

interface CreateBrandDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: (brand: any) => void;
}

export function CreateBrandDialog({
  open,
  onOpenChange,
  onSuccess,
}: CreateBrandDialogProps) {
  const { user } = useAuth();
  const [isLoading, setIsLoading] = useState(false);

  const form = useForm<CreateBrandForm>({
    resolver: zodResolver(createBrandSchema),
    defaultValues: {
      name: "",
      code: "",
    },
  });

  const nameValue = form.watch("name");

  // Auto-generate code from name
  useEffect(() => {
    if (nameValue && !form.formState.dirtyFields.code) {
      const generatedCode = nameValue
        .toUpperCase()
        .replace(/[^A-Z0-9]/g, "-") // Replace non-alphanumeric with hyphen
        .replace(/-+/g, "-") // Replace multiple hyphens with single
        .replace(/^-|-$/g, "") // Trim hyphens
        .slice(0, 20); // Limit length

      form.setValue("code", generatedCode);
    }
  }, [nameValue, form.formState.dirtyFields.code, form]);

  const onSubmit = async (values: CreateBrandForm) => {
    if (!user?.organization_id) return;
    try {
      setIsLoading(true);
      const brand = await brandService.createBrand({
        name: values.name,
        code: values.code,
        organization_id: user.organization_id,
        is_active: true,
      });
      toast.success("Brand created successfully");
      form.reset();
      onOpenChange(false);
      onSuccess?.(brand);
    } catch (error: any) {
      toast.error("Failed to create brand", {
        description: error.response?.data?.message || "Please try again",
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>Create New Brand</DialogTitle>
          <DialogDescription>
            Add a new brand or manufacturer to your catalog.
          </DialogDescription>
        </DialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    Name <span className="text-destructive">*</span>
                  </FormLabel>
                  <FormControl>
                    <Input placeholder="Brand Name" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="code"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>
                    Code <span className="text-destructive">*</span>
                  </FormLabel>
                  <FormControl>
                    <Input placeholder="e.g. BRAND-CODE" {...field} />
                  </FormControl>
                  <FormDescription className="text-xs">
                    Unique identifier for system use.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => onOpenChange(false)}
                disabled={isLoading}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={isLoading}>
                {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create Brand
              </Button>
            </DialogFooter>
          </form>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
