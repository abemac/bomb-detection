import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import {HttpClientModule} from '@angular/common/http'
import {NodesService} from './nodes.service';
import { NodeViewComponent } from './nodeview/nodeview.component'
import {MatTabsModule} from '@angular/material/tabs';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatGridListModule} from '@angular/material/grid-list';
import {MatSliderModule} from '@angular/material/slider';
import 'hammerjs'
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import {MatCardModule} from '@angular/material/card';
import {MatInputModule} from '@angular/material/input';
import {MatFormFieldModule} from '@angular/material/form-field';
import {MouseListenerDirective } from './mouselistener.directive';
import {MatButtonModule} from '@angular/material/button';
import {MatSlideToggleModule} from '@angular/material/slide-toggle';
import {MatCheckboxModule} from '@angular/material/checkbox';
import {MatTableModule} from '@angular/material/table';
import { DetailsComponent } from './details/details.component';
import {MatPaginatorModule} from '@angular/material/paginator';
import {MatExpansionModule} from '@angular/material/expansion';
import { CreatorComponent } from './creator/creator.component';
import {MatStepperModule} from '@angular/material/stepper';
import { SimchooserComponent } from './simchooser/simchooser.component';
import {MatDialogModule} from '@angular/material/dialog';
import {MatSelectModule} from '@angular/material/select';
import {TabControlService} from './tab-control.service'
import {MatChipsModule} from '@angular/material/chips';
import {MatProgressSpinnerModule} from '@angular/material/progress-spinner';
import {MatIconModule} from '@angular/material/icon';


@NgModule({
  declarations: [
    AppComponent,
    NodeViewComponent,
    MouseListenerDirective,
    DetailsComponent,
    CreatorComponent,
    SimchooserComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    MatTabsModule,
    BrowserAnimationsModule,
    MatGridListModule,
    MatSliderModule,
    FormsModule,
    ReactiveFormsModule,
    MatCardModule,
    MatInputModule,
    MatFormFieldModule,
    MatButtonModule,
    MatSlideToggleModule,
    MatCheckboxModule,
    MatTableModule,
    MatPaginatorModule,
    MatExpansionModule,
    MatStepperModule,
    MatDialogModule,
    MatSelectModule,
    MatChipsModule,
    MatProgressSpinnerModule,
    MatIconModule
  ],
  providers: [NodesService,TabControlService],
  entryComponents: [SimchooserComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }
