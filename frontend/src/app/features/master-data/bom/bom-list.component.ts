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
import { MessageService } from 'primeng/api';
import { MasterDataService } from '../../../core/services/master-data.service';
import { BOM, FinishedGood, RawMaterial } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-bom-list',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule,
    InputTextModule, InputNumberModule, SelectModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService],
  templateUrl: './bom-list.component.html',
  styleUrl: './bom-list.component.scss',
})
export class BomListComponent implements OnInit {
  private svc = inject(MasterDataService);
  private msg = inject(MessageService);
  private fb = inject(FormBuilder);

  items = signal<BOM[]>([]);
  finishedGoods = signal<FinishedGood[]>([]);
  rawMaterials = signal<RawMaterial[]>([]);
  loading = signal(false);
  showDialog = false;
  editingId: string | null = null;
  searchText = signal('');

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(i => i.fg_code.toLowerCase().includes(q) || i.fg_name.toLowerCase().includes(q)) : all;
  });

  form = this.fb.group({
    finished_good_id: ['', Validators.required],
    version: [1, [Validators.required, Validators.min(1)]],
    is_active: [true],
    effective_date: ['', Validators.required],
    notes: [''],
    lines: this.fb.array<ReturnType<typeof this.newLine>>([]),
  });

  get dialogTitle() { return this.editingId ? 'แก้ไข BOM' : 'สร้าง BOM'; }
  get lines(): FormArray { return this.form.get('lines') as FormArray; }

  ngOnInit() {
    this.load();
    this.svc.listFinishedGoods().subscribe(d => this.finishedGoods.set((d ?? []).filter(f => f.is_active)));
    this.svc.listRawMaterials().subscribe(d => this.rawMaterials.set((d ?? []).filter(r => r.is_active)));
  }

  load() {
    this.loading.set(true);
    this.svc.listBOMs().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  newLine() {
    return this.fb.group({
      raw_material_id: ['', Validators.required],
      qty_per_unit: [null as number | null, [Validators.required, Validators.min(0.0001)]],
      waste_factor: [0],
    });
  }

  addLine() { this.lines.push(this.newLine()); }
  removeLine(i: number) { this.lines.removeAt(i); }

  openCreate() {
    this.editingId = null;
    this.form.reset({ finished_good_id: '', version: 1, is_active: true, effective_date: '', notes: '' });
    this.lines.clear();
    this.addLine();
    this.form.get('finished_good_id')?.enable();
    this.showDialog = true;
  }

  openEdit(item: BOM) {
    this.editingId = item.id;
    this.svc.getBOM(item.id).subscribe(bom => {
      this.form.patchValue({ finished_good_id: bom.finished_good_id, version: bom.version, is_active: bom.is_active, effective_date: bom.effective_date, notes: bom.notes });
      this.form.get('finished_good_id')?.disable();
      this.lines.clear();
      (bom.lines ?? []).forEach(l => {
        const g = this.newLine();
        g.patchValue({ raw_material_id: l.raw_material_id, qty_per_unit: l.qty_per_unit, waste_factor: l.waste_factor });
        this.lines.push(g);
      });
      if (this.lines.length === 0) this.addLine();
      this.showDialog = true;
    });
  }

  save() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    const obs = this.editingId ? this.svc.updateBOM(this.editingId, val) : this.svc.createBOM(val);
    obs.subscribe({
      next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึก BOM เรียบร้อย' }); this.showDialog = false; this.load(); },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }
}
