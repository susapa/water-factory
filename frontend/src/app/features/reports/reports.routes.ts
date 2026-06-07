import { Routes } from '@angular/router';
import { ReportsShellComponent } from './reports-shell.component';

export const REPORTS_ROUTES: Routes = [
  {
    path: '',
    component: ReportsShellComponent,
    children: [
      { path: '', redirectTo: 'stock', pathMatch: 'full' },
      {
        path: 'stock',
        loadComponent: () =>
          import('./stock/stock-report.component').then(m => m.StockReportComponent),
      },
      {
        path: 'sales',
        loadComponent: () =>
          import('./sales/sales-report.component').then(m => m.SalesReportComponent),
      },
      {
        path: 'production',
        loadComponent: () =>
          import('./production/production-report.component').then(m => m.ProductionReportComponent),
      },
    ],
  },
];
