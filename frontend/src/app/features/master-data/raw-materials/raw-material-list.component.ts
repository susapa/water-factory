import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
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
import { RawMaterial, UOM, RawMaterialCategory } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-raw-material-list',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule,
    InputTextModule, InputNumberModule, SelectModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService],
  templateUrl: './raw-material-list.component.html',
  styleUrl: './raw-material-list.component.scss',
})
export class RawMaterialListComponent implements OnInit {
  private svc = inject(MasterDataService);
  private msg = inject(MessageService);
  private fb = inject(FormBuilder);

  items = signal<RawMaterial[]>([]);
  uoms = signal<UOM[]>([]);
  categories = signal<RawMaterialCategory[]>([]);
  loading = signal(false);
  showDialog = false;
  editingId: string | null = null;
  searchText = signal('');

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q ? all.filter(i => i.code.toLowerCase().includes(q) || i.name.toLowerCase().includes(q)) : all;
  });

  form = this.fb.group({
    code: ['', Validators.required],
    name: ['', Validators.required],
    category_id: [null as number | null],
    uom_id: [null as number | null, Validators.required],
    min_stock_qty: [0],
    reorder_qty: [0],
    description: [''],
  });

  get dialogTitle() { return this.editingId ? 'แก้ไขวัตถุดิบ' : 'เพิ่มวัตถุดิบ'; }

  ngOnInit() {
    this.load();
    this.svc.listUOM().subscribe(d => this.uoms.set(d ?? []));
    this.svc.listCategories().subscribe(d => this.categories.set(d ?? []));
  }

  load() {
    this.loading.set(true);
    this.svc.listRawMaterials().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  openCreate() {
    this.editingId = null;
    this.form.reset({ code: '', name: '', category_id: null, uom_id: null, min_stock_qty: 0, reorder_qty: 0, description: '' });
    this.form.get('code')?.enable();
    this.showDialog = true;
  }

  openEdit(item: RawMaterial) {
    this.editingId = item.id;
    this.form.patchValue({ code: item.code, name: item.name, category_id: item.category_id, uom_id: item.uom_id, min_stock_qty: item.min_stock_qty, reorder_qty: item.reorder_qty, description: item.description });
    this.form.get('code')?.disable();
    this.showDialog = true;
  }

  save() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    const obs = this.editingId ? this.svc.updateRawMaterial(this.editingId, val) : this.svc.createRawMaterial(val);
    obs.subscribe({
      next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึกข้อมูลเรียบร้อย' }); this.showDialog = false; this.load(); },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  toggleActive(item: RawMaterial) {
    this.svc.setRawMaterialActive(item.id, !item.is_active).subscribe({
      next: () => this.load(),
      error: () => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: 'ไม่สามารถเปลี่ยนสถานะได้' }),
    });
  }
}
