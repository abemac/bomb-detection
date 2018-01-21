import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';


import { AppComponent } from './app.component';

import {HttpClientModule} from '@angular/common/http'
import {ApiService} from './api.service';
import { NodeViewComponent } from './nodeview/nodeview.component'

import {MatTabsModule} from '@angular/material/tabs';
import {NoopAnimationsModule} from '@angular/platform-browser/animations';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatSliderModule} from '@angular/material/slider';
import 'hammerjs'
import { FormsModule, ReactiveFormsModule } from '@angular/forms';


@NgModule({
  declarations: [
    AppComponent,
    NodeViewComponent,
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    MatTabsModule,
    NoopAnimationsModule,
    MatGridListModule,
    MatSliderModule,
    FormsModule
  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule { }
