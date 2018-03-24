import { Component, OnInit } from '@angular/core';
import {FormBuilder, FormGroup, Validators} from '@angular/forms';
import { HttpClient,HttpHeaders } from '@angular/common/http';

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


  constructor(private _formBuilder: FormBuilder,private http: HttpClient) {  
    
  
  }

  ngOnInit() {
    this.firstFormGroup = this._formBuilder.group({
      filename: ['', Validators.required],
      description: ['']
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
  
  saveConfig(){
    var config=[];
    config.push(
      `{"filename":"${this.firstFormGroup.get('filename').value}.json",`,
      `"description":"${this.firstFormGroup.get('description').value}",`,
    )
    config.push(`"rows": [`)
    var first=true;
    this.nodeConfigRows.forEach((row,index,rows)=>{
      if (!first){
        config.push(",");
      }else{
        first=false;
      }
      config.push(
        `{"north":${row.north},`,
        `"east":${row.east},`,
        `"south":${row.south},`,
        `"west":${row.west},`,
        `"num":${row.num},`,
        `"supernode":${row.supernode},`,
        `"group":${row.group}}`
      );
    });
    config.push(`]}`)

    var jsonStr= config.join("")
    var escapedJson= jsonStr.replace(/[\b]/g, '\\b')
                            .replace(/[\f]/g, '\\f')
                            .replace(/[\n]/g, '\\n')
                            .replace(/[\r]/g, '\\r')
                            .replace(/[\t]/g, '\\t')
    var httpOptions = {
      headers: new HttpHeaders({ 'Filename': `${this.firstFormGroup.get('filename').value}.json`})
    };
    this.http.post('/UploadConfig',escapedJson,httpOptions ).toPromise().then(resp=>{});
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