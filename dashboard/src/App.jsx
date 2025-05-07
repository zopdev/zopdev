import { HashRouter, Route, Routes } from 'react-router-dom';
import { ErrorBoundary } from './components/organisms/ErrorBoundary';
import { Suspense } from 'react';
import { routes } from './routes';
import CompleteLoader from '@/components/atom/Loaders/CompleteLoader.jsx';
import TopBar from '@/components/molecules/TopBar/index.jsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import ToastContainer from '@/components/organisms/ToastContainer.jsx';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <HashRouter>
        <ErrorBoundary>
          <Suspense fallback={<CompleteLoader />}>
            <ToastContainer classNameParent={`absolute right-5 top-10 `} stacked={false} />
            <div className="flex flex-col h-screen">
              <div className="sticky top-0 z-50 bg-white shadow">
                <TopBar />
              </div>

              <div className="flex-1 overflow-y-auto">
                <Routes>
                  {routes.map(({ path, component: Component }) => (
                    <Route key={path} path={path} element={<Component />} />
                  ))}
                </Routes>
              </div>
            </div>
          </Suspense>
        </ErrorBoundary>
      </HashRouter>
    </QueryClientProvider>
  );
}

export default App;
