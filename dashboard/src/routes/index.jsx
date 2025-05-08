import { lazy } from 'react';
import MainLayout from '@/components/layouts/MainLayout.jsx';

const DashboardHome = lazy(() => import('../pages/dashboard.jsx'));
const CloudSetup = lazy(() => import('../pages/cloud-setup.jsx'));
const NotFound = lazy(() => import('../pages/404.jsx'));

export const appRoutes = [
  {
    path: '/',
    element: <MainLayout />,
    children: [
      {
        index: true,
        element: <DashboardHome />,
      },
      {
        path: 'cloud-setup',
        element: <CloudSetup />,
      },
    ],
  },
  {
    path: '*',
    element: <NotFound />,
  },
];
