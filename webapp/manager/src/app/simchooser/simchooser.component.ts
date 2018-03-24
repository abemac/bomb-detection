import { Component, OnInit,Inject} from '@angular/core';
import {MatDialog, MatDialogRef, MAT_DIALOG_DATA} from '@angular/material';
import { HttpClient,HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-simchooser',
  templateUrl: './simchooser.component.html',
  styleUrls: ['./simchooser.component.css']
})
export class SimchooserComponent implements OnInit {

  files : string[]=new Array<string>() 
  jsontext:any
  description:string;
  constructor(
    public dialogRef: MatDialogRef<SimchooserComponent>,
    @Inject(MAT_DIALOG_DATA) public data: any,private http: HttpClient) { 
      this.http.get("/GetConfig").toPromise().then( resp =>{
        for (var file of resp['files']){
          this.files.push(file)
        }
      })
    }

  onChoose(event){
    this.http.get("/GetConfig?filename="+event.value).toPromise().then( resp =>{
      this.jsontext=resp
      this.description=resp['description']
    })
  }
  ngOnInit() {
  }
  
  onCancelClick(): void {
    this.dialogRef.close();
  }
  goToCreator(){
    this.dialogRef.close();
  }
}
