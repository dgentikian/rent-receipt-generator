import api from './api';
import { Tenant, TenantCreateRequest, TenantUpdateRequest } from '../types/tenant';

export const tenantService = {
  async create(data: TenantCreateRequest): Promise<Tenant> {
    const response = await api.post<Tenant>('/tenants', data);
    return response.data;
  },

  async getByProperty(propertyId: number): Promise<Tenant[]> {
    const response = await api.get<Tenant[]>(`/tenants?property_id=${propertyId}`);
    return response.data;
  },

  async getById(id: number): Promise<Tenant> {
    const response = await api.get<Tenant>(`/tenants/${id}`);
    return response.data;
  },

  async update(id: number, data: TenantUpdateRequest): Promise<Tenant> {
    const response = await api.put<Tenant>(`/tenants/${id}`, data);
    return response.data;
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/tenants/${id}`);
  },
};
