import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { tenantService } from '../services/tenant';
import { propertyService } from '../services/property';
import { TenantCreateRequest } from '../types/tenant';

export const TenantsPage: React.FC = () => {
  const [showForm, setShowForm] = useState(false);
  const [selectedPropertyId, setSelectedPropertyId] = useState<number | null>(null);
  const queryClient = useQueryClient();

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
    mutationFn: tenantService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
      setShowForm(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: tenantService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tenants'] });
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: TenantCreateRequest = {
      property_id: parseInt(formData.get('property_id') as string),
      first_name: formData.get('first_name') as string,
      last_name: formData.get('last_name') as string,
      email: formData.get('email') as string,
      phone: formData.get('phone') as string,
      move_in_date: formData.get('move_in_date') as string,
    };
    createMutation.mutate(data);
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-800">Locataires</h1>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Annuler' : '+ Nouveau locataire'}
          </Button>
        </div>

        {showForm && (
          <Card>
            <h2 className="text-xl font-bold mb-4">Ajouter un locataire</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block mb-2 font-semibold text-gray-700">
                  Propriété *
                </label>
                <select name="property_id" className="input" required>
                  <option value="">Sélectionner...</option>
                  {properties.map((p) => (
                    <option key={p.id} value={p.id}>
                      {p.address}
                    </option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Input name="first_name" label="Prénom" required />
                <Input name="last_name" label="Nom" required />
              </div>

              <Input name="email" type="email" label="Email" />
              <Input name="phone" type="tel" label="Téléphone" />
              <Input name="move_in_date" type="date" label="Date d'entrée" />

              <Button type="submit" isLoading={createMutation.isPending}>
                Ajouter
              </Button>
            </form>
          </Card>
        )}

        <Card>
          <div className="mb-4">
            <label className="block mb-2 font-semibold text-gray-700">
              Filtrer par propriété
            </label>
            <select
              className="input"
              value={selectedPropertyId || ''}
              onChange={(e) =>
                setSelectedPropertyId(e.target.value ? parseInt(e.target.value) : null)
              }
            >
              <option value="">Toutes les propriétés</option>
              {properties.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.address}
                </option>
              ))}
            </select>
          </div>

          {!selectedPropertyId ? (
            <p className="text-gray-500 text-center py-8">
              Veuillez sélectionner une propriété
            </p>
          ) : tenants.length === 0 ? (
            <p className="text-gray-500 text-center py-8">
              Aucun locataire pour cette propriété
            </p>
          ) : (
            <div className="space-y-4">
              {tenants.map((tenant) => (
                <div
                  key={tenant.id}
                  className="flex justify-between items-center p-4 bg-gray-50 rounded-lg"
                >
                  <div>
                    <h3 className="font-bold text-lg">
                      {tenant.first_name} {tenant.last_name}
                      {tenant.is_active && (
                        <span className="ml-2 text-xs bg-green-100 text-green-700 px-2 py-1 rounded">
                          Actif
                        </span>
                      )}
                    </h3>
                    {tenant.email && (
                      <p className="text-sm text-gray-600">📧 {tenant.email}</p>
                    )}
                    {tenant.phone && (
                      <p className="text-sm text-gray-600">📱 {tenant.phone}</p>
                    )}
                    {tenant.move_in_date && (
                      <p className="text-sm text-gray-600">
                        Entrée: {new Date(tenant.move_in_date).toLocaleDateString('fr-FR')}
                      </p>
                    )}
                  </div>
                  <Button
                    variant="danger"
                    onClick={() => deleteMutation.mutate(tenant.id)}
                    className="text-sm"
                  >
                    Supprimer
                  </Button>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </Layout>
  );
};
