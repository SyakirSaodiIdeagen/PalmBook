import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { MsalService } from '@azure/msal-angular';
import { AuthService } from '../services/auth/auth.service';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styleUrl: './login.component.css'
})
export class LoginComponent {
  constructor(private msalService: MsalService, private authService: AuthService, private router: Router) { }
  ngOnInit(): void {
    const accounts = this.msalService.instance.getAllAccounts();
    if (accounts.length > 0) {
      // Redirect if already logged in
      this.router.navigate(['/search-bar']);
    }
  }

  // Initiates Microsoft login using redirect
  login() {
    this.authService.microsoftLogin(); // Call the microsoftLogin() method from the AuthService

  }

  redirect() {
    this.router.navigate(['/']);
  }

}
