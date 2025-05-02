'use client';

import React from 'react';
import Sidebar from '../../../../partials/sideBar';

const Layout = ({ children }) => {
  return (
    <div className=" xs:inline-block md:flex h-[90vh] w-[100vw] overflow-hidden ">
      <div>
        <Sidebar />
      </div>
      <main className="px-4 sm:px-6 lg:px-8 w-full overflow-auto text-left pt-8">{children}</main>
    </div>
  );
};

export default Layout;
