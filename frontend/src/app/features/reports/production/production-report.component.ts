import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ProductionService } from '../../../core/services/production.service';
import { ProductionOrder } from '../../../core/models/production.model';

@Component({
  selector: 'app-production-report',
  standalone: true,
  imports: [TableModule, TagModule, DecimalPipe],
  templateUrl: './production-report.component.html',
  styleUrl: './production-report.component.scss',
})
export class ProductionReportComponent implements OnInit {
  private svc = inject(ProductionService);

  allOrders = signal<ProductionOrder[]>([]);
  loading = signal(false);
  dateFrom = signal('');
  dateTo = signal('');

  filtered = computed(() => {
    const from = this.dateFrom();
    const to = this.dateTo();
    return this.allOrders().filter(o => {
      const d = o.planned_start_date ?? '';
      if (from && d < from) return false;
      if (to && d > to) return false;
      return true;
    });
  });

  totalPlanned = computed(() => this.filtered().reduce((s, o) => s + o.planned_qty, 0));
  totalYield = computed(() => this.filtered().reduce((s, o) => s + o.actual_yield_qty, 0));
  totalDefect = computed(() => this.filtered().reduce((s, o) => s + o.defect_qty, 0));

  statusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    const map: Record<string, 'success' | 'secondary' | 'info' | 'warn' | 'danger'> = {
      draft: 'secondary',
      confirmed: 'info',
      rm_issued: 'info',
      in_progress: 'warn',
      completed: 'success',
      cancelled: 'danger',
    };
    return map[status] ?? 'secondary';
  }

  ngOnInit() {
    this.loading.set(true);
    this.svc.listOrders().subscribe({
      next: d => { this.allOrders.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }
}
