export interface SalesOrderLine {
  id: string;
  sales_order_id: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  uom_code: string;
  ordered_qty: number;
  unit_price: number;
  discount_pct: number;
  line_total: number;
  status: 'pending' | 'allocated' | 'dispatched';
}

export interface SalesOrder {
  id: string;
  order_number: string;
  customer_id: string;
  customer_name: string;
  order_date: string;
  requested_date: string | null;
  status: 'draft' | 'confirmed' | 'picking' | 'dispatched' | 'invoiced' | 'cancelled';
  subtotal: number;
  vat_amount: number;
  total_amount: number;
  vat_rate: number;
  notes: string;
  created_by: string;
  created_by_name: string;
  lines: SalesOrderLine[];
  created_at: string;
  updated_at: string;
}

export interface Vehicle {
  id: number;
  license_plate: string;
  type: string;
  capacity_units: number;
  is_active: boolean;
}

export interface DeliveryOrderLine {
  id: string;
  delivery_order_id: string;
  sales_order_line_id: string;
  fg_stock_lot_id: string;
  batch_number: string;
  finished_good_id: string;
  fg_code: string;
  fg_name: string;
  picked_qty: number;
  dispatched_qty: number;
}

export interface DeliveryOrder {
  id: string;
  do_number: string;
  sales_order_id: string;
  so_number: string;
  customer_name: string;
  delivery_date: string | null;
  driver_user_id: string | null;
  driver_name: string;
  vehicle_id: number | null;
  license_plate: string;
  status: 'pending' | 'loading' | 'dispatched' | 'delivered' | 'failed';
  dispatched_at: string | null;
  delivered_at: string | null;
  route_notes: string;
  lines: DeliveryOrderLine[];
  created_at: string;
  updated_at: string;
}

export interface Invoice {
  id: string;
  invoice_number: string;
  sales_order_id: string;
  so_number: string;
  customer_name: string;
  delivery_order_id: string | null;
  do_number: string;
  invoice_date: string;
  due_date: string | null;
  subtotal: number;
  vat_amount: number;
  total_amount: number;
  status: 'issued' | 'paid' | 'overdue' | 'cancelled';
  payment_date: string | null;
  payment_method: string;
  created_by: string;
  created_by_name: string;
  created_at: string;
}

// DTOs
export interface CreateSOLineRequest {
  finished_good_id: string;
  ordered_qty: number;
  unit_price: number;
  discount_pct: number;
}

export interface CreateSalesOrderRequest {
  customer_id: string;
  order_date: string;
  requested_date?: string;
  vat_rate: number;
  notes: string;
  lines: CreateSOLineRequest[];
}

export interface CreateDeliveryOrderRequest {
  sales_order_id: string;
  delivery_date?: string;
  driver_user_id?: string;
  vehicle_id?: number | null;
  route_notes?: string;
}

export interface CreateInvoiceRequest {
  sales_order_id: string;
  delivery_order_id?: string;
  invoice_date: string;
  due_date?: string;
}

export interface MarkPaidRequest {
  payment_date: string;
  payment_method: string;
}
