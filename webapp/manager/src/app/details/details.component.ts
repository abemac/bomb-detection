import { Component,ViewChild,AfterViewInit } from '@angular/core';
import {MatTableDataSource,MatPaginator} from '@angular/material';
import { ApiService } from '../api.service';
import {NODEDATA} from '../types'
@Component({
  selector: 'app-details',
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.css']
})
export class DetailsComponent implements AfterViewInit {
  displayedColumns = ['id', 'latitude', 'longitude'];
  private dataSource :any;
  
  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(private api: ApiService) {
   
  
   }


  ngAfterViewInit() {
    this.api.updateNodes().then(data=>{
      this.dataSource = new MatTableDataSource<NODEDATA>(this.api.nodeData);
      this.dataSource.paginator = this.paginator;

    })
    
  }
  
  
}