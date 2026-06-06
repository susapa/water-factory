export interface ProductionRequirement {
  id: string;
  production_order_id: string;
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  uom_code: string;
  required_qty: number;
  issued_qty: number;
  status: 'pending' | 'partial' | 'fully_issued';
}

export interface ProductionOrder {
  id: string;
  order_number: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  bom_id: string | null;
  planned_qty: number;
  actual_yield_qty: number;
  defect_qty: number;
  waste_qty: number;
  planned_start_date: string | null;
  planned_end_date: string | null;
  actual_start_date: string | null;
  actual_end_date: string | null;
  status: 'draft' | 'confirmed' | 'rm_issued' | 'in_progress' | 'completed' | 'cancelled';
  created_by: string;
  created_by_name: string;
  sales_order_id: string | null;
  notes: string;
  requirements: ProductionRequirement[];
  created_at: string;
  updated_at: string;
}

export interface DefectDetail {
  id: string;
  yield_id: string;
  defect_type: string;
  qty: number;
  description: string;
}

export interface ProductionYield {
  id: string;
  production_order_id: string;
  order_number: string;
  fg_name: string;
  yield_qty: number;
  defect_qty: number;
  waste_qty: number;
  batch_number: string;
  production_date: string;
  expiry_date: string;
  recorded_by: string;
  recorded_by_name: string;
  recorded_at: string;
  notes: string;
  defect_details: DefectDetail[];
}
