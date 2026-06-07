import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { SalesService } from '../../../core/services/sales.service';
import { SalesOrder } from '../../../core/models/sales.model';

@Component({
  selector: 'app-sales-report',
  standalone: true,
  imports: [TableModule, TagModule, DecimalPipe],
  templateUrl: './sales-report.component.html',
  styleUrl: './sales-report.component.scss',
})
export class SalesReportComponent implements OnInit {
  private svc = inject(SalesService);

  allOrders = signal<SalesOrder[]>([]);
  loading = signal(false);
  dateFrom = signal('');
  dateTo = signal('');

  filtered = computed(() => {
    const from = this.dateFrom();
    const to = this.dateTo();
    return this.allOrders().filter(o => {
      if (from && o.order_date < from) return false;
      if (to && o.order_date > to) return false;
      return true;
    });
  });

  total = computed(() => this.filtered().reduce((s, o) => s + o.total_amount, 0));

  statusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    const map: Record<string, 'success' | 'secondary' | 'info' | 'warn' | 'danger'> = {
      draft: 'secondary',
      confirmed: 'info',
      picking: 'warn',
      dispatched: 'warn',
      invoiced: 'success',
      cancelled: 'danger',
    };
    return map[status] ?? 'secondary';
  }

  ngOnInit() {
    this.loading.set(true);
    this.svc.listSalesOrders().subscribe({
      next: d => { this.allOrders.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }
}
