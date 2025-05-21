import { HashRouter, useRoutes } from 'react-router-dom';
import { ErrorBoundary } from './components/organisms/ErrorBoundary';
import { Suspense } from 'react';
import { appRoutes } from './routes';
import CompleteLoader from '@/components/atom/Loaders/CompleteLoader.jsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ToastContainer from '@/components/organisms/ToastContainer.jsx';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 0,
      cacheTime: 0,
      refetchOnWindowFocus: false,
      refetchOnReconnect: false,
      retry: false,
    },
  },
});

function AppRoutes() {
  return useRoutes(appRoutes);
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <HashRouter>
        <ErrorBoundary>
          <Suspense fallback={<CompleteLoader />}>
            <ToastContainer classNameParent="absolute right-5 top-10" stacked={false} />
            <AppRoutes />
          </Suspense>
        </ErrorBoundary>
      </HashRouter>
    </QueryClientProvider>
  );
}

export default App;
