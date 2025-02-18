import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { SearchBarComponent } from './search-bar/search-bar.component';
import { LoginComponent } from './login/login.component';
import { MsalGuard } from '@azure/msal-angular';
import { AuthGuard } from './auth.guard';
import { Router } from '@angular/router';
import { InteractionType } from '@azure/msal-browser';

const routes: Routes = [
  {
    path: '',
    pathMatch: 'full',
    redirectTo: '', // Default route redirects to search-bar
  },
  {
    path: '',
    component: SearchBarComponent,
    canActivate: [MsalGuard], // Protect search-bar with MSAL Guard
  },
  {
    path: 'login',
    component: LoginComponent,
    canActivate: [MsalGuard], // Protect login page from logged-in users
    data: {
      // Custom data to redirect users if they're logged in
      msalGuardConfig: {
        interactionType: InteractionType.Redirect,
        authRequest: {
          scopes: ['user.read'],
        },
      },
    },
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: false })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
