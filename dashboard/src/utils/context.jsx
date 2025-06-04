import { createContext, useState } from 'react';

function ContextProvider({ children }) {
  const [appData, setAppData] = useState({
    CLOUD_ACCOUNT_DATA: { data: [], isLoaded: false, isError: false },
    APPLICATION_DATA: { data: [], isLoaded: false, isError: false },
  });

  return <AppContext.Provider value={{ appData, setAppData }}>{children}</AppContext.Provider>;
}

export default ContextProvider;

export const AppContext = createContext();
