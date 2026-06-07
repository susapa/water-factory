export interface DailySales {
  date: string;
  total: number;
}

export interface FGExpiringSoonItem {
  lot_number: string;
  fg_name: string;
  qty: number;
  expiry_date: string;
}

export interface DashboardSummary {
  pending_so_count: number;
  in_progress_po_count: number;
  pending_do_count: number;
  issued_invoice_count: number;
  issued_invoice_total: number;
  sales_last_7_days: DailySales[];
  po_status_counts: Record<string, number>;
  fg_expiring_soon: FGExpiringSoonItem[];
}
