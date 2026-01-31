import api from './api';
import { LoginRequest, LoginResponse, LandlordCreateRequest, Landlord } from '../types/landlord';

export const authService = {
  async register(data: LandlordCreateRequest): Promise<Landlord> {
    const response = await api.post<Landlord>('/auth/register', data);
    return response.data;
  },

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>('/auth/login', data);
    if (response.data.token) {
      localStorage.setItem('token', response.data.token);
      localStorage.setItem('landlord', JSON.stringify(response.data.landlord));
    }
    return response.data;
  },

  logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('landlord');
    window.location.href = '/login';
  },

  getStoredLandlord(): Landlord | null {
    const landlordStr = localStorage.getItem('landlord');
    return landlordStr ? JSON.parse(landlordStr) : null;
  },

  getToken(): string | null {
    return localStorage.getItem('token');
  },

  isAuthenticated(): boolean {
    return !!this.getToken();
  },
};
