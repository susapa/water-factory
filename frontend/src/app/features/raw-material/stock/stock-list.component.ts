import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { DatePipe, DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { RawMaterialInventoryService } from '../../../core/services/raw-material-inventory.service';
import { StockSummary, StockLot, StockLotDetail } from '../../../core/models/raw-material-inventory.model';

@Component({
  selector: 'app-stock-list',
  standalone: true,
  imports: [
    DatePipe, DecimalPipe,
    TableModule, ButtonModule, DialogModule,
    TagModule, ToastModule,
  ],
  providers: [MessageService],
  templateUrl: './stock-list.component.html',
  styleUrl: './stock-list.component.scss',
})
export class StockListComponent implements OnInit {
  private svc = inject(RawMaterialInventoryService);
  private msg = inject(MessageService);

  summary = signal<StockSummary[]>([]);
  lots = signal<StockLot[]>([]);
  selectedSummary: StockSummary | null = null;
  selectedLotDetail = signal<StockLotDetail | null>(null);
  loadingSummary = signal(false);
  loadingLots = signal(false);
  showLotDialog = false;
  searchText = signal('');

  filteredSummary = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.summary();
    return q ? all.filter(s =>
      s.rm_code.toLowerCase().includes(q) || s.rm_name.toLowerCase().includes(q)
    ) : all;
  });

  selectedRMId: string | null = null;

  lotStatusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    if (status === 'available') return 'success';
    if (status === 'depleted') return 'secondary';
    if (status === 'quarantine') return 'warn';
    return 'info';
  }
  lotStatusLabel(status: string): string {
    if (status === 'available') return 'มีสินค้า';
    if (status === 'depleted') return 'หมด';
    if (status === 'quarantine') return 'กักกัน';
    return 'จอง';
  }
  movementTypeLabel(type: string): string {
    const map: Record<string, string> = {
      GRN: 'รับวัตถุดิบ',
      ISSUE_TO_PROD: 'เบิกผลิต',
      RETURN_FROM_PROD: 'คืนจากผลิต',
      ADJUSTMENT_IN: 'ปรับเพิ่ม',
      ADJUSTMENT_OUT: 'ปรับลด',
      CYCLE_COUNT: 'นับสต็อก',
    };
    return map[type] ?? type;
  }

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

  selectRM(item: StockSummary) {
    this.selectedSummary = item;
    this.selectedRMId = item.raw_material_id;
    this.loadingLots.set(true);
    this.svc.listStockLots(item.raw_material_id).subscribe({
      next: data => { this.lots.set(data ?? []); this.loadingLots.set(false); },
      error: () => this.loadingLots.set(false),
    });
  }

  openLotDetail(lot: StockLot) {
    this.svc.getStockLotDetail(lot.id).subscribe({
      next: detail => {
        this.selectedLotDetail.set(detail);
        this.showLotDialog = true;
      },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'ไม่พบ lot' }),
    });
  }
}
