import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { MsalGuard, MsalInterceptor, MsalModule, MsalRedirectComponent, MsalService } from '@azure/msal-angular';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { FormsModule } from '@angular/forms';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';

import { AppRoutingModule } from './app-routing.module';
import { msalConfig, msalGuardConfig, msalInterceptorConfig } from './auth-config';
import { InteractionType, PublicClientApplication } from '@azure/msal-browser';
import { SearchBarComponent } from './search-bar/search-bar.component';
import { LoginComponent } from './login/login.component';
import { AppComponent } from './app.component';
import { AuthGuard } from './auth.guard';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { MatSortModule } from '@angular/material/sort';




@NgModule({
  declarations: [
    AppComponent,
    SearchBarComponent,
    LoginComponent
    ],
    schemas: [
  ],
    imports: [
      MatCardModule,
      MatTableModule,
      MatPaginatorModule,
      MatSortModule,
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    FormsModule,
    MatFormFieldModule,
      MatInputModule,
    MatIconModule,
      MatButtonModule,
    AppRoutingModule,
    MsalModule.forRoot(
      new PublicClientApplication(msalConfig),
        {
            interactionType: InteractionType.Popup, // Msal Guard Configuration
            authRequest: {
                scopes: ["user.read"],
            },
        },
        {
            interactionType: InteractionType.Popup, // MSAL Interceptor Configuration
            protectedResourceMap: new Map([
                ["https://graph.microsoft.com/v1.0/me", ["user.read"]]
            ]),
        }
    ),
  ],
    providers: [
        {
            provide: HTTP_INTERCEPTORS,
            useClass: MsalInterceptor,
            multi: true,
        },
        MsalGuard,
        AuthGuard,
        provideAnimationsAsync(),
    ],

  bootstrap: [AppComponent, MsalRedirectComponent]
})
export class AppModule { }
