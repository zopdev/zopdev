import { HashRouter, Route, Routes } from 'react-router-dom';
import { ErrorBoundary } from './components/organisms/ErrorBoundary';
import { Suspense } from 'react';
import { routes } from './routes';

function App() {
  return (
    <HashRouter>
      <ErrorBoundary>
        <Suspense fallback={<div>Loading...</div>}>
          <Routes>
            {routes.map(({ path, component: Component }) => (
              <Route key={path} path={path} element={<Component />} />
            ))}
          </Routes>
        </Suspense>
      </ErrorBoundary>
    </HashRouter>
  );
}

export default App;
