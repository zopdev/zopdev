import { lazy } from 'react';

export const routes = [
  { path: '/', component: lazy(() => import('../pages/home')) },
  { path: '/about', component: lazy(() => import('../pages/about')) },
  { path: '*', component: lazy(() => import('../pages/404')) },
];
