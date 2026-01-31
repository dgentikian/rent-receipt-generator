export interface Tenant {
  id: number;
  property_id: number;
  first_name: string;
  last_name: string;
  email: string;
  phone: string;
  move_in_date: string | null;
  move_out_date: string | null;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface TenantCreateRequest {
  property_id: number;
  first_name: string;
  last_name: string;
  email?: string;
  phone?: string;
  move_in_date?: string;
}

export interface TenantUpdateRequest {
  first_name: string;
  last_name: string;
  email?: string;
  phone?: string;
  move_in_date?: string;
  move_out_date?: string;
  is_active?: boolean;
}
