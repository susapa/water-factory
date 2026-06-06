export interface User {
  id: string;
  email: string;
  full_name: string;
  role: UserRole;
  is_active: boolean;
}

export type UserRole =
  | 'admin'
  | 'warehouse_manager'
  | 'production_manager'
  | 'sales_admin'
  | 'delivery';

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface Role {
  id: number;
  name: string;
  description: string;
}
