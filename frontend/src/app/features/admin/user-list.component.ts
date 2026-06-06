import { Component, OnInit, inject, signal, computed } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { TableModule } from 'primeng/table';
import { ButtonModule } from 'primeng/button';
import { DialogModule } from 'primeng/dialog';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { SelectModule } from 'primeng/select';
import { TagModule } from 'primeng/tag';
import { ToastModule } from 'primeng/toast';
import { MessageService } from 'primeng/api';
import { UserService } from '../../core/services/user.service';
import { User, Role } from '../../core/models/user.model';

@Component({
  selector: 'app-user-list',
  standalone: true,
  imports: [
    ReactiveFormsModule,
    TableModule, ButtonModule, DialogModule,
    InputTextModule, PasswordModule, SelectModule,
    TagModule, ToastModule,
  ],
  providers: [MessageService],
  templateUrl: './user-list.component.html',
  styleUrl: './user-list.component.scss',
})
export class UserListComponent implements OnInit {
  private svc = inject(UserService);
  private msg = inject(MessageService);
  private fb = inject(FormBuilder);

  items = signal<User[]>([]);
  roles = signal<Role[]>([]);
  loading = signal(false);
  showDialog = false;
  searchText = signal('');

  filtered = computed(() => {
    const q = this.searchText().toLowerCase();
    const all = this.items();
    return q
      ? all.filter(u =>
          u.full_name.toLowerCase().includes(q) ||
          u.email.toLowerCase().includes(q) ||
          u.role.toLowerCase().includes(q)
        )
      : all;
  });

  form = this.fb.group({
    full_name: ['', Validators.required],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    role_id: [null as number | null, Validators.required],
  });

  roleLabel: Record<string, string> = {
    admin: 'ผู้ดูแลระบบ',
    warehouse_manager: 'ผู้จัดการคลัง',
    production_manager: 'ผู้จัดการผลิต',
    sales_admin: 'ฝ่ายขาย',
    delivery: 'จัดส่ง',
  };

  ngOnInit() {
    this.load();
    this.svc.listRoles().subscribe(d => this.roles.set(d ?? []));
  }

  load() {
    this.loading.set(true);
    this.svc.listUsers().subscribe({
      next: data => { this.items.set(data ?? []); this.loading.set(false); },
      error: () => this.loading.set(false),
    });
  }

  openCreate() {
    this.form.reset({ full_name: '', email: '', password: '', role_id: null });
    this.showDialog = true;
  }

  save() {
    if (this.form.invalid) { this.form.markAllAsTouched(); return; }
    const val = this.form.getRawValue();
    this.svc.createUser({
      full_name: val.full_name!,
      email: val.email!,
      password: val.password!,
      role_id: val.role_id!,
    }).subscribe({
      next: () => {
        this.msg.add({ severity: 'success', summary: 'สำเร็จ', detail: 'สร้างผู้ใช้เรียบร้อย' });
        this.showDialog = false;
        this.load();
      },
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }

  toggleActive(user: User) {
    this.svc.setUserActive(user.id, !user.is_active).subscribe({
      next: () => this.load(),
      error: err => this.msg.add({ severity: 'error', summary: 'ผิดพลาด', detail: err.error?.message ?? 'เกิดข้อผิดพลาด' }),
    });
  }
}
