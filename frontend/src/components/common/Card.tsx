import React, { HTMLAttributes } from 'react';

interface CardProps extends HTMLAttributes<HTMLDivElement> {
  hover?: boolean;
}

export const Card: React.FC<CardProps> = ({
  children,
  hover = false,
  className = '',
  ...props
}) => {
  const hoverClass = hover ? 'card-hover' : 'card';
  
  return (
    <div className={`${hoverClass} ${className}`} {...props}>
      {children}
    </div>
  );
};
