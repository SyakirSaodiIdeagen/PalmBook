import { HttpClientModule } from '@angular/common/http';
import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { MsalModule, MsalRedirectComponent } from '@azure/msal-angular';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { msalConfig, msalGuardConfig, msalInterceptorConfig } from './auth-config';
import { InteractionType, PublicClientApplication } from '@azure/msal-browser';

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule, HttpClientModule,
    AppRoutingModule,
    MsalModule.forRoot(
      new PublicClientApplication(msalConfig),
      msalGuardConfig,
      msalInterceptorConfig
    ),
  ],
  providers: [],
  bootstrap: [AppComponent, MsalRedirectComponent]
})
export class AppModule { }
