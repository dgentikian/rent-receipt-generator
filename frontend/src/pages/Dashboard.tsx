import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { propertyService } from '../services/property';
import { receiptService } from '../services/receipt';
import { useAuth } from '../context/AuthContext';

export const Dashboard: React.FC = () => {
  const { landlord } = useAuth();

  const { data: properties = [] } = useQuery({
    queryKey: ['properties'],
    queryFn: propertyService.getAll,
  });

  const { data: receipts = [] } = useQuery({
    queryKey: ['receipts-recent'],
    queryFn: () => receiptService.getAll({ limit: 10 }),
  });

  const totalProperties = properties.length;
  const totalRent = properties.reduce((sum, p) => sum + (p.rent_amount || 0), 0);
  const receiptsThisMonth = receipts.filter(r => {
    const date = new Date();
    return r.period_month === date.getMonth() + 1 && r.period_year === date.getFullYear();
  }).length;

  return (
    <Layout>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800">
            Bienvenue, {landlord?.first_name} !
          </h1>
          <p className="text-gray-600 mt-2">
            Voici un aperçu de votre activité
          </p>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">Propriétés</p>
                <p className="text-3xl font-bold text-primary-600 mt-2">
                  {totalProperties}
                </p>
              </div>
              <div className="text-5xl">🏠</div>
            </div>
          </Card>

          <Card>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">Loyer total mensuel</p>
                <p className="text-3xl font-bold text-green-600 mt-2">
                  {totalRent.toFixed(2)} €
                </p>
              </div>
              <div className="text-5xl">💰</div>
            </div>
          </Card>

          <Card>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-gray-600 text-sm">Quittances ce mois</p>
                <p className="text-3xl font-bold text-blue-600 mt-2">
                  {receiptsThisMonth}
                </p>
              </div>
              <div className="text-5xl">📄</div>
            </div>
          </Card>
        </div>

        {/* Recent Receipts */}
        <Card>
          <h2 className="text-xl font-bold mb-4">Dernières quittances</h2>
          {receipts.length === 0 ? (
            <p className="text-gray-500">Aucune quittance générée</p>
          ) : (
            <div className="space-y-3">
              {receipts.slice(0, 5).map((receipt) => (
                <div
                  key={receipt.id}
                  className="flex justify-between items-center p-3 bg-gray-50 rounded-lg"
                >
                  <div>
                    <p className="font-semibold">{receipt.receipt_number}</p>
                    <p className="text-sm text-gray-600">
                      {receipt.period_month}/{receipt.period_year}
                    </p>
                  </div>
                  <div className="text-right">
                    <p className="font-bold text-primary-600">
                      {receipt.total_amount.toFixed(2)} €
                    </p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </Layout>
  );
};
