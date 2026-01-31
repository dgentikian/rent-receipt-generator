import React from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Layout } from '../components/layout/Layout';
import { Card } from '../components/common/Card';
import { Button } from '../components/common/Button';
import { Input } from '../components/common/Input';
import { useAuth } from '../context/AuthContext';
import { landlordService } from '../services/landlord';
import { LandlordUpdateRequest } from '../types/landlord';

export const LandlordPage: React.FC = () => {
  const { landlord, updateLandlord } = useAuth();
  const queryClient = useQueryClient();

  const updateMutation = useMutation({
    mutationFn: (data: LandlordUpdateRequest) => landlordService.updateProfile(data),
    onSuccess: (data) => {
      updateLandlord(data);
      alert('Profil mis à jour avec succès');
    },
  });

  const uploadSignatureMutation = useMutation({
    mutationFn: (file: File) => landlordService.uploadSignature(file),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['landlord'] });
      alert('Signature ajoutée avec succès');
    },
  });

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const data: LandlordUpdateRequest = {
      first_name: formData.get('first_name') as string,
      last_name: formData.get('last_name') as string,
      address: formData.get('address') as string,
      city: formData.get('city') as string,
      postal_code: formData.get('postal_code') as string,
      phone: formData.get('phone') as string,
    };
    updateMutation.mutate(data);
  };

  const handleSignatureUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      uploadSignatureMutation.mutate(file);
    }
  };

  if (!landlord) return null;

  return (
    <Layout>
      <div className="space-y-6">
        <h1 className="text-3xl font-bold text-gray-800">Mes Informations</h1>

        <Card>
          <h2 className="text-xl font-bold mb-4">Informations personnelles</h2>
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <Input
                name="first_name"
                label="Prénom"
                defaultValue={landlord.first_name}
                required
              />
              <Input
                name="last_name"
                label="Nom"
                defaultValue={landlord.last_name}
                required
              />
            </div>

            <Input
              name="email"
              type="email"
              label="Email"
              value={landlord.email}
              disabled
            />

            <Input
              name="address"
              label="Adresse"
              defaultValue={landlord.address}
            />

            <div className="grid grid-cols-2 gap-4">
              <Input
                name="city"
                label="Ville"
                defaultValue={landlord.city}
              />
              <Input
                name="postal_code"
                label="Code postal"
                defaultValue={landlord.postal_code}
              />
            </div>

            <Input
              name="phone"
              type="tel"
              label="Téléphone"
              defaultValue={landlord.phone}
            />

            <Button type="submit" isLoading={updateMutation.isPending}>
              Mettre à jour
            </Button>
          </form>
        </Card>

        <Card>
          <h2 className="text-xl font-bold mb-4">Signature</h2>
          {landlord.signature_url && (
            <div className="mb-4">
              <p className="text-sm text-gray-600 mb-2">Signature actuelle:</p>
              <img
                src={landlord.signature_url}
                alt="Signature"
                className="max-w-xs border rounded"
              />
            </div>
          )}
          <div>
            <label className="block mb-2 font-semibold text-gray-700">
              Télécharger une nouvelle signature
            </label>
            <input
              type="file"
              accept="image/*"
              onChange={handleSignatureUpload}
              className="input"
            />
          </div>
        </Card>
      </div>
    </Layout>
  );
};
