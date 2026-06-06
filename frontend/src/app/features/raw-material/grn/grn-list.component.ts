import { Component, OnInit, inject, signal, computed } from '@angular/core';
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
import { MasterDataService } from '../../../core/services/master-data.service';
import { RawMaterialInventoryService } from '../../../core/services/raw-material-inventory.service';
import { GRN, WarehouseLocation } from '../../../core/models/raw-material-inventory.model';
import { RawMaterial, Supplier } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-grn-list',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule, ConfirmDialogModule,
    InputTextModule, InputNumberModule, SelectModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService, ConfirmationService],
  templateUrl: './grn-list.component.html',
  styleUrl: './grn-list.component.scss',
})
export class GrnListComponent implements OnInit {
  private svc = inject(RawMaterialInventoryService);
  private mdSvc = inject(MasterDataService);
  private msg = inject(MessageService);
  private confirm = inject(ConfirmationService);
  private fb = inject(FormBuilder);

  items = signal<GRN[]>([]);
  suppliers = signal<Supplier[]>([]);
  rawMaterials = signal<RawMaterial[]>([]);
  locations = signal<WarehouseLocation[]>([]);
  loading = signal(false);
  showDialog = false;
  editingId: string | null = null;
  viewOnly = false;
  searchText = signal('');

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(i =>
      i.grn_number.toLowerCase().includes(q) ||
      i.supplier_name.toLowerCase().includes(q) ||
      i.po_reference.toLowerCase().includes(q)
    ) : all;
  });

  form = this.fb.group({
    supplier_id: [null as string | null],
    received_date: ['', Validators.required],
    po_reference: [''],
    notes: [''],
    lines: this.fb.array<ReturnType<typeof this.newLine>>([]),
  });

  get dialogTitle() {
    if (this.viewOnly) return 'ดูรายละเอียด GRN';
    return this.editingId ? 'แก้ไข GRN (Draft)' : 'สร้างใบรับวัตถุดิบ (GRN)';
  }
  get lines(): FormArray { return this.form.get('lines') as FormArray; }

  grnStatusSeverity(status: string): 'success' | 'secondary' | 'info' | 'warn' | 'danger' {
    if (status === 'confirmed') return 'success';
    if (status === 'cancelled') return 'danger';
    return 'warn';
  }
  grnStatusLabel(status: string): string {
    if (status === 'confirmed') return 'ยืนยันแล้ว';
    if (status === 'cancelled') return 'ยกเลิก';
    return 'ร่าง';
  }

  ngOnInit() {
    this.load();
    this.mdSvc.listSuppliers().subscribe(d => this.suppliers.set(d ?? []));
    this.mdSvc.listRawMaterials().subscribe(d => this.rawMaterials.set((d ?? []).filter(r => r.is_active)));
    this.svc.listWarehouseLocations().subscribe(d => this.locations.set(d ?? []));
  }

  load() {
    this.loading.set(true);
    this.svc.listGRNs().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  newLine() {
    return this.fb.group({
      raw_material_id: ['', Validators.required],
      lot_number: [''],
      received_qty: [null as number | null, [Validators.required, Validators.min(0.0001)]],
      unit_cost: [null as number | null],
      expiry_date: [''],
      location_id: [null as number | null],
      notes: [''],
    });
  }

  addLine() { this.lines.push(this.newLine()); }
  removeLine(i: number) { if (!this.viewOnly) this.lines.removeAt(i); }

  openCreate() {
    this.editingId = null;
    this.viewOnly = false;
    this.form.reset({ supplier_id: null, received_date: '', po_reference: '', notes: '' });
    this.lines.clear();
    this.addLine();
    this.form.enable();
    this.showDialog = true;
  }

  openView(item: GRN) {
    this.editingId = item.id;
    this.viewOnly = true;
    this.svc.getGRN(item.id).subscribe(grn => {
      this.form.patchValue({
        supplier_id: grn.supplier_id,
        received_date: grn.received_date,
        po_reference: grn.po_reference,
        notes: grn.notes,
      });
      this.lines.clear();
      (grn.lines ?? []).forEach(l => {
        const g = this.newLine();
        g.patchValue({
          raw_material_id: l.raw_material_id,
          lot_number: l.lot_number,
          received_qty: l.received_qty,
          unit_cost: l.unit_cost,
          expiry_date: l.expiry_date ?? '',
          location_id: l.location_id,
          notes: l.notes,
        });
        this.lines.push(g);
      });
      if (this.lines.length === 0) this.addLine();
      if (item.status !== 'draft') {
        this.form.disable();
      } else {
        this.form.enable();
        this.viewOnly = false;
        this.editingId = item.id;
      }
      this.showDialog = true;
    });
  }

  save() {
    if (this.viewOnly) { this.showDialog = false; return; }
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    const body = {
      ...val,
      supplier_id: val.supplier_id || null,
      lines: val.lines.map((l: any) => ({
        ...l,
        expiry_date: l.expiry_date || null,
        location_id: l.location_id || null,
      })),
    };
    const obs = this.editingId
      ? this.svc.updateGRN(this.editingId, body)
      : this.svc.createGRN(body);
    obs.subscribe({
      next: () => {
        this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึก GRN เรียบร้อย' });
        this.showDialog = false;
        this.load();
      },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  confirmGRN(item: GRN) {
    this.confirm.confirm({
      message: `ยืนยันการรับวัตถุดิบ ${item.grn_number}? ระบบจะสร้าง stock lots และบันทึกการเคลื่อนไหว`,
      header: 'ยืนยัน GRN',
      icon: 'pi pi-check-circle',
      acceptLabel: 'ยืนยัน',
      rejectLabel: 'ยกเลิก',
      accept: () => {
        this.svc.confirmGRN(item.id).subscribe({
          next: () => {
            this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: `ยืนยัน ${item.grn_number} เรียบร้อย สต็อกถูกอัปเดตแล้ว` });
            this.load();
          },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  cancelGRN(item: GRN) {
    this.confirm.confirm({
      message: `ยกเลิก GRN ${item.grn_number}?`,
      header: 'ยกเลิก GRN',
      icon: 'pi pi-times-circle',
      acceptLabel: 'ยืนยัน',
      rejectLabel: 'ไม่',
      acceptButtonStyleClass: 'p-button-danger',
      accept: () => {
        this.svc.cancelGRN(item.id).subscribe({
          next: () => {
            this.msg.add({ severity: 'info', summary: 'ยกเลิกแล้ว', detail: `${item.grn_number} ถูกยกเลิก` });
            this.load();
          },
          error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
        });
      },
    });
  }

  locationOptions() {
    return this.locations().map(l => ({
      label: l.zone + (l.row_no ? '-' + l.row_no : '') + (l.bay_no ? '-' + l.bay_no : ''),
      value: l.id,
    }));
  }
}
