import { Landlord } from './landlord';
import { Property } from './property';
import { Tenant } from './tenant';

export interface Receipt {
  id: number;
  landlord_id: number;
  property_id: number;
  tenant_id: number;
  receipt_number: string;
  period_month: number;
  period_year: number;
  rent_amount: number;
  charges_amount: number;
  total_amount: number;
  payment_method: string;
  payment_date: string | null;
  notes: string;
  pdf_url: string;
  created_at: string;
}

export interface ReceiptWithDetails extends Receipt {
  landlord: Landlord;
  property: Property;
  tenant: Tenant;
}

export interface ReceiptCreateRequest {
  property_id: number;
  tenant_id: number;
  period_month: number;
  period_year: number;
  rent_amount: number;
  charges_amount?: number;
  payment_method?: string;
  payment_date?: string;
  notes?: string;
}

export interface ReceiptListQuery {
  property_id?: number;
  tenant_id?: number;
  year?: number;
  month?: number;
  limit?: number;
  offset?: number;
}
