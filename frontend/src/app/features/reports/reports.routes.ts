import { Routes } from '@angular/router';

export const REPORTS_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () => import('./reports-placeholder.component').then(m => m.ReportsPlaceholderComponent),
  },
];
