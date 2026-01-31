export interface Property {
  id: number;
  landlord_id: number;
  address: string;
  city: string;
  postal_code: string;
  rent_amount: number;
  charges_amount: number;
  syndic_name: string;
  syndic_address: string;
  created_at: string;
  updated_at: string;
}

export interface PropertyCreateRequest {
  address: string;
  city: string;
  postal_code: string;
  rent_amount: number;
  charges_amount?: number;
  syndic_name?: string;
  syndic_address?: string;
}

export interface PropertyUpdateRequest extends PropertyCreateRequest {}
