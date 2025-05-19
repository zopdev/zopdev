import { lazy } from 'react';
import MainLayout from '@/components/layouts/MainLayout.jsx';

const NotFoundPage = lazy(() => import('../pages/404.jsx'));
const DashboardPage = lazy(() => import('../pages/index.jsx'));
const CloudSetupPage = lazy(() => import('../pages/cloud-setup.jsx'));
const CloudAccountsPage = lazy(() => import('../pages/cloud-accounts.jsx'));

export const appRoutes = [
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <DashboardPage />,
      },
      {
        path: 'cloud-setup',
        element: <CloudSetupPage />,
      },
    ],
  },
  {
    path: '/cloud-accounts',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <CloudAccountsPage />,
      },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
];
