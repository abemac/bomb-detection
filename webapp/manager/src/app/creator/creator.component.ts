import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
declare var jquery:any;
declare var $ :any;

@Component({
  selector: 'app-creator',
  templateUrl: './creator.component.html',
  styleUrls: ['./creator.component.css']
})
export class CreatorComponent implements OnInit {
  firstFormGroup: FormGroup;
  secondFormGroup: FormGroup;

  
  nodeConfigRows: ConfigRow[]=new Array<ConfigRow>();

  constructor(private _formBuilder: FormBuilder) {  
    
  
  }

  ngOnInit() {
    this.firstFormGroup = this._formBuilder.group({
      firstCtrl: ['', Validators.required]
    });
    this.secondFormGroup = this._formBuilder.group({
      
    });
    this.addRow();
    
  }
  addRow(){
    this.nodeConfigRows.push(new ConfigRow());
  }
  delrow(index){
    this.nodeConfigRows.splice(index,1)
  }

  downloadConfig(){

  }
  saveConfig(){
    
  }

}
class ConfigRow{
  north:number=0.0;
  east:number=0.0;
  south:number=0.0;
  west:number=0.0;
  num:number=0.0;
  supernode:boolean=false;
  group:boolean=false;
  constructor(){
    this.north=0.0;
    this.east=0.0;
    this.south=0.0;
    this.west=0.0;
    this.num=0;
    this.supernode=false;
    this.group=false;

  }
}