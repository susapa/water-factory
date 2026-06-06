import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { ProductionService } from '../../../core/services/production.service';
import { ProductionYield } from '../../../core/models/production.model';

@Component({
  selector: 'app-yield-list',
  standalone: true,
  imports: [DecimalPipe, TableModule, TagModule, ToastModule],
  providers: [MessageService],
  templateUrl: './yield-list.component.html',
  styleUrl: './yield-list.component.scss',
})
export class YieldListComponent implements OnInit {
  private svc = inject(ProductionService);

  items = signal<ProductionYield[]>([]);
  loading = signal(false);
  searchText = signal('');
  expandedRows: Record<string, boolean> = {};

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(y =>
      y.batch_number.toLowerCase().includes(q) ||
      y.order_number.toLowerCase().includes(q) ||
      y.fg_name.toLowerCase().includes(q)
    ) : all;
  });

  ngOnInit() {
    this.loading.set(true);
    this.svc.listYields().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }
}
