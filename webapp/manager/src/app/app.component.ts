import { Component,ViewChild,OnInit} from '@angular/core';
import { TabControlService } from './tab-control.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit{
  tabIndex=0;
  constructor(private tabs : TabControlService){
    
  }
  title = 'app';
  
  ngOnInit(){
    this.tabs.tabindex.subscribe(index=>{
      this.tabIndex=index
      console.log(index)
    } )
  }
    
    
}
