import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { Landlord, LoginRequest } from '../types/landlord';
import { authService } from '../services/auth';

interface AuthContextType {
  landlord: Landlord | null;
  isAuthenticated: boolean;
  login: (data: LoginRequest) => Promise<void>;
  logout: () => void;
  updateLandlord: (landlord: Landlord) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [landlord, setLandlord] = useState<Landlord | null>(null);

  useEffect(() => {
    // Load landlord from localStorage on mount
    const storedLandlord = authService.getStoredLandlord();
    if (storedLandlord) {
      setLandlord(storedLandlord);
    }
  }, []);

  const login = async (data: LoginRequest) => {
    const response = await authService.login(data);
    setLandlord(response.landlord);
  };

  const logout = () => {
    authService.logout();
    setLandlord(null);
  };

  const updateLandlord = (updatedLandlord: Landlord) => {
    setLandlord(updatedLandlord);
    localStorage.setItem('landlord', JSON.stringify(updatedLandlord));
  };

  return (
    <AuthContext.Provider
      value={{
        landlord,
        isAuthenticated: !!landlord,
        login,
        logout,
        updateLandlord,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
