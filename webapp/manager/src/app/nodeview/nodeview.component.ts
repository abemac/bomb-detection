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
  private canvasWidthPixels: number =window.innerWidth/2-100
  private canvasHeightPixels:number = window.innerHeight-250
  private blockSizePixels:number=25
  private purpleIntensityPerNode=.1
  private nodeData : NODEDATA[]
  private blockIntensites: number[][]

  private scale : number=1

  private focuslat : number;
  private focuslong : number;

  constructor(private api: ApiService){
    this.focuslat=this.canvasHeightPixels/2
    this.focuslong=this.canvasWidthPixels/2
    this.blockIntensites=new Array<Array<number>>()
    for(var r=0;r*this.blockSizePixels<this.canvasHeightPixels;r++){
      var row :number[]=new Array<number>()
      for(var c=0;c*this.blockSizePixels<this.canvasWidthPixels;c++){
        row.push(0)
      }
      this.blockIntensites.push(row)
    }
  }

  ngOnInit() {
    
    
  }
  
  ngAfterViewInit() {
    this.canvas=(this.c.nativeElement as HTMLCanvasElement)
    this.context = this.canvas.getContext('2d');
    this.update()
  }

  drawGrid() {
    //draw rows
    for(var i=0;i<this.canvas.height;i+=this.blockSizePixels){
      this.context.beginPath();
      this.context.moveTo(0,i);
      this.context.lineTo(this.canvas.width,i);
      this.context.stroke();
    }
    //draw columns
    for(var i=0;i<this.canvas.width;i+=this.blockSizePixels){
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
    this.context.arc(x,y,4/this.scale,0,2*Math.PI)
    this.context.fill()
  }

  colorSections(){
    var intensity=0
    for(var r=0;r*this.blockSizePixels<this.canvasHeightPixels;r++){
      for(var c=0;c*this.blockSizePixels<this.canvasWidthPixels;c++){
        intensity=Math.min(this.blockIntensites[r][c],1)
        this.context.fillStyle="hsla(260,100%,43%,"+intensity.toString()+")"
        this.context.fillRect(c*this.blockSizePixels+1,r*this.blockSizePixels+1,this.blockSizePixels-2,this.blockSizePixels-2)
      }
    }  
  }

  drawNodes(){
    this.context.save()
    this.context.translate(this.focuslong,this.focuslat)
    this.context.scale(this.scale,this.scale)
    this.context.translate(this.canvasWidthPixels/2-this.focuslong,this.canvasHeightPixels/2-this.focuslat)
    //this.context.translate(this.focuslong,this.focuslat)
    //this.context.translate(this.canvasWidthPixels/2,this.canvasHeightPixels/2)
    
    //this.context.translate((this.canvasWidthPixels/2-this.focuslong)/this.scale,(this.canvasHeightPixels/2-this.focuslat)/this.scale)
    //this.context.tra
  
    
    for(var i=0;i<this.nodeData.length;i++){
      this.drawNode(this.nodeData[i].longitude,this.nodeData[i].latitude)
    }
    this.context.restore()
  }

  getRandomColor() {
    //see https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/
  }

  updateBlockCounts(){
    for(var r=0;r*this.blockSizePixels<this.canvasHeightPixels;r++){
      for(var c=0;c*this.blockSizePixels<this.canvasWidthPixels;c++){
      
        this.blockIntensites[r][c]=0
      }
    }
    for(let node of this.nodeData){
      var rowi=Math.floor((this.focuslat+(node.latitude+this.canvasHeightPixels/2-this.focuslat)*this.scale)/this.blockSizePixels)
      var coli=Math.floor((this.focuslong+(node.longitude+this.canvasWidthPixels/2-this.focuslong)*this.scale)/this.blockSizePixels)
      
      if(this.blockIntensites[rowi] != undefined){
        if(this.blockIntensites[rowi][coli]!=undefined){
          this.blockIntensites[rowi][coli]+=this.purpleIntensityPerNode
        }
      }
    }
    
  }

  update(){
    this.api.getNodes().then(data=>{
      this.nodeData=data
      this.updateView()
    });
  }

  updateView(){
    this.context.clearRect(0,0,this.canvasWidthPixels,this.canvasHeightPixels)
    this.drawGrid()
    
    this.updateBlockCounts()
    this.colorSections()
    this.drawNodes()
    
    
  }
  onSliderChange(event) {
    this.blockSizePixels = event.value;
    this.blockIntensites=new Array<Array<number>>()
    for(var r=0;r*this.blockSizePixels<this.canvasHeightPixels;r++){
      var row :number[]=new Array<number>()
      for(var c=0;c*this.blockSizePixels<this.canvasWidthPixels;c++){
        row.push(0)
      }
      this.blockIntensites.push(row)
    }
    this.updateView()
    
  }
  onMouseWheelUp(event){
    this.scale*=1.2
    this.focuslong=event.offsetX
    this.focuslat=event.offsetY
    this.updateView()
  }
  onMouseWheelDown(event){
    if(this.scale>.4){ 
      this.scale/=1.2
      this.focuslong=event.offsetX
      this.focuslat=event.offsetY
      this.updateView()
    }
   
  }

}
