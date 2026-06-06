import { Routes } from '@angular/router';
import { authGuard } from './core/auth/auth.guard';
import { roleGuard } from './core/auth/role.guard';
import { LayoutComponent } from './layout/layout.component';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/login/login.component').then(m => m.LoginComponent),
  },
  {
    path: '',
    component: LayoutComponent,
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./features/dashboard/dashboard.component').then(m => m.DashboardComponent),
        canActivate: [roleGuard],
        data: { roles: ['admin','warehouse_manager','production_manager','sales_admin'] },
      },
      {
        path: 'master-data',
        loadChildren: () =>
          import('./features/master-data/master-data.routes').then(m => m.MASTER_DATA_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin'] },
      },
      {
        path: 'raw-material',
        loadChildren: () =>
          import('./features/raw-material/raw-material.routes').then(m => m.RAW_MATERIAL_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin','warehouse_manager','production_manager'] },
      },
      {
        path: 'production',
        loadChildren: () =>
          import('./features/production/production.routes').then(m => m.PRODUCTION_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin','warehouse_manager','production_manager'] },
      },
      {
        path: 'finished-goods',
        loadChildren: () =>
          import('./features/finished-goods/finished-goods.routes').then(m => m.FINISHED_GOODS_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin','warehouse_manager','sales_admin'] },
      },
      {
        path: 'sales',
        loadChildren: () =>
          import('./features/sales/sales.routes').then(m => m.SALES_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin','sales_admin','delivery'] },
      },
      {
        path: 'reports',
        loadChildren: () =>
          import('./features/reports/reports.routes').then(m => m.REPORTS_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin','warehouse_manager','production_manager','sales_admin'] },
      },
      {
        path: 'admin',
        loadChildren: () =>
          import('./features/admin/admin.routes').then(m => m.ADMIN_ROUTES),
        canActivate: [roleGuard],
        data: { roles: ['admin'] },
      },
    ],
  },
  { path: '**', redirectTo: '' },
];
