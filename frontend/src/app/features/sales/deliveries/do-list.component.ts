import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { CommonModule, DecimalPipe } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MessageService, ConfirmationService } from 'primeng/api';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { SelectModule } from 'primeng/select';
import { InputTextModule } from 'primeng/inputtext';

import { SalesService } from '../../../core/services/sales.service';
import { DeliveryOrder, SalesOrder, Vehicle } from '../../../core/models/sales.model';

@Component({
  selector: 'app-do-list',
  standalone: true,
  imports: [
    CommonModule, ReactiveFormsModule, DecimalPipe,
    TableModule, ButtonModule, TagModule, DialogModule,
    ToastModule, ConfirmDialogModule, SelectModule, InputTextModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './do-list.component.html',
  styleUrl: './do-list.component.scss',
})
export class DoListComponent implements OnInit {
  private svc = inject(SalesService);
  private fb = inject(FormBuilder);
  private toast = inject(MessageService);
  private confirm = inject(ConfirmationService);

  items = signal<DeliveryOrder[]>([]);
  loading = signal(false);
  searchText = signal('');
  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    return this.items().filter(d =>
      d.do_number.toLowerCase().includes(q) ||
      d.so_number.toLowerCase().includes(q) ||
      d.customer_name.toLowerCase().includes(q) ||
      d.status.toLowerCase().includes(q)
    );
  });

  confirmedSOs = signal<SalesOrder[]>([]);
  vehicles = signal<Vehicle[]>([]);

  showCreateDialog = false;
  showDetailDialog = false;
  detailDO = signal<DeliveryOrder | null>(null);

  form = this.fb.group({
    sales_order_id: ['', Validators.required],
    delivery_date:  [''],
    vehicle_id:     [null as number | null],
    route_notes:    [''],
  });

  ngOnInit() {
    this.load();
    this.loadSOOptions();
    this.svc.listVehicles().subscribe(v => this.vehicles.set(v));
  }

  load() {
    this.loading.set(true);
    this.svc.listDeliveryOrders().subscribe({
      next: d => { this.items.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  loadSOOptions() {
    this.svc.listSalesOrders().subscribe(orders =>
      this.confirmedSOs.set(orders.filter(o => o.status === 'confirmed'))
    );
  }

  soOptions = computed(() =>
    this.confirmedSOs().map(o => ({ id: o.id, label: `${o.order_number} — ${o.customer_name}` }))
  );

  vehicleOptions = computed(() =>
    this.vehicles().map(v => ({ id: v.id, label: `${v.license_plate} (${v.type})` }))
  );

  openCreate() {
    this.form.reset({ sales_order_id: '', delivery_date: '', vehicle_id: null, route_notes: '' });
    this.loadSOOptions();
    this.showCreateDialog = true;
  }

  openDetail(d: DeliveryOrder) {
    this.detailDO.set(d);
    this.showDetailDialog = true;
  }

  save() {
    this.form.markAllAsTouched();
    if (this.form.invalid) return;
    const v = this.form.value;
    this.svc.createDeliveryOrder({
      sales_order_id: v.sales_order_id!,
      delivery_date: v.delivery_date || undefined,
      vehicle_id: v.vehicle_id ?? undefined,
      route_notes: v.route_notes || '',
    }).subscribe({
      next: () => {
        this.toast.add({ severity: 'success', summary: 'สำเร็จ', detail: 'สร้างใบส่งสินค้าแล้ว' });
        this.showCreateDialog = false;
        this.load();
      },
      error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
    });
  }

  dispatch(d: DeliveryOrder) {
    this.confirm.confirm({
      message: `จัดส่ง ${d.do_number} และหักสต็อกสินค้า?`,
      accept: () => this.svc.dispatchDeliveryOrder(d.id).subscribe({
        next: () => { this.toast.add({ severity: 'success', summary: 'จัดส่งแล้ว' }); this.load(); },
        error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
      }),
    });
  }

  deliver(d: DeliveryOrder) {
    this.confirm.confirm({
      message: `ยืนยันการส่งมอบ ${d.do_number}?`,
      accept: () => this.svc.deliverDeliveryOrder(d.id).subscribe({
        next: () => { this.toast.add({ severity: 'success', summary: 'ส่งมอบแล้ว' }); this.load(); },
        error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
      }),
    });
  }

  statusLabel(s: string): string {
    return { pending: 'รอจัดส่ง', loading: 'กำลังโหลด', dispatched: 'จัดส่งแล้ว', delivered: 'ส่งมอบแล้ว', failed: 'ล้มเหลว' }[s] ?? s;
  }
  statusSeverity(s: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    return ({ pending: 'info', loading: 'warn', dispatched: 'success', delivered: 'success', failed: 'danger' } as any)[s] ?? 'secondary';
  }
}
