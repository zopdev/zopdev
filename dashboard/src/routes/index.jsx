import { lazy } from 'react';
import MainLayout from '@/components/layouts/MainLayout.jsx';

const DashboardPage = lazy(() => import('../pages/index.jsx'));
const CloudSetupPage = lazy(() => import('../pages/cloud-setup.jsx'));
const NotFoundPage = lazy(() => import('../pages/404.jsx'));

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
    path: '*',
    element: <NotFoundPage />,
  },
];
