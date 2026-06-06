import { Routes } from '@angular/router';

export const SALES_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () => import('./sales-placeholder.component').then(m => m.SalesPlaceholderComponent),
  },
];
