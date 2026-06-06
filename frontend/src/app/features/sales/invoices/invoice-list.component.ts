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
import { Invoice, SalesOrder } from '../../../core/models/sales.model';

@Component({
  selector: 'app-invoice-list',
  standalone: true,
  imports: [
    CommonModule, ReactiveFormsModule, DecimalPipe,
    TableModule, ButtonModule, TagModule, DialogModule,
    ToastModule, ConfirmDialogModule, SelectModule, InputTextModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './invoice-list.component.html',
  styleUrl: './invoice-list.component.scss',
})
export class InvoiceListComponent implements OnInit {
  private svc = inject(SalesService);
  private fb = inject(FormBuilder);
  private toast = inject(MessageService);
  private confirm = inject(ConfirmationService);

  items = signal<Invoice[]>([]);
  loading = signal(false);
  searchText = signal('');
  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    return this.items().filter(i =>
      i.invoice_number.toLowerCase().includes(q) ||
      i.so_number.toLowerCase().includes(q) ||
      i.customer_name.toLowerCase().includes(q) ||
      i.status.toLowerCase().includes(q)
    );
  });

  dispatchedSOs = signal<SalesOrder[]>([]);
  showCreateDialog = false;
  showPayDialog = false;
  payingInvoice = signal<Invoice | null>(null);
  showDetailDialog = false;
  detailInvoice = signal<Invoice | null>(null);

  createForm = this.fb.group({
    sales_order_id:    ['', Validators.required],
    invoice_date:      [new Date().toISOString().slice(0, 10), Validators.required],
    due_date:          [''],
  });

  payForm = this.fb.group({
    payment_date:   [new Date().toISOString().slice(0, 10), Validators.required],
    payment_method: ['เงินสด'],
  });

  soOptions = computed(() =>
    this.dispatchedSOs().map(o => ({ id: o.id, label: `${o.order_number} — ${o.customer_name} (${(o.total_amount).toLocaleString()} บ.)` }))
  );

  ngOnInit() { this.load(); }

  load() {
    this.loading.set(true);
    this.svc.listInvoices().subscribe({
      next: d => { this.items.set(d); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  loadSOOptions() {
    this.svc.listSalesOrders().subscribe(orders =>
      this.dispatchedSOs.set(orders.filter(o => o.status === 'dispatched'))
    );
  }

  openCreate() {
    this.createForm.reset({
      sales_order_id: '',
      invoice_date: new Date().toISOString().slice(0, 10),
      due_date: '',
    });
    this.loadSOOptions();
    this.showCreateDialog = true;
  }

  openPay(inv: Invoice) {
    this.payingInvoice.set(inv);
    this.payForm.reset({
      payment_date: new Date().toISOString().slice(0, 10),
      payment_method: 'เงินสด',
    });
    this.showPayDialog = true;
  }

  openDetail(inv: Invoice) {
    this.detailInvoice.set(inv);
    this.showDetailDialog = true;
  }

  save() {
    this.createForm.markAllAsTouched();
    if (this.createForm.invalid) return;
    const v = this.createForm.value;
    this.svc.createInvoice({
      sales_order_id: v.sales_order_id!,
      invoice_date: v.invoice_date!,
      due_date: v.due_date || undefined,
    }).subscribe({
      next: () => {
        this.toast.add({ severity: 'success', summary: 'สำเร็จ', detail: 'ออกใบแจ้งหนี้แล้ว' });
        this.showCreateDialog = false;
        this.load();
      },
      error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
    });
  }

  savePay() {
    this.payForm.markAllAsTouched();
    if (this.payForm.invalid || !this.payingInvoice()) return;
    const v = this.payForm.value;
    this.svc.markInvoicePaid(this.payingInvoice()!.id, {
      payment_date: v.payment_date!,
      payment_method: v.payment_method || 'เงินสด',
    }).subscribe({
      next: () => {
        this.toast.add({ severity: 'success', summary: 'บันทึกการรับชำระแล้ว' });
        this.showPayDialog = false;
        this.load();
      },
      error: e => this.toast.add({ severity: 'error', summary: 'ผิดพลาด', detail: e.error?.error || e.message }),
    });
  }

  statusLabel(s: string): string {
    return { issued: 'ออกใบแล้ว', paid: 'ชำระแล้ว', overdue: 'เกินกำหนด', cancelled: 'ยกเลิก' }[s] ?? s;
  }
  statusSeverity(s: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    return ({ issued: 'info', paid: 'success', overdue: 'danger', cancelled: 'secondary' } as any)[s] ?? 'secondary';
  }
}
