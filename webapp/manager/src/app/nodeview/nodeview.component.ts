import { Component, OnInit,ViewChild,ElementRef,AfterViewInit  } from '@angular/core';
import {NODEDATA} from '../types'
import { ApiService } from '../api.service';
@Component({
  selector: 'app-nodeview',
  templateUrl: './nodeview.component.html',
  styleUrls: ['./nodeview.component.css']
})
export class NodeViewComponent implements OnInit {
  
  @ViewChild('canvas') c :ElementRef;
  private context: CanvasRenderingContext2D;
  private canvas: HTMLCanvasElement;
  private canvasWidth: number = 400
  private canvasHeight:number = 400
  private nodeData : NODEDATA[]
  private blockCounts: number[][]

  constructor(private api: ApiService){
    
  }

  ngOnInit() {
    this.blockCounts=new Array<Array<number>>()
    for(var r=0;r*50<this.canvasHeight;r++){
      var row :number[]=new Array<number>()
      for(var c=0;c*50<this.canvasWidth;c++){
        row.push(0)
      }
      this.blockCounts.push(row)
    }
  }
  
  ngAfterViewInit() {
    this.canvas=(this.c.nativeElement as HTMLCanvasElement)
    this.context = this.canvas.getContext('2d');
    this.drawGrid()
    this.update()
  }

  drawGrid() {
    //draw rows
    for(var i=0;i<this.canvas.height;i+=50){
      this.context.beginPath();
      this.context.moveTo(0,i);
      this.context.lineTo(this.canvas.width,i);
      this.context.stroke();
    }
    //draw columns
    for(var i=0;i<this.canvas.width;i+=50){
      this.context.beginPath();
      this.context.moveTo(i,0);
      this.context.lineTo(i,this.canvas.height);
      this.context.stroke();
    }
  }
  
  drawNode(x:number,y:number){
    this.context.fillStyle="hsl(120, 100%, 50%)"
    this.context.strokeStyle="hsl(120, 100%, 50%)"
    this.context.beginPath()
    this.context.arc(x,y,4,0,2*Math.PI)
    this.context.fill()
  }

  colorSections(){
    var intensity=0
    for(var r=0;r*50<this.canvasHeight;r++){
      for(var c=0;c*50<this.canvasWidth;c++){
        intensity=Math.min(this.blockCounts[r][c],1)
        this.context.fillStyle="hsla(260,100%,43%,"+intensity.toString()+")"
        this.context.fillRect(c*50+1,r*50+1,48,48)
      }
    }  
  }

  drawNodes(){
    this.context.save()
    this.context.translate(this.canvasWidth/2,this.canvasHeight/2)
    for(var i=0;i<this.nodeData.length;i++){
      this.drawNode(this.nodeData[i].longitude,this.nodeData[i].latitude)
    }
    this.context.restore()
  }

  getRandomColor() {
    //see https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/
  }

  updateBlockCounts(){
    for(var r=0;r*50<this.canvasHeight;r++){
      for(var c=0;c*50<this.canvasWidth;c++){
        this.blockCounts[r][c]=0
      }
    }
    for(let node of this.nodeData){
      var row=Math.floor((this.canvasWidth/2+node.latitude)/50)
      var col=Math.floor((this.canvasHeight/2+node.longitude)/50)
      this.blockCounts[row][col]+=.1
    }
  }

  update(){
    this.api.getNodes().then(data=>{
      this.nodeData=data
      this.updateBlockCounts()
      this.colorSections()
      this.drawNodes()
    });
  }

}
