import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { propertyService } from '../services/property';
import { PropertyCreateRequest } from '../types/property';

export const PropertiesPage: React.FC = () => {
  const [showForm, setShowForm] = useState(false);
  const queryClient = useQueryClient();

  const { data: properties = [], isLoading } = useQuery({
    queryKey: ['properties'],
    queryFn: propertyService.getAll,
  });

  const createMutation = useMutation({
    mutationFn: propertyService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
      setShowForm(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: propertyService.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['properties'] });
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: PropertyCreateRequest = {
      address: formData.get('address') as string,
      city: formData.get('city') as string,
      postal_code: formData.get('postal_code') as string,
      rent_amount: parseFloat(formData.get('rent_amount') as string),
      charges_amount: parseFloat(formData.get('charges_amount') as string) || 0,
      syndic_name: formData.get('syndic_name') as string,
      syndic_address: formData.get('syndic_address') as string,
    };
    createMutation.mutate(data);
  };

  return (
    <Layout>
      <div className="space-y-6">
        <div className="flex justify-between items-center">
          <h1 className="text-3xl font-bold text-gray-800">Mes Propriétés</h1>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? 'Annuler' : '+ Nouvelle propriété'}
          </Button>
        </div>

        {showForm && (
          <Card>
            <h2 className="text-xl font-bold mb-4">Ajouter une propriété</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <Input name="address" label="Adresse" required />
              <div className="grid grid-cols-2 gap-4">
                <Input name="city" label="Ville" required />
                <Input name="postal_code" label="Code postal" required />
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
                />
              </div>
              <Input name="syndic_name" label="Nom du syndic" />
              <Input name="syndic_address" label="Adresse du syndic" />
              <Button type="submit" isLoading={createMutation.isPending}>
                Ajouter
              </Button>
            </form>
          </Card>
        )}

        {isLoading ? (
          <p>Chargement...</p>
        ) : properties.length === 0 ? (
          <Card>
            <p className="text-gray-500 text-center py-8">
              Aucune propriété enregistrée
            </p>
          </Card>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {properties.map((property) => (
              <Card key={property.id} hover>
                <div className="space-y-3">
                  <div className="flex justify-between items-start">
                    <div>
                      <h3 className="font-bold text-lg">{property.address}</h3>
                      <p className="text-gray-600">
                        {property.postal_code} {property.city}
                      </p>
                    </div>
                    <Button
                      variant="danger"
                      onClick={() => deleteMutation.mutate(property.id)}
                      className="text-sm px-3 py-1"
                    >
                      Supprimer
                    </Button>
                  </div>
                  <div className="pt-3 border-t">
                    <div className="flex justify-between">
                      <span className="text-gray-600">Loyer:</span>
                      <span className="font-bold">{property.rent_amount} €</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-gray-600">Charges:</span>
                      <span className="font-bold">{property.charges_amount} €</span>
                    </div>
                    <div className="flex justify-between pt-2 border-t mt-2">
                      <span className="text-gray-600">Total:</span>
                      <span className="font-bold text-primary-600">
                        {(property.rent_amount + property.charges_amount).toFixed(2)} €
                      </span>
                    </div>
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
