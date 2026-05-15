import { useState, useEffect } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { ArrowLeft, Loader2, Save } from "lucide-react";
import { useNavigate, useParams } from "react-router-dom";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

import { deviceService } from "@/services/device.service";
import { DEVICE_TYPES } from "@/types/device.types";
import { useAuth } from "@/contexts/AuthContext";

const macAddressRegex = /^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$|^[0-9A-Fa-f]{12}$/;

const formSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  mac_address: z
    .string()
    .min(1, "MAC address is required")
    .regex(
      macAddressRegex,
      "Invalid MAC address format. Use XX:XX:XX:XX:XX:XX or XXXXXXXXXXXX",
    ),
  device_type: z.string().min(1, "Device type is required"),
  description: z.string().optional(),
  location: z.string().optional(),
  is_active: z.boolean().default(true),
});

type FormValues = z.infer<typeof formSchema>;

export function DeviceFormPage() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { user } = useAuth();
  const isEditMode = !!id;
  const [isLoading, setIsLoading] = useState(false);
  const [isFetching, setIsFetching] = useState(isEditMode);

  const form = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      mac_address: "",
      device_type: "",
      description: "",
      location: "",
      is_active: true,
    },
  });

  useEffect(() => {
    if (isEditMode && id) {
      fetchDevice(id);
    }
  }, [isEditMode, id]);

  const fetchDevice = async (deviceId: string) => {
    try {
      setIsFetching(true);
      const device = await deviceService.getDevice(deviceId);

      form.reset({
        name: device.name,
        mac_address: device.mac_address,
        device_type: device.device_type,
        description: device.description || "",
        location: device.location || "",
        is_active: device.is_active,
      });
    } catch (error) {
      console.error("Failed to fetch device:", error);
      toast.error("Failed to load device details");
      navigate("/app/devices");
    } finally {
      setIsFetching(false);
    }
  };

  const onSubmit = async (values: FormValues) => {
    if (!user?.organization_id) return;

    try {
      setIsLoading(true);

      if (isEditMode && id) {
        await deviceService.updateDevice(id, {
          name: values.name,
          mac_address: values.mac_address,
          device_type: values.device_type as any,
          description: values.description,
          location: values.location,
          is_active: values.is_active,
        });
        toast.success("Device updated successfully");
      } else {
        await deviceService.createDevice({
          organization_id: user.organization_id,
          name: values.name,
          mac_address: values.mac_address,
          device_type: values.device_type as any,
          description: values.description,
          location: values.location,
        });
        toast.success("Device registered successfully");
      }

      navigate("/app/devices");
    } catch (error: any) {
      console.error("Failed to save device:", error);
      toast.error(
        isEditMode ? "Failed to update device" : "Failed to register device",
        {
          description:
            error.response?.data?.message || "Please try again later",
        },
      );
    } finally {
      setIsLoading(false);
    }
  };

  if (isFetching) {
    return (
      <div className="flex h-[400px] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return (
    <div className="space-y-6 max-w-4xl mx-auto pb-10">
      {/* Header */}
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          onClick={() => navigate("/app/devices")}
        >
          <ArrowLeft className="h-5 w-5" />
        </Button>
        <div>
          <h2 className="text-2xl font-semibold tracking-tight">
            {isEditMode ? "Edit Device" : "Register New Device"}
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditMode
              ? "Update device details"
              : "Register a new device for your organization. Users can sign up from registered devices without entering an organization ID."}
          </p>
        </div>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
          <Card>
            <CardHeader>
              <CardTitle>Device Details</CardTitle>
              <CardDescription>
                Provide the device name and its unique MAC address for
                identification.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Device Name</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Front Counter POS"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        A friendly name to identify this device.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="mac_address"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>MAC Address</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. AA:BB:CC:DD:EE:FF"
                          className="font-mono"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        The hardware MAC address of the device.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="device_type"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Device Type</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                        value={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select device type" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {DEVICE_TYPES.map((type) => (
                            <SelectItem key={type.value} value={type.value}>
                              {type.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="location"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Physical Location (Optional)</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="e.g. Main Branch - Counter 1"
                          {...field}
                        />
                      </FormControl>
                      <FormDescription>
                        Where this device is physically located.
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Description (Optional)</FormLabel>
                    <FormControl>
                      <Textarea
                        placeholder="Add any notes about this device..."
                        className="resize-none"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {isEditMode && (
                <FormField
                  control={form.control}
                  name="is_active"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-lg border p-4">
                      <div className="space-y-0.5">
                        <FormLabel className="text-base">
                          Active Status
                        </FormLabel>
                        <FormDescription>
                          Inactive devices will not allow user registration.
                        </FormDescription>
                      </div>
                      <FormControl>
                        <Switch
                          checked={field.value}
                          onCheckedChange={field.onChange}
                        />
                      </FormControl>
                    </FormItem>
                  )}
                />
              )}
            </CardContent>
          </Card>

          {/* Submit Button */}
          <div className="flex justify-end gap-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => navigate("/app/devices")}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              className="bg-brand text-brand-foreground hover:bg-brand/90"
              disabled={isLoading}
            >
              {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              <Save className="mr-2 h-4 w-4" />
              {isEditMode ? "Update Device" : "Register Device"}
            </Button>
          </div>
        </form>
      </Form>
    </div>
  );
}
