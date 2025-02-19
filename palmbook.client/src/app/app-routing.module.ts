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
    redirectTo: 'login', // Default route redirects to search-bar
  },
  {
    path: 'search-bar',
    component: SearchBarComponent,
  },
  {
    path: 'login',
    component: LoginComponent,
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, { useHash: false })],
  exports: [RouterModule]
})
export class AppRoutingModule { }
