import { Routes } from '@angular/router';
import { RawMaterialShellComponent } from './raw-material-shell.component';

export const RAW_MATERIAL_ROUTES: Routes = [
  {
    path: '',
    component: RawMaterialShellComponent,
    children: [
      { path: '', redirectTo: 'grn', pathMatch: 'full' },
      {
        path: 'grn',
        loadComponent: () => import('./grn/grn-list.component').then(m => m.GrnListComponent),
      },
      {
        path: 'stock',
        loadComponent: () => import('./stock/stock-list.component').then(m => m.StockListComponent),
      },
      {
        path: 'adjustments',
        loadComponent: () => import('./adjustments/adjustment-list.component').then(m => m.AdjustmentListComponent),
      },
    ],
  },
];
