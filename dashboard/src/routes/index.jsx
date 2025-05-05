import { lazy } from 'react';

export const routes = [
  { path: '/', component: lazy(() => import('../pages/home')) },
  { path: '/cloud-setup', component: lazy(() => import('../pages/cloud-setup.jsx')) },
  { path: '/about', component: lazy(() => import('../pages/about')) },
  { path: '*', component: lazy(() => import('../pages/404')) },
];
