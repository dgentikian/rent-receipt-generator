import React from 'react';
import { useAuth } from '../../context/AuthContext';
import { Button } from '../common/Button';

export const Header: React.FC = () => {
  const { landlord, logout } = useAuth();

  return (
    <header className="bg-gradient-to-r from-primary-500 to-secondary-500 text-white px-8 py-4 shadow-lg">
      <div className="max-w-7xl mx-auto flex justify-between items-center">
        <h1 className="text-2xl font-bold">Générateur de Quittances</h1>
        <div className="flex items-center gap-4">
          {landlord && (
            <>
              <span className="text-sm">
                {landlord.first_name} {landlord.last_name}
              </span>
              <Button
                variant="secondary"
                onClick={logout}
                className="text-sm"
              >
                Déconnexion
              </Button>
            </>
          )}
        </div>
      </div>
    </header>
  );
};
