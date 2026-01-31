import api from './api';
import { Receipt, ReceiptWithDetails, ReceiptCreateRequest, ReceiptListQuery } from '../types/receipt';

export const receiptService = {
  async create(data: ReceiptCreateRequest): Promise<Receipt> {
    const response = await api.post<Receipt>('/receipts', data);
    return response.data;
  },

  async getAll(query?: ReceiptListQuery): Promise<Receipt[]> {
    const params = new URLSearchParams();
    if (query?.property_id) params.append('property_id', query.property_id.toString());
    if (query?.tenant_id) params.append('tenant_id', query.tenant_id.toString());
    if (query?.year) params.append('year', query.year.toString());
    if (query?.month) params.append('month', query.month.toString());
    if (query?.limit) params.append('limit', query.limit.toString());
    if (query?.offset) params.append('offset', query.offset.toString());

    const response = await api.get<Receipt[]>(`/receipts?${params.toString()}`);
    return response.data;
  },

  async getById(id: number): Promise<Receipt> {
    const response = await api.get<Receipt>(`/receipts/${id}`);
    return response.data;
  },

  async getWithDetails(id: number): Promise<ReceiptWithDetails> {
    const response = await api.get<ReceiptWithDetails>(`/receipts/${id}/details`);
    return response.data;
  },

  async downloadPDF(id: number): Promise<Blob> {
    const response = await api.get(`/receipts/${id}/pdf`, {
      responseType: 'blob',
    });
    return response.data;
  },
};
