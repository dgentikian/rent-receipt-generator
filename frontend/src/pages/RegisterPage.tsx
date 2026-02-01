import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Input } from '../components/common/Input';
import { Button } from '../components/common/Button';
import { Card } from '../components/common/Card';
import { authService } from '../services/auth';

export const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    confirmPassword: '',
    first_name: '',
    last_name: '',
    address: '',
    city: '',
    postal_code: '',
    phone: '',
  });
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  
  const navigate = useNavigate();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    // Validate passwords match
    if (formData.password !== formData.confirmPassword) {
      setError('Les mots de passe ne correspondent pas');
      return;
    }

    // Validate password length
    if (formData.password.length < 6) {
      setError('Le mot de passe doit contenir au moins 6 caractères');
      return;
    }

    setIsLoading(true);

    try {
      await authService.register({
        email: formData.email,
        password: formData.password,
        first_name: formData.first_name,
        last_name: formData.last_name,
        address: formData.address,
        city: formData.city,
        postal_code: formData.postal_code,
        phone: formData.phone,
      });
      
      // After successful registration, redirect to login
      alert('Compte créé avec succès ! Veuillez vous connecter.');
      navigate('/login');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Échec de la création du compte');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <Card className="w-full max-w-2xl">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            Créer un compte
          </h1>
          <p className="text-gray-600 mt-2">Inscrivez-vous pour commencer</p>
        </div>

        {error && (
          <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              name="first_name"
              label="Prénom *"
              value={formData.first_name}
              onChange={handleChange}
              required
            />

            <Input
              name="last_name"
              label="Nom *"
              value={formData.last_name}
              onChange={handleChange}
              required
            />
          </div>

          <Input
            type="email"
            name="email"
            label="Email *"
            value={formData.email}
            onChange={handleChange}
            placeholder="votre@email.com"
            required
          />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              type="password"
              name="password"
              label="Mot de passe *"
              value={formData.password}
              onChange={handleChange}
              placeholder="••••••••"
              required
            />

            <Input
              type="password"
              name="confirmPassword"
              label="Confirmer le mot de passe *"
              value={formData.confirmPassword}
              onChange={handleChange}
              placeholder="••••••••"
              required
            />
          </div>

          <hr className="my-6" />
          
          <p className="text-sm text-gray-600 mb-4">Informations optionnelles</p>

          <Input
            name="address"
            label="Adresse"
            value={formData.address}
            onChange={handleChange}
            placeholder="123 Rue de la Paix"
          />

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Input
              name="city"
              label="Ville"
              value={formData.city}
              onChange={handleChange}
              placeholder="Paris"
            />

            <Input
              name="postal_code"
              label="Code postal"
              value={formData.postal_code}
              onChange={handleChange}
              placeholder="75001"
            />
          </div>

          <Input
            type="tel"
            name="phone"
            label="Téléphone"
            value={formData.phone}
            onChange={handleChange}
            placeholder="+33 1 23 45 67 89"
          />

          <Button
            type="submit"
            className="w-full mt-6"
            isLoading={isLoading}
          >
            Créer mon compte
          </Button>
        </form>

        <div className="mt-6 text-center text-sm text-gray-600">
          Vous avez déjà un compte ?{' '}
          <a href="/login" className="text-primary-500 hover:underline font-semibold">
            Se connecter
          </a>
        </div>
      </Card>
    </div>
  );
};
