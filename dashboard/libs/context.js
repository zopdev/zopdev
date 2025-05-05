'use client';

import { createContext, useState } from 'react';

function ContextProvider({ children }) {
  const [appData, setAppData] = useState({
    CLOUD_ACCOUNT_DATA: { data: [], isSuccess: false },
    // CLUSTERS_DATA: { data: [], isSuccess: false, isError: false },
    // NAMESPACE_DATA: { data: [], isSuccess: false },
    APPLICATION_DATA: { data: [], isSuccess: false },
  });

  return <AppContext.Provider value={{ appData, setAppData }}>{children}</AppContext.Provider>;
}

export default ContextProvider;

export const AppContext = createContext();
