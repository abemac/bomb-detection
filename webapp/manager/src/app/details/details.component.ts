import { Component,ViewChild,AfterViewInit } from '@angular/core';
import {MatTableDataSource,MatPaginator} from '@angular/material';
import { NodesService } from '../nodes.service';
import {NODEDATA} from '../types'
@Component({
  selector: 'app-details',
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.css']
})
export class DetailsComponent implements AfterViewInit {
  displayedColumns = ['id', 'latitude', 'longitude'];
  dataSource :any;
  
  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(private nodes: NodesService) {
   
   }


  ngAfterViewInit() {
    //this.api.updateNodeData().then(data=>{
      // this.api.fillBuffer()
      // this.dataSource = new MatTableDataSource<NODEDATA>(Array.from(this.api.nextNodes.values()));
      // this.dataSource.paginator = this.paginator;

   // })
    
  }
  
  
}