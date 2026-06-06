export interface FGStockLot {
  id: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  uom_code: string;
  batch_number: string;
  production_order_id: string | null;
  production_date: string;
  expiry_date: string;
  location_id: number | null;
  location_label: string;
  initial_qty: number;
  current_qty: number;
  status: 'available' | 'reserved' | 'dispatched' | 'expired' | 'quarantine';
  created_at: string;
  updated_at: string;
}

export interface FGStockMovement {
  id: string;
  fg_stock_lot_id: string;
  finished_good_id: string;
  movement_type: 'PRODUCTION_RECEIPT' | 'SALES_DISPATCH' | 'RETURN' | 'ADJUSTMENT_IN' | 'ADJUSTMENT_OUT';
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

export interface FGStockLotDetail extends FGStockLot {
  movements: FGStockMovement[];
}

export interface FGStockSummary {
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  uom_code: string;
  total_qty: number;
  lot_count: number;
  nearest_expiry: string | null;
}

export interface FGAdjustmentLine {
  id: string;
  adjustment_id: string;
  fg_stock_lot_id: string;
  batch_number: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  system_qty: number;
  counted_qty: number;
  variance_qty: number;
  reason: string;
}

export interface FGAdjustment {
  id: string;
  adj_number: string;
  adjustment_date: string;
  type: 'cycle_count' | 'write_off' | 'write_in';
  performed_by: string;
  performed_by_name: string;
  approved_by: string | null;
  approved_by_name: string;
  status: 'pending' | 'approved' | 'cancelled';
  notes: string;
  lines: FGAdjustmentLine[];
  created_at: string;
}
