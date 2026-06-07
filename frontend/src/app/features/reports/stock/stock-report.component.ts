import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { RawMaterialInventoryService } from '../../../core/services/raw-material-inventory.service';
import { FinishedGoodsInventoryService } from '../../../core/services/finished-goods-inventory.service';
import { StockSummary } from '../../../core/models/raw-material-inventory.model';
import { FGStockSummary } from '../../../core/models/finished-goods-inventory.model';

@Component({
  selector: 'app-stock-report',
  standalone: true,
  imports: [TableModule, TagModule, DecimalPipe],
  templateUrl: './stock-report.component.html',
  styleUrl: './stock-report.component.scss',
})
export class StockReportComponent implements OnInit {
  private rmSvc = inject(RawMaterialInventoryService);
  private fgSvc = inject(FinishedGoodsInventoryService);

  rmSummary = signal<StockSummary[]>([]);
  fgSummary = signal<FGStockSummary[]>([]);
  loading = signal(false);

  expiryTagSeverity(dateStr: string | null): 'danger' | 'warn' | 'success' {
    if (!dateStr) return 'success';
    const days = (new Date(dateStr).getTime() - Date.now()) / 86400000;
    if (days < 0) return 'danger';
    if (days <= 7) return 'warn';
    return 'success';
  }

  ngOnInit() {
    this.loading.set(true);
    this.rmSvc.listStockSummary().subscribe(d => this.rmSummary.set(d));
    this.fgSvc.listStockSummary().subscribe({
      next: d => { this.fgSummary.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }
}
