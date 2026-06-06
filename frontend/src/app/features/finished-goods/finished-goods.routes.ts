import { Routes } from '@angular/router';
import { FinishedGoodsShellComponent } from './finished-goods-shell.component';

export const FINISHED_GOODS_ROUTES: Routes = [
  {
    path: '',
    component: FinishedGoodsShellComponent,
    children: [
      { path: '', redirectTo: 'stock', pathMatch: 'full' },
      {
        path: 'stock',
        loadComponent: () => import('./stock/fg-stock-list.component').then(m => m.FgStockListComponent),
      },
      {
        path: 'adjustments',
        loadComponent: () => import('./adjustments/fg-adjustment-list.component').then(m => m.FgAdjustmentListComponent),
      },
    ],
  },
];
