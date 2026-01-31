import api from './api';
import { Property, PropertyCreateRequest, PropertyUpdateRequest } from '../types/property';

export const propertyService = {
  async create(data: PropertyCreateRequest): Promise<Property> {
    const response = await api.post<Property>('/properties', data);
    return response.data;
  },

  async getAll(): Promise<Property[]> {
    const response = await api.get<Property[]>('/properties');
    return response.data;
  },

  async getById(id: number): Promise<Property> {
    const response = await api.get<Property>(`/properties/${id}`);
    return response.data;
  },

  async update(id: number, data: PropertyUpdateRequest): Promise<Property> {
    const response = await api.put<Property>(`/properties/${id}`, data);
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/properties/${id}`);
  },
};
