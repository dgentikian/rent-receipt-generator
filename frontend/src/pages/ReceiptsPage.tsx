import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { receiptService } from '../services/receipt';
import { propertyService } from '../services/property';
import { tenantService } from '../services/tenant';
import { ReceiptCreateRequest } from '../types/receipt';

export const ReceiptsPage: React.FC = () => {
  const [showForm, setShowForm] = useState(false);
  const [selectedPropertyId, setSelectedPropertyId] = useState<number | null>(null);
  const queryClient = useQueryClient();

  const { data: receipts = [], isLoading } = useQuery({
    queryKey: ['receipts'],
    queryFn: () => receiptService.getAll({ limit: 100 }),
  });

  const { data: properties = [] } = useQuery({
    queryKey: ['properties'],
    queryFn: propertyService.getAll,
  });

  const { data: tenants = [] } = useQuery({
    queryKey: ['tenants', selectedPropertyId],
    queryFn: () => tenantService.getByProperty(selectedPropertyId!),
    enabled: !!selectedPropertyId,
  });

  const createMutation = useMutation({
    mutationFn: receiptService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['receipts'] });
      setShowForm(false);
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: ReceiptCreateRequest = {
      property_id: parseInt(formData.get('property_id') as string),
      tenant_id: parseInt(formData.get('tenant_id') as string),
      period_month: parseInt(formData.get('period_month') as string),
      period_year: parseInt(formData.get('period_year') as string),
      rent_amount: parseFloat(formData.get('rent_amount') as string),
      charges_amount: parseFloat(formData.get('charges_amount') as string) || 0,
      payment_method: formData.get('payment_method') as string,
      payment_date: formData.get('payment_date') as string,
    };
    createMutation.mutate(data);
  };

  const handleDownload = async (id: number) => {
    try {
      const blob = await receiptService.downloadPDF(id);
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `quittance_${id}.pdf`;
      a.click();
      window.URL.revokeObjectURL(url);
    } catch (error) {
      alert('Erreur lors du téléchargement');
    }
  };

  const currentYear = new Date().getFullYear();
  const currentMonth = new Date().getMonth() + 1;

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-800">Quittances</h1>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Annuler' : '+ Générer une quittance'}
          </Button>
        </div>

        {showForm && (
          <Card>
            <h2 className="text-xl font-bold mb-4">Générer une quittance</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block mb-2 font-semibold text-gray-700">
                  Propriété *
                </label>
                <select
                  name="property_id"
                  className="input"
                  required
                  onChange={(e) => setSelectedPropertyId(parseInt(e.target.value))}
                >
                  <option value="">Sélectionner...</option>
                  {properties.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.address}
                    </option>
                  ))}
                </select>
              </div>

              {selectedPropertyId && (
                <div>
                  <label className="block mb-2 font-semibold text-gray-700">
                    Locataire *
                  </label>
                  <select name="tenant_id" className="input" required>
                    <option value="">Sélectionner...</option>
                    {tenants.filter(t => t.is_active).map((t) => (
                      <option key={t.id} value={t.id}>
                        {t.first_name} {t.last_name}
                      </option>
                    ))}
                  </select>
                </div>
              )}

              <div className="grid grid-cols-2 gap-4">
                <Input
                  name="period_month"
                  type="number"
                  min="1"
                  max="12"
                  label="Mois"
                  defaultValue={currentMonth}
                  required
                />
                <Input
                  name="period_year"
                  type="number"
                  label="Année"
                  defaultValue={currentYear}
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Input
                  name="rent_amount"
                  type="number"
                  step="0.01"
                  label="Loyer (€)"
                  required
                />
                <Input
                  name="charges_amount"
                  type="number"
                  step="0.01"
                  label="Charges (€)"
                  defaultValue="0"
                />
              </div>

              <Input name="payment_method" label="Moyen de paiement" />
              <Input name="payment_date" type="date" label="Date de paiement" />

              <Button type="submit" isLoading={createMutation.isPending}>
                Générer la quittance
              </Button>
            </form>
          </Card>
        )}

        {isLoading ? (
          <p>Chargement...</p>
        ) : receipts.length === 0 ? (
          <Card>
            <p className="text-gray-500 text-center py-8">
              Aucune quittance générée
            </p>
          </Card>
        ) : (
          <div className="space-y-4">
            {receipts.map((receipt) => (
              <Card key={receipt.id}>
                <div className="flex justify-between items-center">
                  <div>
                    <h3 className="font-bold text-lg">{receipt.receipt_number}</h3>
                    <p className="text-gray-600">
                      Période: {receipt.period_month}/{receipt.period_year}
                    </p>
                    <p className="text-sm text-gray-500">
                      Créé le: {new Date(receipt.created_at).toLocaleDateString('fr-FR')}
                    </p>
                  </div>
                  <div className="text-right space-y-2">
                    <p className="text-2xl font-bold text-primary-600">
                      {receipt.total_amount.toFixed(2)} €
                    </p>
                    <Button
                      variant="secondary"
                      onClick={() => handleDownload(receipt.id)}
                      className="text-sm"
                    >
                      📥 Télécharger PDF
                    </Button>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>
    </Layout>
  );
};
