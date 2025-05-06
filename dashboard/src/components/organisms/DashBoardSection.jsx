'use client';

import React from 'react';

const DashboardSection = ({ title, children }) => {
  return (
    <div className="w-full md:w-1/2 flex flex-col gap-6">
      <h2 className="text-xl font-semibold">{title}</h2>
      {children}
    </div>
  );
};

export default DashboardSection;
