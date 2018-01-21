import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';

import {HttpClientModule} from '@angular/common/http'
import {ApiService} from './api.service';
import { NodeViewComponent } from './nodeview/nodeview.component'
@NgModule({
  declarations: [
    AppComponent,
    NodeViewComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }
