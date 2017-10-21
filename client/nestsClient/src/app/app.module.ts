import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { HttpModule } from '@angular/http';

import { HttpService } from './services/http.service'
import { SessionService } from './services/session.service'

import { AppComponent } from './app.component';
import { DrawerComponent } from './drawer/drawer.component';


@NgModule({
  declarations: [
    AppComponent,
    DrawerComponent,
  ],
  imports: [
    BrowserModule,
    HttpModule
  ],
  providers: [HttpService, SessionService],
  bootstrap: [AppComponent]
})
export class AppModule { }
