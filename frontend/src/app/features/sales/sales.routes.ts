import { Routes } from '@angular/router';
import { SalesShellComponent } from './sales-shell.component';

export const SALES_ROUTES: Routes = [
  {
    path: '',
    component: SalesShellComponent,
    children: [
      { path: '', redirectTo: 'orders', pathMatch: 'full' },
      {
        path: 'orders',
        loadComponent: () => import('./orders/so-list.component').then(m => m.SoListComponent),
      },
      {
        path: 'deliveries',
        loadComponent: () => import('./deliveries/do-list.component').then(m => m.DoListComponent),
      },
      {
        path: 'invoices',
        loadComponent: () => import('./invoices/invoice-list.component').then(m => m.InvoiceListComponent),
      },
    ],
  },
];
