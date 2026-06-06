import { Component } from '@angular/core';
import { RouterOutlet, RouterLink, RouterLinkActive } from '@angular/router';

@Component({
  selector: 'app-finished-goods-shell',
  standalone: true,
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './finished-goods-shell.component.html',
  styleUrl: './finished-goods-shell.component.scss',
})
export class FinishedGoodsShellComponent {}
