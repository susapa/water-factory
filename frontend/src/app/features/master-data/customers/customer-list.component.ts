import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { InputNumberModule } from 'primeng/inputnumber';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { TextareaModule } from 'primeng/textarea';
import { MessageService } from 'primeng/api';
import { MasterDataService } from '../../../core/services/master-data.service';
import { Customer } from '../../../core/models/master-data.model';

@Component({
  selector: 'app-customer-list',
  standalone: true,
  imports: [
    DecimalPipe, ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule,
    InputTextModule, InputNumberModule,
    TagModule, ToastModule, TextareaModule,
  ],
  providers: [MessageService],
  templateUrl: './customer-list.component.html',
  styleUrl: './customer-list.component.scss',
})
export class CustomerListComponent implements OnInit {
  private svc = inject(MasterDataService);
  private msg = inject(MessageService);
  private fb = inject(FormBuilder);

  items = signal<Customer[]>([]);
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
    tax_id: [''], address: [''], phone: [''], email: [''],
    credit_limit: [0], credit_days: [30],
  });

  get dialogTitle() { return this.editingId ? 'แก้ไขลูกค้า' : 'เพิ่มลูกค้า'; }

  ngOnInit() { this.load(); }

  load() {
    this.loading.set(true);
    this.svc.listCustomers().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  openCreate() {
    this.editingId = null;
    this.form.reset({ code: '', name: '', tax_id: '', address: '', phone: '', email: '', credit_limit: 0, credit_days: 30 });
    this.form.get('code')?.enable();
    this.showDialog = true;
  }

  openEdit(item: Customer) {
    this.editingId = item.id;
    this.form.patchValue({ code: item.code, name: item.name, tax_id: item.tax_id, address: item.address, phone: item.phone, email: item.email, credit_limit: item.credit_limit, credit_days: item.credit_days });
    this.form.get('code')?.disable();
    this.showDialog = true;
  }

  save() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    const obs = this.editingId ? this.svc.updateCustomer(this.editingId, val) : this.svc.createCustomer(val);
    obs.subscribe({
      next: () => { this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'บันทึกข้อมูลเรียบร้อย' }); this.showDialog = false; this.load(); },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  toggleActive(item: Customer) {
    this.svc.setCustomerActive(item.id, !item.is_active).subscribe({
      next: () => this.load(),
      error: () => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: 'ไม่สามารถเปลี่ยนสถานะได้' }),
    });
  }
}
