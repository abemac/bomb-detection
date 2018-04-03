import { Component, OnInit,Inject} from '@angular/core';
import {MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material';
import { HttpClient,HttpHeaders } from '@angular/common/http';
import { TabControlService } from '../tab-control.service';

@Component({
  selector: 'app-simchooser',
  templateUrl: './simchooser.component.html',
  styleUrls: ['./simchooser.component.css']
})
export class SimchooserComponent implements OnInit {

  files : string[]=new Array<string>() 
  jsontext:any
  description:string;
  filename:string;
  error:boolean=false;
  errorMSG:string=""
  success=false;
  starting=false;
  closeStr="Cancel"
  selected:boolean =false;
  constructor(private tabs : TabControlService,
    public dialogRef: MatDialogRef<SimchooserComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,private http: HttpClient) { 
      this.http.get("/GetConfig").toPromise().then( resp =>{
        for (var file of resp['files']){
          this.files.push(file)
        }
      })
    }

  onChoose(event){
    this.selected=true;
    this.filename=event.value;
    this.http.get("/GetConfig?filename="+event.value).toPromise().then( resp =>{
      this.jsontext=resp
      this.description=resp['description']
    })
  }
  ngOnInit() {
  }
  onStart(){
    this.starting=true;
    this.http.get("/StartSim?filename="+this.filename).toPromise().then( resp =>{
      this.starting=false;
      this.error=false;
      this.success=true;
      this.errorMSG=""
      this.closeStr="Close"
      console.log(resp)
    }).catch(err=>{
      this.starting=false;
      this.error=true;
      this.errorMSG=err.error
      this.success=false
      console.log(err.error)
    })
  }
  onCancelClick(): void {
    if(this.success){
      this.dialogRef.close("SUCCESS");
    }else{
      this.dialogRef.close();
    }
  }
  goToCreator(){
    this.dialogRef.close("CREATE");
    this.tabs.setIndex(2)
  }
}
