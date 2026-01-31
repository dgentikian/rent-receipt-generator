export interface Landlord {
  id: number;
  email: string;
  first_name: string;
  last_name: string;
  address: string;
  city: string;
  postal_code: string;
  phone: string;
  signature_url: string;
  created_at: string;
  updated_at: string;
}

export interface LandlordCreateRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  address?: string;
  city?: string;
  postal_code?: string;
  phone?: string;
}

export interface LandlordUpdateRequest {
  first_name: string;
  last_name: string;
  address: string;
  city: string;
  postal_code: string;
  phone: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  landlord: Landlord;
}
