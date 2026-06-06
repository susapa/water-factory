import { Component, computed, inject } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'app-sales-shell',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './sales-shell.component.html',
  styleUrl: './sales-shell.component.scss',
})
export class SalesShellComponent {
  private auth = inject(AuthService);

  isDelivery = computed(() => this.auth.userRole() === 'delivery');
  isSalesAdmin = computed(() => this.auth.userRole() === 'sales_admin' || this.auth.userRole() === 'admin');
}
