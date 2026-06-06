export interface UOM {
  id: number;
  code: string;
  name: string;
}

export interface RawMaterialCategory {
  id: number;
  name: string;
}

export interface RawMaterial {
  id: string;
  code: string;
  name: string;
  category_id: number | null;
  category_name: string;
  uom_id: number;
  uom_code: string;
  min_stock_qty: number;
  reorder_qty: number;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface FinishedGood {
  id: string;
  code: string;
  name: string;
  uom_id: number;
  uom_code: string;
  shelf_life_days: number;
  min_stock_qty: number;
  description: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface BOMLine {
  id: string;
  bom_id: string;
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  qty_per_unit: number;
  waste_factor: number;
}

export interface BOM {
  id: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  version: number;
  is_active: boolean;
  effective_date: string;
  notes: string;
  lines: BOMLine[];
  created_at: string;
}

export interface Customer {
  id: string;
  code: string;
  name: string;
  tax_id: string;
  address: string;
  phone: string;
  email: string;
  credit_limit: number;
  credit_days: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Supplier {
  id: string;
  code: string;
  name: string;
  tax_id: string;
  address: string;
  phone: string;
  email: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}
