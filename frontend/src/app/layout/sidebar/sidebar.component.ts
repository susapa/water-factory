import { Component, computed, inject } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { NgFor } from '@angular/common';
import { AuthService } from '../../core/auth/auth.service';
import { UserRole } from '../../core/models/user.model';

interface NavItem {
  label: string;
  icon: string;
  route: string;
  roles: UserRole[];
}

const NAV_ITEMS: NavItem[] = [
  { label: 'แดชบอร์ด',         icon: 'pi pi-home',          route: '/dashboard',        roles: ['admin','warehouse_manager','production_manager','sales_admin'] },
  { label: 'ข้อมูลหลัก',       icon: 'pi pi-database',      route: '/master-data',      roles: ['admin'] },
  { label: 'วัตถุดิบ',         icon: 'pi pi-box',           route: '/raw-material',     roles: ['admin','warehouse_manager','production_manager'] },
  { label: 'ผลิต',             icon: 'pi pi-cog',           route: '/production',       roles: ['admin','warehouse_manager','production_manager'] },
  { label: 'สินค้าสำเร็จรูป', icon: 'pi pi-shopping-bag',  route: '/finished-goods',   roles: ['admin','warehouse_manager','sales_admin'] },
  { label: 'ขายและจัดส่ง',    icon: 'pi pi-truck',         route: '/sales',            roles: ['admin','sales_admin','delivery'] },
  { label: 'รายงาน',           icon: 'pi pi-chart-bar',     route: '/reports',          roles: ['admin','warehouse_manager','production_manager','sales_admin'] },
  { label: 'จัดการผู้ใช้',    icon: 'pi pi-users',         route: '/admin',            roles: ['admin'] },
];

@Component({
  selector: 'app-sidebar',
  standalone: true,
  imports: [RouterLink, RouterLinkActive, NgFor],
  templateUrl: './sidebar.component.html',
  styleUrl: './sidebar.component.scss'
})
export class SidebarComponent {
  private auth = inject(AuthService);

  visibleItems = computed(() => {
    const role = this.auth.userRole();
    if (!role) return [];
    return NAV_ITEMS.filter(item => item.roles.includes(role));
  });
}
