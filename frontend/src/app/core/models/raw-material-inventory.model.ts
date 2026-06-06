export interface WarehouseLocation {
  id: number;
  zone: string;
  row_no: string;
  bay_no: string;
  description: string;
}

export interface GRNLine {
  id: string;
  grn_id: string;
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  lot_number: string;
  received_qty: number;
  unit_cost: number | null;
  expiry_date: string | null;
  location_id: number | null;
  location_label: string;
  notes: string;
}

export interface GRN {
  id: string;
  grn_number: string;
  supplier_id: string | null;
  supplier_name: string;
  received_by: string;
  received_by_name: string;
  received_date: string;
  po_reference: string;
  status: 'draft' | 'confirmed' | 'cancelled';
  notes: string;
  lines: GRNLine[];
  created_at: string;
  updated_at: string;
}

export interface StockLot {
  id: string;
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  uom_code: string;
  lot_number: string;
  grn_line_id: string | null;
  received_date: string;
  expiry_date: string | null;
  location_id: number | null;
  location_label: string;
  initial_qty: number;
  current_qty: number;
  unit_cost: number | null;
  status: 'available' | 'reserved' | 'depleted' | 'quarantine';
  created_at: string;
  updated_at: string;
}

export interface StockMovement {
  id: string;
  rm_stock_lot_id: string;
  raw_material_id: string;
  movement_type: string;
  reference_type: string;
  reference_id: string | null;
  qty: number;
  qty_before: number;
  qty_after: number;
  performed_by: string;
  performed_by_name: string;
  notes: string;
  created_at: string;
}

export interface StockLotDetail extends StockLot {
  movements: StockMovement[];
}

export interface StockSummary {
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  uom_code: string;
  total_qty: number;
  lot_count: number;
  earliest_expiry: string | null;
}

export interface AdjustmentLine {
  id: string;
  adjustment_id: string;
  rm_stock_lot_id: string;
  lot_number: string;
  raw_material_id: string;
  rm_code: string;
  rm_name: string;
  system_qty: number;
  counted_qty: number;
  variance_qty: number;
  reason: string;
}

export interface Adjustment {
  id: string;
  adj_number: string;
  adjustment_date: string;
  type: 'cycle_count' | 'damage' | 'correction';
  performed_by: string;
  performed_by_name: string;
  approved_by: string | null;
  approved_by_name: string;
  status: 'pending' | 'approved' | 'cancelled';
  notes: string;
  lines: AdjustmentLine[];
  created_at: string;
}
