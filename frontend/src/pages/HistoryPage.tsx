import React, { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { receiptService } from '../services/receipt';
import { propertyService } from '../services/property';

export const HistoryPage: React.FC = () => {
  const currentYear = new Date().getFullYear();
  const [selectedYear, setSelectedYear] = useState(currentYear);
  const [selectedPropertyId, setSelectedPropertyId] = useState<number | undefined>();

  const { data: receipts = [], isLoading } = useQuery({
    queryKey: ['receipts-history', selectedYear, selectedPropertyId],
    queryFn: () =>
      receiptService.getAll({
        year: selectedYear,
        property_id: selectedPropertyId,
        limit: 200,
      }),
  });

  const { data: properties = [] } = useQuery({
    queryKey: ['properties'],
    queryFn: propertyService.getAll,
  });

  const years = Array.from({ length: 5 }, (_, i) => currentYear - i);

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

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-800">Historique des Quittances</h1>

        <Card>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block mb-2 font-semibold text-gray-700">Année</label>
              <select
                className="input"
                value={selectedYear}
                onChange={(e) => setSelectedYear(parseInt(e.target.value))}
              >
                {years.map((year) => (
                  <option key={year} value={year}>
                    {year}
                  </option>
                ))}
              </select>
            </div>

            <div>
              <label className="block mb-2 font-semibold text-gray-700">
                Propriété
              </label>
              <select
                className="input"
                value={selectedPropertyId || ''}
                onChange={(e) =>
                  setSelectedPropertyId(e.target.value ? parseInt(e.target.value) : undefined)
                }
              >
                <option value="">Toutes</option>
                {properties.map((p) => (
                  <option key={p.id} value={p.id}>
                    {p.address}
                  </option>
                ))}
              </select>
            </div>
          </div>
        </Card>

        {isLoading ? (
          <Card>
            <p className="text-center py-8">Chargement...</p>
          </Card>
        ) : receipts.length === 0 ? (
          <Card>
            <p className="text-gray-500 text-center py-8">
              Aucune quittance trouvée pour les filtres sélectionnés
            </p>
          </Card>
        ) : (
          <div className="space-y-4">
            {receipts.map((receipt) => (
              <Card key={receipt.id}>
                <div className="flex justify-between items-center">
                  <div className="flex-1">
                    <div className="flex items-center gap-4">
                      <div className="text-3xl font-bold text-primary-500">
                        {String(receipt.period_month).padStart(2, '0')}
                      </div>
                      <div>
                        <h3 className="font-bold">{receipt.receipt_number}</h3>
                        <p className="text-sm text-gray-600">
                          {receipt.period_month}/{receipt.period_year}
                        </p>
                        <p className="text-xs text-gray-500">
                          Créé: {new Date(receipt.created_at).toLocaleDateString('fr-FR')}
                        </p>
                      </div>
                    </div>
                  </div>
                  <div className="text-right space-y-2">
                    <p className="text-xl font-bold text-primary-600">
                      {receipt.total_amount.toFixed(2)} €
                    </p>
                    <Button
                      variant="secondary"
                      onClick={() => handleDownload(receipt.id)}
                      className="text-sm"
                    >
                      📥 Télécharger
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
