import { Component, OnInit, signal, computed, inject } from '@angular/core';
import { CommonModule, DecimalPipe } from '@angular/common';
import { FormBuilder, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { MessageService, ConfirmationService } from 'primeng/api';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { TagModule } from 'primeng/tag';
import { DialogModule } from 'primeng/dialog';
import { ToastModule } from 'primeng/toast';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { SelectModule } from 'primeng/select';
import { InputTextModule } from 'primeng/inputtext';
import { InputNumberModule } from 'primeng/inputnumber';

import { SalesService } from '../../../core/services/sales.service';
import { MasterDataService } from '../../../core/services/master-data.service';
import { SalesOrder, SalesOrderLine, CreateSOLineRequest } from '../../../core/models/sales.model';
import { FinishedGood } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-so-list',
  standalone: true,
  imports: [
    CommonModule, ReactiveFormsModule, DecimalPipe,
    TableModule, ButtonModule, TagModule, DialogModule,
    ToastModule, ConfirmDialogModule, SelectModule, InputTextModule, InputNumberModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './so-list.component.html',
  styleUrl: './so-list.component.scss',
})
export class SoListComponent implements OnInit {
  private svc = inject(SalesService);
  private mdSvc = inject(MasterDataService);
  private fb = inject(FormBuilder);
  private toast = inject(MessageService);
  private confirm = inject(ConfirmationService);

  items = signal<SalesOrder[]>([]);
  loading = signal(false);
  searchText = signal('');
  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    return this.items().filter(o =>
      o.order_number.toLowerCase().includes(q) ||
      o.customer_name.toLowerCase().includes(q) ||
      o.status.toLowerCase().includes(q)
    );
  });

  fgOptions = signal<FinishedGood[]>([]);
  showDialog = false;
  viewOnly = false;
  dialogTitle = '';
  editingOrder = signal<SalesOrder | null>(null);

  form = this.fb.group({
    customer_id:    ['', Validators.required],
    order_date:     [new Date().toISOString().slice(0, 10), Validators.required],
    requested_date: [''],
    vat_rate:       [7],
    notes:          [''],
    lines:          this.fb.array([]),
  });

  customers = signal<{ id: string; name: string }[]>([]);

  get lines(): FormArray { return this.form.get('lines') as FormArray; }

  ngOnInit() {
    this.load();
    this.mdSvc.listFinishedGoods().subscribe(g => this.fgOptions.set(g.filter(x => x.is_active)));
    this.mdSvc.listCustomers().subscribe(c => this.customers.set(c.filter(x => x.is_active).map(x => ({ id: x.id, name: x.name }))));
  }

  load() {
    this.loading.set(true);
    this.svc.listSalesOrders().subscribe({
      next: d => { this.items.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  openCreate() {
    this.viewOnly = false;
    this.editingOrder.set(null);
    this.dialogTitle = 'สร้างใบสั่งขาย';
    this.form.reset({
      customer_id: '', order_date: new Date().toISOString().slice(0, 10),
      requested_date: '', vat_rate: 7, notes: '',
    });
    this.lines.clear();
    this.addLine();
    this.showDialog = true;
  }

  openView(order: SalesOrder) {
    this.viewOnly = true;
    this.editingOrder.set(order);
    this.dialogTitle = `ใบสั่งขาย ${order.order_number}`;
    this.form.patchValue({
      customer_id: order.customer_id,
      order_date: order.order_date,
      requested_date: order.requested_date ?? '',
      vat_rate: order.vat_rate,
      notes: order.notes,
    });
    this.lines.clear();
    order.lines.forEach(() => this.addLine());
    this.showDialog = true;
  }

  addLine() {
    this.lines.push(this.fb.group({
      finished_good_id: ['', Validators.required],
      ordered_qty:      [null, [Validators.required, Validators.min(0.001)]],
      unit_price:       [null, [Validators.required, Validators.min(0)]],
      discount_pct:     [0],
    }));
  }

  removeLine(i: number) { this.lines.removeAt(i); }

  getFGName(id: string): string {
    return this.fgOptions().find(f => f.id === id)?.name ?? '';
  }

  lineTotal(i: number): number {
    const l = this.lines.at(i).value;
    if (!l.ordered_qty || !l.unit_price) return 0;
    return l.ordered_qty * l.unit_price * (1 - (l.discount_pct || 0) / 100);
  }

  subtotal(): number {
    let s = 0;
    for (let i = 0; i < this.lines.length; i++) s += this.lineTotal(i);
    return s;
  }

  vatAmount(): number { return this.subtotal() * ((this.form.value.vat_rate ?? 7) / 100); }
  totalAmount(): number { return this.subtotal() + this.vatAmount(); }

  save() {
    this.form.markAllAsTouched();
    if (this.form.invalid) return;
    const v = this.form.value;
    const lines: CreateSOLineRequest[] = this.lines.value.map((l: any) => ({
      finished_good_id: l.finished_good_id,
      ordered_qty: l.ordered_qty,
      unit_price: l.unit_price,
      discount_pct: l.discount_pct || 0,
    }));
    this.svc.createSalesOrder({
      customer_id: v.customer_id!,
      order_date: v.order_date!,
      requested_date: v.requested_date || undefined,
      vat_rate: v.vat_rate ?? 7,
      notes: v.notes || '',
      lines,
    }).subscribe({
      next: () => {
        this.toast.add({ severity: 'success', summary: 'สำเร็จ', detail: 'สร้างใบสั่งขายแล้ว' });
        this.showDialog = false;
        this.load();
      },
      error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
    });
  }

  confirm_(order: SalesOrder) {
    this.confirm.confirm({
      message: `ยืนยันใบสั่งขาย ${order.order_number}?`,
      accept: () => this.svc.confirmSalesOrder(order.id).subscribe({
        next: () => { this.toast.add({ severity: 'success', summary: 'ยืนยันแล้ว' }); this.load(); },
        error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
      }),
    });
  }

  cancel(order: SalesOrder) {
    this.confirm.confirm({
      message: `ยกเลิกใบสั่งขาย ${order.order_number}?`,
      accept: () => this.svc.cancelSalesOrder(order.id).subscribe({
        next: () => { this.toast.add({ severity: 'success', summary: 'ยกเลิกแล้ว' }); this.load(); },
        error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
      }),
    });
  }

  statusLabel(s: string): string {
    return { draft: 'ร่าง', confirmed: 'ยืนยันแล้ว', picking: 'กำลังจัดส่ง', dispatched: 'จัดส่งแล้ว', invoiced: 'ออกใบแจ้งหนี้แล้ว', cancelled: 'ยกเลิก' }[s] ?? s;
  }
  statusSeverity(s: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    return ({ draft: 'secondary', confirmed: 'info', picking: 'warn', dispatched: 'success', invoiced: 'success', cancelled: 'danger' } as any)[s] ?? 'secondary';
  }
}
