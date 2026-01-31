import React from 'react';
import { NavLink } from 'react-router-dom';

interface MenuItem {
  path: string;
  label: string;
  icon: string;
}

const menuItems: MenuItem[] = [
  { path: '/dashboard', label: 'Tableau de bord', icon: '📊' },
  { path: '/landlord', label: 'Mes informations', icon: '👤' },
  { path: '/properties', label: 'Propriétés', icon: '🏠' },
  { path: '/tenants', label: 'Locataires', icon: '👥' },
  { path: '/receipts', label: 'Quittances', icon: '📄' },
  { path: '/history', label: 'Historique', icon: '📋' },
];

export const Sidebar: React.FC = () => {
  return (
    <aside className="w-64 bg-gray-50 border-r-2 border-gray-200 min-h-screen">
      <nav className="py-4">
        {menuItems.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) =>
              `flex items-center gap-3 px-6 py-3 transition-all duration-200 border-l-4 ${
                isActive
                  ? 'bg-gradient-to-r from-primary-50 to-white border-primary-500 text-primary-600 font-semibold'
                  : 'border-transparent text-gray-700 hover:bg-gray-100 hover:text-primary-500'
              }`
            }
          >
            <span className="text-xl">{item.icon}</span>
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>
    </aside>
  );
};
