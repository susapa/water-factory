import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-master-data-shell',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './master-data-shell.component.html',
  styleUrl: './master-data-shell.component.scss',
})
export class MasterDataShellComponent {}
