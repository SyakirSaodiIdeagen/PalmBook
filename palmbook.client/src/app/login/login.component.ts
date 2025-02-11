import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  constructor(private authService: MsalService, private router: Router) { }

  // Initiates Microsoft login using redirect
  login() {
    this.authService.loginRedirect();
  }
  redirect() {
    this.router.navigate(['/search-bar']);
  }

}
