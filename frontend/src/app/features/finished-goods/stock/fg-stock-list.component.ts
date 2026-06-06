import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { DatePipe, DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { FinishedGoodsInventoryService } from '../../../core/services/finished-goods-inventory.service';
import { FGStockSummary, FGStockLot, FGStockLotDetail } from '../../../core/models/finished-goods-inventory.model';

@Component({
  selector: 'app-fg-stock-list',
  standalone: true,
  imports: [
    DatePipe, DecimalPipe,
    TableModule, ButtonModule, DialogModule,
    TagModule, ToastModule,
  ],
  providers: [MessageService],
  templateUrl: './fg-stock-list.component.html',
  styleUrl: './fg-stock-list.component.scss',
})
export class FgStockListComponent implements OnInit {
  private svc = inject(FinishedGoodsInventoryService);
  private msg = inject(MessageService);

  summary = signal<FGStockSummary[]>([]);
  lots = signal<FGStockLot[]>([]);
  selectedSummary: FGStockSummary | null = null;
  selectedLotDetail = signal<FGStockLotDetail | null>(null);
  loadingSummary = signal(false);
  loadingLots = signal(false);
  showLotDialog = false;
  searchText = signal('');

  filteredSummary = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.summary();
    return q ? all.filter(s =>
      s.fg_code.toLowerCase().includes(q) || s.fg_name.toLowerCase().includes(q)
    ) : all;
  });

  selectedFGId: string | null = null;

  ngOnInit() {
    this.loadSummary();
  }

  loadSummary() {
    this.loadingSummary.set(true);
    this.svc.listStockSummary().subscribe({
      next: data => { this.summary.set(data ?? []); this.loadingSummary.set(false); },
      error: () => this.loadingSummary.set(false),
    });
  }

  selectFG(item: FGStockSummary) {
    this.selectedSummary = item;
    this.selectedFGId = item.finished_good_id;
    this.loadingLots.set(true);
    this.svc.listStockLots(item.finished_good_id).subscribe({
      next: data => { this.lots.set(data ?? []); this.loadingLots.set(false); },
      error: () => this.loadingLots.set(false),
    });
  }

  openLotDetail(lot: FGStockLot) {
    this.svc.getStockLotDetail(lot.id).subscribe({
      next: detail => {
        this.selectedLotDetail.set(detail);
        this.showLotDialog = true;
      },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'ไม่พบ lot' }),
    });
  }

  lotStatusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    if (status === 'available') return 'success';
    if (status === 'dispatched') return 'secondary';
    if (status === 'expired') return 'danger';
    if (status === 'quarantine') return 'warn';
    return 'info';
  }

  lotStatusLabel(status: string): string {
    const map: Record<string, string> = {
      available: 'มีสินค้า',
      reserved: 'จอง',
      dispatched: 'จัดส่งแล้ว',
      expired: 'หมดอายุ',
      quarantine: 'กักกัน',
    };
    return map[status] ?? status;
  }

  movementTypeLabel(type: string): string {
    const map: Record<string, string> = {
      PRODUCTION_RECEIPT: 'รับจากผลิต',
      SALES_DISPATCH: 'จัดส่งขาย',
      RETURN: 'รับคืน',
      ADJUSTMENT_IN: 'ปรับเพิ่ม',
      ADJUSTMENT_OUT: 'ปรับลด',
    };
    return map[type] ?? type;
  }

  expiryAlertSeverity(expiry: string | null): 'success' | 'warn' | 'danger' | null {
    if (!expiry) return null;
    const today = new Date();
    const exp = new Date(expiry);
    const diffDays = Math.ceil((exp.getTime() - today.getTime()) / (1000 * 60 * 60 * 24));
    if (diffDays < 0) return 'danger';
    if (diffDays <= 7) return 'warn';
    return null;
  }
}
