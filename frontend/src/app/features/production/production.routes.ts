import { Routes } from '@angular/router';
import { ProductionShellComponent } from './production-shell.component';

export const PRODUCTION_ROUTES: Routes = [
  {
    path: '',
    component: ProductionShellComponent,
    children: [
      { path: '', redirectTo: 'orders', pathMatch: 'full' },
      {
        path: 'orders',
        loadComponent: () => import('./orders/order-list.component').then(m => m.OrderListComponent),
      },
      {
        path: 'yields',
        loadComponent: () => import('./yields/yield-list.component').then(m => m.YieldListComponent),
      },
    ],
  },
];
