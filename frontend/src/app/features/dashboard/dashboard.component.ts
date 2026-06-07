import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { ChartModule } from 'primeng/chart';
import { CardModule } from 'primeng/card';
import { TableModule } from 'primeng/table';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { DashboardService } from '../../core/services/dashboard.service';
import { DashboardSummary } from '../../core/models/dashboard.model';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [ChartModule, CardModule, TableModule, TagModule, ToastModule, DecimalPipe],
  providers: [MessageService],
  templateUrl: './dashboard.component.html',
  styleUrl: './dashboard.component.scss',
})
export class DashboardComponent implements OnInit {
  private svc = inject(DashboardService);
  private toast = inject(MessageService);

  summary = signal<DashboardSummary | null>(null);
  loading = signal(false);

  barChartData = computed(() => {
    const days = this.summary()?.sales_last_7_days ?? [];
    return {
      labels: days.map(d => d.date),
      datasets: [{
        label: 'ยอดขาย (บาท)',
        data: days.map(d => d.total),
        backgroundColor: '#6366f1',
        borderRadius: 4,
      }],
    };
  });

  donutChartData = computed(() => {
    const entries = Object.entries(this.summary()?.po_status_counts ?? {});
    return {
      labels: entries.map(([k]) => k),
      datasets: [{
        data: entries.map(([, v]) => v),
        backgroundColor: ['#6366f1', '#22c55e', '#f59e0b', '#ef4444', '#3b82f6', '#a855f7'],
      }],
    };
  });

  barChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { display: false } },
    scales: { y: { beginAtZero: true, ticks: { precision: 0 } } },
  };

  hasPOCounts = computed(() => Object.keys(this.summary()?.po_status_counts ?? {}).length > 0);

  donutChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: { legend: { position: 'bottom' } },
  };

  ngOnInit() {
    this.loading.set(true);
    this.svc.getDashboardSummary().subscribe({
      next: d => { this.summary.set(d); this.loading.set(false); },
      error: () => {
        this.loading.set(false);
        this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: 'โหลดข้อมูล dashboard ไม่ได้' });
      },
    });
  }
}
