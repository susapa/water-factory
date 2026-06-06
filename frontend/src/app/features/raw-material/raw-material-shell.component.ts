import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-raw-material-shell',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './raw-material-shell.component.html',
  styleUrl: './raw-material-shell.component.scss',
})
export class RawMaterialShellComponent {}
