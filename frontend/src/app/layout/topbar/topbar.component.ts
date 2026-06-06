import { Component, computed, inject } from '@angular/core';
import { NgIf } from '@angular/common';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-topbar',
  standalone: true,
  imports: [NgIf],
  templateUrl: './topbar.component.html',
  styleUrl: './topbar.component.scss'
})
export class TopbarComponent {
  private auth = inject(AuthService);

  user = this.auth.currentUser;

  roleLabel = computed(() => {
    const roleMap: Record<string, string> = {
      admin: 'ผู้ดูแลระบบ',
      warehouse_manager: 'คลังสินค้า',
      production_manager: 'ฝ่ายผลิต',
      sales_admin: 'ฝ่ายขาย',
      delivery: 'จัดส่ง',
    };
    return roleMap[this.auth.userRole() ?? ''] ?? '';
  });

  logout(): void {
    this.auth.logout();
  }
}
