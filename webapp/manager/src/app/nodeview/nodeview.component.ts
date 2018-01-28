import { Component,ViewChild,ElementRef,AfterViewInit  } from '@angular/core';
import {NODEDATA} from '../types'
import { ApiService } from '../api.service';
@Component({
  selector: 'app-nodeview',
  templateUrl: './nodeview.component.html',
  styleUrls: ['./nodeview.component.css']
})
export class NodeViewComponent implements AfterViewInit {
  
  @ViewChild('canvas') c :ElementRef;
  private context: CanvasRenderingContext2D;
  private canvas: HTMLCanvasElement;
  private canvasWidthPixels: number =window.innerWidth/2-100
  private canvasHeightPixels:number = window.innerHeight-250
  private blockSizePixels:number=25
  private purpleIntensityPerNode=.1
  private blockIntensites: number[][]

  private nodeSizePixels :number=4
  private supernodeSizePixels :number=6
  private scale : number=1
  private oldscale : number=1
  private focusy : number;
  private focusx : number;
  private trx : number=this.canvasWidthPixels/2
  private try: number=this.canvasHeightPixels/2

  private autorefresh:boolean;
  protected refreshThreadHandle : any;
  private btnColor: string="green"
  private drawgrid: boolean=true;
  private drawnodes: boolean=true;
  private drawsupernodes: boolean=true;

  public interpolation_ms=40
  public updateInterval=7000
  private interpolation_step=0
  private interpolation_handle :any;
  private interpolate: boolean=true

  constructor(private api: ApiService){
    this.focusy=this.canvasHeightPixels/2
    this.focusx=this.canvasWidthPixels/2
    this.blockIntensites=new Array<Array<number>>()
    for(var r=0;r*this.blockSizePixels<this.canvasHeightPixels;r++){
      var row :number[]=new Array<number>()
      for(var c=0;c*this.blockSizePixels<this.canvasWidthPixels;c++){
        row.push(0)
      }
      this.blockIntensites.push(row)
    }
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
  
  drawNode(x:number,y:number,supernode:boolean){
    if(supernode && this.drawsupernodes){
      this.context.fillStyle="#ff4081"
      this.context.strokeStyle="#ff4081"
      this.context.beginPath()
      this.context.arc(x,y,(this.supernodeSizePixels)/this.scale,0,2*Math.PI)
     this.context.fill()
    }else if (!supernode && this.drawnodes){
      this.context.fillStyle="hsl(120, 100%, 50%)"
      this.context.strokeStyle="hsl(120, 100%, 50%)"
      this.context.beginPath()
      this.context.arc(x,y,this.nodeSizePixels/this.scale,0,2*Math.PI)
      this.context.fill()
    }
    
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
    this.context.translate(this.trx,this.try)
    this.context.scale(this.scale,this.scale)
    this.api.Nodes().forEach((node,key,nodes)=>{
      this.drawNode(node.long,node.lat,node.sn)
    })
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
    this.api.Nodes().forEach((node,key,nodes)=>{
      var rowi=Math.floor((this.try+node.lat*this.scale)/this.blockSizePixels)
      var coli=Math.floor((this.trx+node.long*this.scale)/this.blockSizePixels)
      
      if(this.blockIntensites[rowi] != undefined){
        if(this.blockIntensites[rowi][coli]!=undefined){
          this.blockIntensites[rowi][coli]+=this.purpleIntensityPerNode
        }
      }
    })
    
  }

  update(){
    clearInterval(this.interpolation_handle)
    this.api.shiftBuffer()
    this.interpolation_step=0
    if(this.interpolate){
      this.interpolation_handle=setInterval(()=>this.incrementInterpolationStep(),this.interpolation_ms)
    }else{
      this.updateView()
    }
    this.api.updateNodeData(this.updateInterval / this.interpolation_ms);

  }
  incrementInterpolationStep(){
    this.updateView()
    this.api.Nodes().forEach((node,key,nodes)=>{
        node.lat+=node.dlat;
        node.long+=node.dlong;
    })
    this.interpolation_step++;
    if(this.interpolation_step > (this.updateInterval / this.interpolation_ms-1) ){
      clearInterval(this.interpolation_handle)
      this.interpolation_step=0
    }
  }
  onStartToggle(event){
    if(!this.autorefresh){
      this.autorefresh=true;
      this.btnColor="red"
      this.refreshThreadHandle = setInterval(() => {this.update();},this.updateInterval);
    }else{
      this.autorefresh=false;
      this.btnColor="green"
      clearInterval(this.refreshThreadHandle);
    }
  }
  toggleGrid(){
    
  }
  updateView(){
    this.context.clearRect(0,0,this.canvasWidthPixels,this.canvasHeightPixels)
    var ratio=this.scale/this.oldscale
    this.trx=this.focusx+(this.trx-this.focusx)*ratio
    this.try=this.focusy+(this.try-this.focusy)*ratio
    //ALT
    //this.trx=this.trx*ratio+this.focusx*(1-ratio)
    //this.try=this.try*ratio+this.focusy*(1-ratio)
    if(this.drawgrid){
      this.drawGrid()
      this.updateBlockCounts()
      this.colorSections()
    }else{
      this.context.fillStyle="ghostwhite"
      this.context.fillRect(0,0,this.canvasWidthPixels,this.canvasHeightPixels)
    }

    this.drawNodes()
    
  
    
    this.oldscale=this.scale
    
    
    
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
  onNodeSliderChange(event){
    this.nodeSizePixels=event.value
    this.updateView()
  }
  onSuperNodeSliderChange(event){
    this.supernodeSizePixels=event.value
    this.updateView()
  }
  onMouseWheelUp(event){
    this.scale*=1.2
    this.focusx=event.offsetX
    this.focusy=event.offsetY
    this.updateView()
  }
  onMouseWheelDown(event){
    if(this.scale>1){ 
      this.scale/=1.2
      this.focusx=event.offsetX
      this.focusy=event.offsetY
      this.updateView()
    }else{
      this.scale=1
    }
   
  }

  private mouseDown :boolean=false
  private handPtr : string="default-cursor"
  private startX: number
  private startY: number

  onMouseDown(event){
    this.mouseDown=true
    this.handPtr="grab-cursor"
    this.startX=event.offsetX
    this.startY=event.offsetY
  }
  onMouseUp(event){
    this.mouseDown=false
    this.handPtr="default-cursor"
  }
  onMouseMove(event){
    if(this.mouseDown){
      var dx = event.offsetX - this.startX
      var dy = event.offsetY - this.startY
      if(Math.abs(dx)>2 || Math.abs(dy) >2){
        this.trx+=dx
        this.try+=dy
        this.updateView()
        this.startX=event.offsetX
        this.startY=event.offsetY
      }
      
    }
  }
  onMouseOut(event){
    this.mouseDown=false
    this.handPtr="default-cursor"
  }
  onMouseIn(event){
    if(event.button==0 && event.buttons==1){
      this.mouseDown=true
      this.handPtr="grab-cursor"
    }
    
  }

}
