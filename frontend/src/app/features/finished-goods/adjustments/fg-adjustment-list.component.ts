import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { FormBuilder, FormArray, ReactiveFormsModule, Validators } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { InputNumberModule } from 'primeng/inputnumber';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { TextareaModule } from 'primeng/textarea';
import { ConfirmDialogModule } from 'primeng/confirmdialog';
import { MessageService, ConfirmationService } from 'primeng/api';
import { FinishedGoodsInventoryService } from '../../../core/services/finished-goods-inventory.service';
import { FGAdjustment, FGStockLot } from '../../../core/models/finished-goods-inventory.model';
import { AuthService } from '../../../core/auth/auth.service';

@Component({
  selector: 'app-fg-adjustment-list',
  standalone: true,
  imports: [
    DecimalPipe,
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule, ConfirmDialogModule,
    InputTextModule, InputNumberModule, SelectModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './fg-adjustment-list.component.html',
  styleUrl: './fg-adjustment-list.component.scss',
})
export class FgAdjustmentListComponent implements OnInit {
  private svc = inject(FinishedGoodsInventoryService);
  private auth = inject(AuthService);
  private msg = inject(MessageService);
  private confirm = inject(ConfirmationService);
  private fb = inject(FormBuilder);

  items = signal<FGAdjustment[]>([]);
  lots = signal<FGStockLot[]>([]);
  loading = signal(false);
  showDialog = false;
  viewOnly = false;
  editingAdj = signal<FGAdjustment | null>(null);
  searchText = signal('');

  canApprove = computed(() => {
    const role = this.auth.userRole();
    return role === 'admin' || role === 'warehouse_manager';
  });

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(a =>
      a.adj_number.toLowerCase().includes(q) ||
      a.type.toLowerCase().includes(q)
    ) : all;
  });

  typeOptions = [
    { label: 'นับสต็อก (Cycle Count)', value: 'cycle_count' },
    { label: 'ตัดออก (Write-off)', value: 'write_off' },
    { label: 'เพิ่มเข้า (Write-in)', value: 'write_in' },
  ];
  typeLabel = (type: string) => this.typeOptions.find(o => o.value === type)?.label ?? type;

  form = this.fb.group({
    adjustment_date: ['', Validators.required],
    type: ['cycle_count', Validators.required],
    notes: [''],
    lines: this.fb.array<ReturnType<typeof this.newLine>>([]),
  });

  get dialogTitle() {
    if (this.viewOnly) return 'รายละเอียดการปรับปรุงสต็อก';
    return 'สร้างใบปรับปรุงสต็อก';
  }
  get lines(): FormArray { return this.form.get('lines') as FormArray; }

  adjStatusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    if (status === 'approved') return 'success';
    if (status === 'cancelled') return 'danger';
    return 'warn';
  }
  adjStatusLabel(status: string): string {
    if (status === 'approved') return 'อนุมัติแล้ว';
    if (status === 'cancelled') return 'ยกเลิก';
    return 'รอดำเนินการ';
  }

  ngOnInit() {
    this.load();
    this.svc.listStockLots().subscribe(d => this.lots.set((d ?? []).filter(l => l.status === 'available')));
  }

  load() {
    this.loading.set(true);
    this.svc.listAdjustments().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  newLine() {
    return this.fb.group({
      fg_stock_lot_id: ['', Validators.required],
      counted_qty: [null as number | null, [Validators.required, Validators.min(0)]],
      reason: [''],
    });
  }

  addLine() { this.lines.push(this.newLine()); }
  removeLine(i: number) { this.lines.removeAt(i); }

  openCreate() {
    this.viewOnly = false;
    this.editingAdj.set(null);
    this.form.reset({ adjustment_date: '', type: 'cycle_count', notes: '' });
    this.lines.clear();
    this.addLine();
    this.form.enable();
    this.showDialog = true;
  }

  openView(item: FGAdjustment) {
    this.viewOnly = true;
    this.svc.getAdjustment(item.id).subscribe(adj => {
      this.editingAdj.set(adj);
      this.form.patchValue({ adjustment_date: adj.adjustment_date, type: adj.type, notes: adj.notes });
      this.lines.clear();
      (adj.lines ?? []).forEach(l => {
        const g = this.newLine();
        g.patchValue({ fg_stock_lot_id: l.fg_stock_lot_id, counted_qty: l.counted_qty, reason: l.reason });
        this.lines.push(g);
      });
      this.form.disable();
      this.showDialog = true;
    });
  }

  save() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    this.svc.createAdjustment(val).subscribe({
      next: () => {
        this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'สร้างใบปรับปรุงสต็อกเรียบร้อย' });
        this.showDialog = false;
        this.load();
      },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  approveAdj(item: FGAdjustment) {
    this.confirm.confirm({
      message: `อนุมัติใบปรับปรุง ${item.adj_number}? ระบบจะอัปเดตปริมาณสต็อกทันที`,
      header: 'อนุมัติการปรับปรุงสต็อก',
      icon: 'pi pi-check-circle',
      acceptLabel: 'อนุมัติ',
      rejectLabel: 'ยกเลิก',
      accept: () => {
        this.svc.approveAdjustment(item.id).subscribe({
          next: () => {
            this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: `อนุมัติ ${item.adj_number} เรียบร้อย` });
            this.load();
          },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  cancelAdj(item: FGAdjustment) {
    this.confirm.confirm({
      message: `ยกเลิกใบปรับปรุง ${item.adj_number}?`,
      header: 'ยกเลิก',
      icon: 'pi pi-times-circle',
      acceptLabel: 'ยืนยัน',
      rejectLabel: 'ไม่',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.svc.cancelAdjustment(item.id).subscribe({
          next: () => {
            this.msg.add({ severity: 'info', summary: 'ยกเลิกแล้ว', detail: `${item.adj_number} ถูกยกเลิก` });
            this.load();
          },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  getLotLabel(lotId: string): string {
    const lot = this.lots().find(l => l.id === lotId);
    return lot ? `${lot.batch_number} — ${lot.fg_code} (${lot.current_qty} ${lot.uom_code})` : lotId;
  }

  getSystemQty(lotId: string): number | null {
    const lot = this.lots().find(l => l.id === lotId);
    return lot ? lot.current_qty : null;
  }

  varianceClass(variance: number): string {
    if (variance > 0) return 'var-positive';
    if (variance < 0) return 'var-negative';
    return '';
  }
}
