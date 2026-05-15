export interface Role {
  id: string;
  name: string;
  display_name: string;
  description: string;
  is_system: boolean;
  priority?: number;
  is_active?: boolean;
  organization_id?: string;
  permissions: string[];
}
