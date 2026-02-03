import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { procurementService } from "@/services/procurement.service";
import { GoodsReceiptNote, GRNStatus } from "@/types/procurement.types";
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
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { ArrowLeft, CalendarDays, ClipboardCheck, User } from "lucide-react";
import { format } from "date-fns";
import { Separator } from "@/components/ui/separator";
import { useToast } from "@/components/ui/use-toast";

const inspectionSchema = z.object({
  qc_status: z.enum(["passed", "failed"]),
  qc_notes: z.string().optional(),
});

type InspectionValues = z.infer<typeof inspectionSchema>;

export function GRNDetailsPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();

  const [grn, setGrn] = useState<GoodsReceiptNote | null>(null);
  const [loading, setLoading] = useState(true);
  const [inspectionOpen, setInspectionOpen] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const form = useForm<InspectionValues>({
    resolver: zodResolver(inspectionSchema),
    defaultValues: {
      qc_status: "passed",
      qc_notes: "",
    },
  });

  useEffect(() => {
    if (id) {
      fetchGRN(id);
    }
  }, [id]);

  const fetchGRN = async (grnId: string) => {
    try {
      setLoading(true);
      const data = await procurementService.getGRN(grnId);
      setGrn(data);
    } catch (error) {
      console.error("Failed to fetch GRN details", error);
    } finally {
      setLoading(false);
    }
  };

  const onInspectSubmit = async (values: InspectionValues) => {
    if (!grn?.id) return;
    try {
      setSubmitting(true);
      await procurementService.inspectGRN(grn.id, {
        qc_status: values.qc_status,
        qc_notes: values.qc_notes,
      });
      toast({
        title: "Inspection Completed",
        description: `GRN marked as ${values.qc_status}`,
      });
      setInspectionOpen(false);
      fetchGRN(grn.id); // Refresh
    } catch (error) {
      console.error("Failed to submit inspection", error);
      toast({
        title: "Error",
        description: "Failed to submit inspection",
        variant: "destructive",
      });
    } finally {
      setSubmitting(false);
    }
  };

  const getStatusBadgeVariant = (status: GRNStatus) => {
    switch (status) {
      case GRNStatus.Draft:
        return "secondary";
      case GRNStatus.Received:
        return "default";
      case GRNStatus.Inspected:
        return "outline";
      case GRNStatus.Accepted:
        return "default";
      case GRNStatus.Rejected:
        return "destructive";
      default:
        return "outline";
    }
  };

  if (loading) {
    return <div className="p-8 text-center">Loading GRN details...</div>;
  }

  if (!grn) {
    return (
      <div className="p-8 text-center">
        <h3 className="text-lg font-medium text-destructive">GRN not found</h3>
        <Button
          variant="outline"
          onClick={() => navigate("/app/procurement/grn")}
          className="mt-4"
        >
          <ArrowLeft className="mr-2 h-4 w-4" /> Back to GRNs
        </Button>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => navigate("/app/procurement/grn")}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <div>
            <h2 className="text-3xl font-bold tracking-tight">
              {grn.grn_number}
            </h2>
            <div className="flex items-center gap-2 mt-1">
              <Badge variant={getStatusBadgeVariant(grn.status)}>
                {grn.status?.toUpperCase()}
              </Badge>
              <span className="text-sm text-muted-foreground flex items-center gap-1 ml-2">
                <CalendarDays className="h-3 w-3" />
                {format(new Date(grn.receipt_date), "PPP")}
              </span>
            </div>
          </div>
        </div>

        <div>
          {/* Allow inspection if status is RECEIVED (or whatever flow dictates) */}
          {grn.status === GRNStatus.Received && (
            <Dialog open={inspectionOpen} onOpenChange={setInspectionOpen}>
              <DialogTrigger asChild>
                <Button>
                  <ClipboardCheck className="mr-2 h-4 w-4" /> Perform QC
                  Inspection
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>QC Inspection</DialogTitle>
                  <DialogDescription>
                    Record the Quality Control results for this delivery.
                  </DialogDescription>
                </DialogHeader>
                <Form {...form}>
                  <form
                    onSubmit={form.handleSubmit(onInspectSubmit)}
                    className="space-y-4"
                  >
                    <FormField
                      control={form.control}
                      name="qc_status"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>QC Outcome</FormLabel>
                          <Select
                            onValueChange={field.onChange}
                            defaultValue={field.value}
                          >
                            <FormControl>
                              <SelectTrigger>
                                <SelectValue placeholder="Select outcome" />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="passed">Passed</SelectItem>
                              <SelectItem value="failed">Failed</SelectItem>
                            </SelectContent>
                          </Select>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="qc_notes"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Notes</FormLabel>
                          <FormControl>
                            <Textarea
                              placeholder="Any issues found..."
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <DialogFooter>
                      <Button
                        type="button"
                        variant="outline"
                        onClick={() => setInspectionOpen(false)}
                      >
                        Cancel
                      </Button>
                      <Button type="submit" disabled={submitting}>
                        Submit Inspection
                      </Button>
                    </DialogFooter>
                  </form>
                </Form>
              </DialogContent>
            </Dialog>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Received Items</CardTitle>
            </CardHeader>
            <CardContent className="p-0">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Product</TableHead>
                    <TableHead className="text-right">Ordered</TableHead>
                    <TableHead className="text-right">Received</TableHead>
                    <TableHead>Batch Info</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {grn.items.map((item, index) => (
                    <TableRow key={index}>
                      <TableCell>
                        <div className="font-medium">
                          {item.description ||
                            `Product ${item.product_id.slice(-6)}`}
                        </div>
                        <div className="text-xs text-muted-foreground flex gap-2">
                          <span>SKU: {item.sku || "-"}</span>
                          {item.condition && (
                            <Badge
                              variant="outline"
                              className={
                                item.condition === "good"
                                  ? "text-green-600 border-green-200 bg-green-50"
                                  : item.condition === "damaged"
                                    ? "text-red-600 border-red-200 bg-red-50"
                                    : "text-amber-600 border-amber-200 bg-amber-50"
                              }
                            >
                              {item.condition.toUpperCase()}
                            </Badge>
                          )}
                        </div>
                      </TableCell>
                      <TableCell className="text-right">
                        {item.ordered_quantity ?? "-"}
                      </TableCell>
                      <TableCell className="text-right font-bold">
                        {item.received_quantity}
                      </TableCell>
                      <TableCell>
                        {item.batch_number ? (
                          <div className="text-xs">
                            <div>Batch: {item.batch_number}</div>
                            {item.expiry_date && (
                              <div>Exp: {item.expiry_date}</div>
                            )}
                          </div>
                        ) : (
                          "-"
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>

          {grn.notes && (
            <Card>
              <CardHeader>
                <CardTitle>Notes</CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {grn.notes}
                </p>
              </CardContent>
            </Card>
          )}

          {grn.qc_notes && (
            <Card>
              <CardHeader>
                <CardTitle>QC Inspection Notes</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex gap-2 mb-2">
                  <Badge variant="outline">{grn.qc_status}</Badge>
                </div>
                <p className="text-sm text-muted-foreground whitespace-pre-wrap">
                  {grn.qc_notes}
                </p>
              </CardContent>
            </Card>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>Info</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <span className="text-sm text-muted-foreground">PO Number</span>
                <div className="font-medium">
                  {grn.po_number || grn.purchase_order_id.slice(-6)}
                </div>
                {/* Link to PO */}
                <Button
                  variant="link"
                  className="p-0 h-auto"
                  onClick={() =>
                    navigate(`/app/procurement/orders/${grn.purchase_order_id}`)
                  }
                >
                  View Purchase Order
                </Button>
              </div>
              <Separator />
              <div>
                <span className="text-sm text-muted-foreground flex items-center gap-1">
                  <User className="h-3 w-3" /> Received By
                </span>
                <div className="font-medium">
                  {grn.received_by_name || grn.received_by}
                </div>
                <div className="text-xs text-muted-foreground">
                  {grn.receipt_date &&
                    format(new Date(grn.receipt_date), "PPP 'at' p")}
                </div>
              </div>
              {grn.inspected_by && (
                <div>
                  <span className="text-sm text-muted-foreground flex items-center gap-1">
                    <ClipboardCheck className="h-3 w-3" /> Inspected By
                  </span>
                  <div className="font-medium">
                    {grn.inspected_by_name || grn.inspected_by}
                  </div>
                  {grn.inspected_date && (
                    <div className="text-xs text-muted-foreground">
                      {format(new Date(grn.inspected_date), "PPP 'at' p")}
                    </div>
                  )}
                </div>
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
