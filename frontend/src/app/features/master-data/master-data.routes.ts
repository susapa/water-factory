import { Routes } from '@angular/router';

export const MASTER_DATA_ROUTES: Routes = [
  {
    path: '',
    loadComponent: () =>
      import('./master-data-shell.component').then(m => m.MasterDataShellComponent),
    children: [
      { path: '', redirectTo: 'raw-materials', pathMatch: 'full' },
      {
        path: 'raw-materials',
        loadComponent: () =>
          import('./raw-materials/raw-material-list.component').then(m => m.RawMaterialListComponent),
      },
      {
        path: 'finished-goods',
        loadComponent: () =>
          import('./finished-goods/finished-good-list.component').then(m => m.FinishedGoodListComponent),
      },
      {
        path: 'bom',
        loadComponent: () =>
          import('./bom/bom-list.component').then(m => m.BomListComponent),
      },
      {
        path: 'customers',
        loadComponent: () =>
          import('./customers/customer-list.component').then(m => m.CustomerListComponent),
      },
      {
        path: 'suppliers',
        loadComponent: () =>
          import('./suppliers/supplier-list.component').then(m => m.SupplierListComponent),
      },
    ],
  },
];
