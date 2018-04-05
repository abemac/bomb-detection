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
  displayedColumns = ['id', 'latitude', 'longitude','battery','value'];
  dataSource :any;
  filterStr:string;

  @ViewChild(MatPaginator) paginator: MatPaginator;

  constructor(private nodes: NodesService) {
   
   }


  ngAfterViewInit() {
    
    
  }
  applyFilter(){
    this.dataSource.filter=this.filterStr;
  }
  refresh(){
    this.nodes.updateNodeData(-1).then(data=>{
      this.dataSource = new MatTableDataSource<NODEDATA>(Array.from(this.nodes.nodesBuffer[0].values()).sort((a,b)=>{
        return a.id - b.id;
      }));
      this.dataSource.paginator = this.paginator;
      this.dataSource.filter=this.filterStr;
      this.dataSource.filterPredicate = (data)=>{
        return String(data.id).startsWith(this.filterStr);
      };
   })
  }


  
  
}