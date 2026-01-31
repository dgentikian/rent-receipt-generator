import api from './api';
import { Landlord, LandlordUpdateRequest } from '../types/landlord';

export const landlordService = {
  async getProfile(): Promise<Landlord> {
    const response = await api.get<Landlord>('/landlord/profile');
    return response.data;
  },

  async updateProfile(data: LandlordUpdateRequest): Promise<Landlord> {
    const response = await api.put<Landlord>('/landlord/profile', data);
    // Update stored landlord
    localStorage.setItem('landlord', JSON.stringify(response.data));
    return response.data;
  },

  async uploadSignature(file: File): Promise<{ signature_url: string }> {
    const formData = new FormData();
    formData.append('signature', file);
    const response = await api.post<{ signature_url: string }>('/landlord/signature', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },
};
