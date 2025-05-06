import { lazy } from 'react';

export const routes = [
  { path: '/', component: lazy(() => import('../pages/dashboard.jsx')) },
  { path: '/cloud-setup', component: lazy(() => import('../pages/cloud-setup.jsx')) },
  { path: '*', component: lazy(() => import('../pages/404')) },
];
