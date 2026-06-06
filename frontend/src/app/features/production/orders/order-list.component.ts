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
import { ProductionService } from '../../../core/services/production.service';
import { MasterDataService } from '../../../core/services/master-data.service';
import { ProductionOrder, ProductionRequirement } from '../../../core/models/production.model';
import { FinishedGood, BOM } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-order-list',
  standalone: true,
  imports: [
    DecimalPipe,
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule, ConfirmDialogModule,
    InputTextModule, InputNumberModule, SelectModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './order-list.component.html',
  styleUrl: './order-list.component.scss',
})
export class OrderListComponent implements OnInit {
  private svc = inject(ProductionService);
  private mdSvc = inject(MasterDataService);
  private msg = inject(MessageService);
  private confirm = inject(ConfirmationService);
  private fb = inject(FormBuilder);

  items = signal<ProductionOrder[]>([]);
  fgOptions = signal<FinishedGood[]>([]);
  allBoms = signal<BOM[]>([]);
  loading = signal(false);
  searchText = signal('');
  selectedFGId = signal<string>('');

  // Dialog flags
  showCreateDialog = false;
  showRequirementsDialog = false;
  showIssueDialog = false;
  showYieldDialog = false;

  selectedOrder = signal<ProductionOrder | null>(null);
  isEditMode = false;

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(o =>
      o.order_number.toLowerCase().includes(q) ||
      o.fg_name.toLowerCase().includes(q)
    ) : all;
  });

  filteredBoms = computed(() => {
    const fgId = this.selectedFGId();
    return fgId ? this.allBoms().filter(b => b.finished_good_id === fgId) : [];
  });

  // Create / Edit form
  createForm = this.fb.group({
    finished_good_id: ['', Validators.required],
    bom_id: [null as string | null],
    planned_qty: [null as number | null, [Validators.required, Validators.min(0.001)]],
    planned_start_date: [null as string | null],
    planned_end_date: [null as string | null],
    notes: [''],
  });

  // Issue RM form — dynamic lines based on requirements
  issueLines = this.fb.array<ReturnType<typeof this.newIssueLine>>([]);

  // Record yield form
  yieldForm = this.fb.group({
    yield_qty: [null as number | null, [Validators.required, Validators.min(0.001)]],
    defect_qty: [0 as number | null],
    waste_qty: [0 as number | null],
    batch_number: ['', Validators.required],
    production_date: ['', Validators.required],
    expiry_date: ['', Validators.required],
    notes: [''],
    defect_details: this.fb.array<ReturnType<typeof this.newDefectLine>>([]),
  });

  get defectDetails(): FormArray { return this.yieldForm.get('defect_details') as FormArray; }

  ngOnInit() {
    this.load();
    this.mdSvc.listFinishedGoods().subscribe(data => this.fgOptions.set((data ?? []).filter(fg => fg.is_active)));
    this.mdSvc.listBOMs().subscribe(data => this.allBoms.set(data ?? []));
    this.createForm.get('finished_good_id')!.valueChanges.subscribe(v => {
      this.selectedFGId.set(v ?? '');
      this.createForm.get('bom_id')?.setValue(null);
    });
  }

  load() {
    this.loading.set(true);
    this.svc.listOrders().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  // ── Status helpers ───────────────────────────────────────────────────────────

  statusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    const map: Record<string, 'success' | 'secondary' | 'info' | 'warn' | 'danger'> = {
      draft: 'secondary',
      confirmed: 'info',
      rm_issued: 'warn',
      in_progress: 'warn',
      completed: 'success',
      cancelled: 'danger',
    };
    return map[status] ?? 'secondary';
  }

  statusLabel(status: string): string {
    const map: Record<string, string> = {
      draft: 'แบบร่าง',
      confirmed: 'ยืนยันแล้ว',
      rm_issued: 'เบิกวัตถุดิบแล้ว',
      in_progress: 'กำลังผลิต',
      completed: 'เสร็จสิ้น',
      cancelled: 'ยกเลิก',
    };
    return map[status] ?? status;
  }

  reqStatusLabel(status: string): string {
    const map: Record<string, string> = { pending: 'รอเบิก', partial: 'เบิกบางส่วน', fully_issued: 'เบิกครบ' };
    return map[status] ?? status;
  }

  reqStatusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    if (status === 'fully_issued') return 'success';
    if (status === 'partial') return 'warn';
    return 'secondary';
  }

  // ── Create / Edit ────────────────────────────────────────────────────────────

  openCreate() {
    this.isEditMode = false;
    this.selectedOrder.set(null);
    this.createForm.reset({ finished_good_id: '', bom_id: null, planned_qty: null, planned_start_date: null, planned_end_date: null, notes: '' });
    this.createForm.enable();
    this.showCreateDialog = true;
  }

  openEdit(order: ProductionOrder) {
    this.isEditMode = true;
    this.selectedOrder.set(order);
    this.createForm.patchValue({
      finished_good_id: order.finished_good_id,
      bom_id: order.bom_id,
      planned_qty: order.planned_qty,
      planned_start_date: order.planned_start_date,
      planned_end_date: order.planned_end_date,
      notes: order.notes,
    });
    this.createForm.get('finished_good_id')?.disable();
    this.showCreateDialog = true;
  }

  saveOrder() {
    if (this.createForm.invalid) { this.createForm.markAllAsTouched(); return; }
    const val = this.createForm.getRawValue();
    const payload = {
      finished_good_id: val.finished_good_id,
      bom_id: val.bom_id || null,
      planned_qty: val.planned_qty,
      planned_start_date: val.planned_start_date || null,
      planned_end_date: val.planned_end_date || null,
      notes: val.notes ?? '',
    };

    if (this.isEditMode) {
      const id = this.selectedOrder()!.id;
      this.svc.updateOrder(id, payload).subscribe({
        next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึกแล้ว' }); this.showCreateDialog = false; this.load(); },
        error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
      });
    } else {
      this.svc.createOrder(payload).subscribe({
        next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'สร้างใบสั่งผลิตแล้ว' }); this.showCreateDialog = false; this.load(); },
        error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
      });
    }
  }

  // ── Requirements dialog ──────────────────────────────────────────────────────

  openRequirements(order: ProductionOrder) {
    this.selectedOrder.set(order);
    this.showRequirementsDialog = true;
  }

  // ── Confirm order ────────────────────────────────────────────────────────────

  confirmOrder(order: ProductionOrder) {
    this.confirm.confirm({
      message: `ยืนยันใบสั่งผลิต ${order.order_number}? ระบบจะคำนวณวัตถุดิบที่ต้องใช้จาก BOM`,
      header: 'ยืนยันใบสั่งผลิต',
      icon: 'pi pi-check-circle',
      acceptLabel: 'ยืนยัน',
      rejectLabel: 'ยกเลิก',
      accept: () => {
        this.svc.confirmOrder(order.id).subscribe({
          next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: `${order.order_number} ยืนยันแล้ว` }); this.load(); },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  // ── Issue RM ─────────────────────────────────────────────────────────────────

  newIssueLine(req: ProductionRequirement) {
    return this.fb.group({
      raw_material_id: [req.raw_material_id],
      rm_name: [req.rm_name],
      uom_code: [req.uom_code],
      required_qty: [req.required_qty],
      issued_qty: [req.issued_qty],
      qty: [req.required_qty - req.issued_qty, [Validators.required, Validators.min(0.001)]],
    });
  }

  openIssueRM(order: ProductionOrder) {
    this.selectedOrder.set(order);
    this.issueLines.clear();
    (order.requirements ?? [])
      .filter(r => r.status !== 'fully_issued')
      .forEach(r => this.issueLines.push(this.newIssueLine(r)));
    this.showIssueDialog = true;
  }

  submitIssueRM() {
    if (this.issueLines.invalid) { this.issueLines.markAllAsTouched(); return; }
    const order = this.selectedOrder()!;
    const lines = this.issueLines.getRawValue().map(l => ({
      raw_material_id: l.raw_material_id,
      qty: l.qty,
    }));
    this.svc.issueRM(order.id, { lines }).subscribe({
      next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'เบิกวัตถุดิบเรียบร้อย' }); this.showIssueDialog = false; this.load(); },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'ไม่สามารถเบิกได้' }),
    });
  }

  // ── Start production ─────────────────────────────────────────────────────────

  startProduction(order: ProductionOrder) {
    this.confirm.confirm({
      message: `เริ่มผลิต ${order.order_number}?`,
      header: 'เริ่มผลิต',
      icon: 'pi pi-play',
      acceptLabel: 'เริ่ม',
      rejectLabel: 'ยกเลิก',
      accept: () => {
        this.svc.startProduction(order.id).subscribe({
          next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: `${order.order_number} เริ่มผลิตแล้ว` }); this.load(); },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  // ── Record yield ─────────────────────────────────────────────────────────────

  newDefectLine() {
    return this.fb.group({
      defect_type: ['', Validators.required],
      qty: [null as number | null, [Validators.required, Validators.min(0.001)]],
      description: [''],
    });
  }

  openRecordYield(order: ProductionOrder) {
    this.selectedOrder.set(order);
    const today = new Date().toISOString().slice(0, 10);
    this.yieldForm.reset({ yield_qty: null, defect_qty: 0, waste_qty: 0, batch_number: '', production_date: today, expiry_date: '', notes: '' });
    this.defectDetails.clear();
    this.showYieldDialog = true;
  }

  addDefectLine() { this.defectDetails.push(this.newDefectLine()); }
  removeDefectLine(i: number) { this.defectDetails.removeAt(i); }

  submitYield() {
    if (this.yieldForm.invalid) { this.yieldForm.markAllAsTouched(); return; }
    const order = this.selectedOrder()!;
    const val = this.yieldForm.getRawValue();
    this.svc.recordYield(order.id, val).subscribe({
      next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึกผลผลิตเรียบร้อย' }); this.showYieldDialog = false; this.load(); },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  // ── Complete ─────────────────────────────────────────────────────────────────

  completeOrder(order: ProductionOrder) {
    this.confirm.confirm({
      message: `เสร็จสิ้นการผลิต ${order.order_number}?`,
      header: 'เสร็จสิ้นการผลิต',
      icon: 'pi pi-flag-fill',
      acceptLabel: 'เสร็จสิ้น',
      rejectLabel: 'ยกเลิก',
      accept: () => {
        this.svc.completeOrder(order.id).subscribe({
          next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: `${order.order_number} เสร็จสิ้นแล้ว` }); this.load(); },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  // ── Cancel ───────────────────────────────────────────────────────────────────

  cancelOrder(order: ProductionOrder) {
    this.confirm.confirm({
      message: `ยกเลิกใบสั่งผลิต ${order.order_number}?`,
      header: 'ยกเลิก',
      icon: 'pi pi-times-circle',
      acceptLabel: 'ยืนยัน',
      rejectLabel: 'ไม่',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.svc.cancelOrder(order.id).subscribe({
          next: () => { this.msg.add({ severity: 'info', summary: 'ยกเลิกแล้ว', detail: `${order.order_number} ถูกยกเลิก` }); this.load(); },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }
}
