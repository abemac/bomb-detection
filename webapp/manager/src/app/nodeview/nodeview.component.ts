import { Component, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { NODEDATA } from '../types'
import { NodesService } from '../nodes.service';
import { MatDialog, MatDialogRef, MAT_DIALOG_DATA, } from '@angular/material';
import { SimchooserComponent } from '../simchooser/simchooser.component'
import { HttpClient,HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-nodeview',
  templateUrl: './nodeview.component.html',
  styleUrls: ['./nodeview.component.css']
})
export class NodeViewComponent implements AfterViewInit {


  blockSizePixels: number = 25
  purpleIntensityPerNode = 100
  smoothPurple: boolean = true
  smoothingTime: number = 500
  drawgrid: boolean = true;
  drawnodes: boolean = true;
  drawsupernodes: boolean = true;
  nodeSizePixels: number = 4
  supernodeSizePixels: number = 6
  public fps: number = 25
  blockIntensites: number[][]
  public updateInterval = 7000
  interpolate: boolean = true


  autorefresh: boolean;
  paused: boolean = false;
  loading: boolean = true;
  refreshThreadHandle: any;
  numNodes: number = 0;
  numSuperNodes: number = 0
  public fps_actual: string = "0";
  interpolation_step = 0
  interpolation_handle: any;
  lastTime: number;

  @ViewChild('canvas') c: ElementRef;
  context: CanvasRenderingContext2D;
  canvas: HTMLCanvasElement;
  canvasWidthPixels: number = window.innerWidth / 2 - 100
  canvasHeightPixels: number = window.innerHeight - 250
  scale: number = 1
  oldscale: number = 1
  focusy: number;
  focusx: number;
  trx: number = this.canvasWidthPixels / 2
  try: number = this.canvasHeightPixels / 2

  constructor(public nodes: NodesService, public dialog: MatDialog,private http:HttpClient) {
    this.focusy = this.canvasHeightPixels / 2
    this.focusx = this.canvasWidthPixels / 2
    this.blockIntensites = new Array<Array<number>>()
    for (var r = 0; r * this.blockSizePixels < this.canvasHeightPixels; r++) {
      var row: number[] = new Array<number>()
      for (var c = 0; c * this.blockSizePixels < this.canvasWidthPixels; c++) {
        row.push(0)
      }
      this.blockIntensites.push(row)
    }
    this.nodes.updateNodeData(-1);
  }

  ngAfterViewInit() {
    this.canvas = (this.c.nativeElement as HTMLCanvasElement)
    this.context = this.canvas.getContext('2d');
    this.update()
    this.onStartToggle(null)
    setTimeout(() => this.loading = false, 2 * this.updateInterval)
  }
  openDialog(): void {

    let dialogRef = this.dialog.open(SimchooserComponent, {
      width: '1000px',
      disableClose: true
      
      //data: { name: this.name, animal: this.animal }
    });

    dialogRef.afterClosed().subscribe(result => {
      if(result=="SUCCESS"){

      }
    });
  }

  
  stopSim(){
    var httpOptions = {
      responseType: 'text' as 'text'
    };
    this.http.get('/StopSim',httpOptions).toPromise().then( resp =>{
      console.log(resp)
      this.nodes.reset()
      this.updateView()
    }).catch(err=>{
      console.log(err.error)
    })

  }

  update() {
    clearInterval(this.interpolation_handle)
    var time = performance.now()
    if (!isNaN(this.interpolation_step / (time - this.lastTime) * 1000)) {
      this.fps_actual = (this.interpolation_step / (time - this.lastTime) * 1000).toFixed(3)
    }
    this.nodes.shiftBuffer()
    this.interpolation_step = 0
    if (this.interpolate) {
      this.lastTime = performance.now()
      this.interpolation_handle = setInterval(() => this.incrementInterpolationStep(), 1000 / this.fps)
    } else {
      this.updateView()
    }
    this.nodes.updateNodeData(this.updateInterval * this.fps / 1000);

  }

  updateView() {
    this.context.clearRect(0, 0, this.canvasWidthPixels, this.canvasHeightPixels)
    var ratio = this.scale / this.oldscale
    this.trx = this.focusx + (this.trx - this.focusx) * ratio
    this.try = this.focusy + (this.try - this.focusy) * ratio
    //ALT
    //this.trx=this.trx*ratio+this.focusx*(1-ratio)
    //this.try=this.try*ratio+this.focusy*(1-ratio)
    if (this.drawgrid) {
      this.drawGrid()
      this.updateBlockCounts()
      this.colorSections()
    } else {
      this.context.fillStyle = "ghostwhite"
      this.context.fillRect(0, 0, this.canvasWidthPixels, this.canvasHeightPixels)
    }
    this.drawNodes()
    this.oldscale = this.scale
  }
  updateBlockCounts() {
    var old
    if (this.smoothPurple && this.autorefresh && this.interpolate) {
      old = JSON.parse(JSON.stringify(this.blockIntensites))
    }
    for (var r = 0; r * this.blockSizePixels < this.canvasHeightPixels; r++) {
      for (var c = 0; c * this.blockSizePixels < this.canvasWidthPixels; c++) {

        this.blockIntensites[r][c] = 0
      }
    }
    if (this.paused) {
      this.nodes.SavedNodes().forEach((node, key, nodes) => {
        var rowi = Math.floor((this.try - node.lat * this.scale) / this.blockSizePixels)
        var coli = Math.floor((this.trx + node.long * this.scale) / this.blockSizePixels)

        if (this.blockIntensites[rowi] != undefined) {
          if (this.blockIntensites[rowi][coli] != undefined) {
            this.blockIntensites[rowi][coli] += this.purpleIntensityPerNode / 1000.0
          }
        }
      });
    } else {
      this.nodes.Nodes().forEach((node, key, nodes) => {
        var rowi = Math.floor((this.try - node.lat * this.scale) / this.blockSizePixels)
        var coli = Math.floor((this.trx + node.long * this.scale) / this.blockSizePixels)

        if (this.blockIntensites[rowi] != undefined) {
          if (this.blockIntensites[rowi][coli] != undefined) {
            this.blockIntensites[rowi][coli] += this.purpleIntensityPerNode / 1000.0
          }
        }
      });
    }
    if (this.smoothPurple && this.autorefresh && this.interpolate) {
      for (var r = 0; r * this.blockSizePixels < this.canvasHeightPixels; r++) {
        for (var c = 0; c * this.blockSizePixels < this.canvasWidthPixels; c++) {
          this.blockIntensites[r][c] = old[r][c] + (this.blockIntensites[r][c] - old[r][c]) / (this.smoothingTime * this.fps / 1000)
        }
      }
    }
  }
  incrementInterpolationStep() {
    this.updateView()
    this.nodes.Nodes().forEach((node, key, nodes) => {
      node.lat += node.dlat;
      node.long += node.dlong;
    })
    this.interpolation_step++;
    if (this.interpolation_step > this.updateInterval * this.fps / 1000) {
      clearInterval(this.interpolation_handle)
    }
  }
  showLoadingSpinner(event) {
    this.loading = true;
    setTimeout(() => this.loading = false, this.updateInterval)
  }
  onSliderChange(event) {
    this.blockSizePixels = event.value;
    this.blockIntensites = new Array<Array<number>>()
    for (var r = 0; r * this.blockSizePixels < this.canvasHeightPixels; r++) {
      var row: number[] = new Array<number>()
      for (var c = 0; c * this.blockSizePixels < this.canvasWidthPixels; c++) {
        row.push(0)
      }
      this.blockIntensites.push(row)
    }
    var oldSmooth = this.smoothPurple
    this.smoothPurple = false
    this.updateView()
    this.smoothPurple = oldSmooth;

  }
  onFPSchange(event) {
    this.showLoadingSpinner(event);
    this.nodes.recalc(this.updateInterval * event.value / 1000);

  }
  onNodeSliderChange(event) {
    this.nodeSizePixels = event.value
    this.updateView()
  }
  onSuperNodeSliderChange(event) {
    this.supernodeSizePixels = event.value
    this.updateView()
  }
  onPurpleChange(event) {
    this.purpleIntensityPerNode = event.value;
    this.updateView()
  }
  onStartToggle(event) {
    if (!this.autorefresh) {
      this.autorefresh = true;
      this.refreshThreadHandle = setInterval(() => { this.update(); }, this.updateInterval);
    } else {
      this.autorefresh = false;
      clearInterval(this.refreshThreadHandle);
    }
  }
  onPause(event) {
    this.paused = !this.paused;
    if (this.paused) {
      this.nodes.saveCurrentFrame();
    }
  }

  drawGrid() {
    //draw rows
    for (var i = 0; i < this.canvas.height; i += this.blockSizePixels) {
      this.context.beginPath();
      this.context.moveTo(0, i);
      this.context.lineTo(this.canvas.width, i);
      this.context.stroke();
    }
    //draw columns
    for (var i = 0; i < this.canvas.width; i += this.blockSizePixels) {
      this.context.beginPath();
      this.context.moveTo(i, 0);
      this.context.lineTo(i, this.canvas.height);
      this.context.stroke();
    }
  }
  drawNodes() {
    this.context.save()
    this.context.translate(this.trx, this.try)
    if (this.paused) {
      this.nodes.SavedNodes().forEach((node, key, nodes) => {
        this.drawNode(node.long * this.scale, -node.lat * this.scale, node.sn)
      });
    } else {
      this.nodes.Nodes().forEach((node, key, nodes) => {
        this.drawNode(node.long * this.scale, -node.lat * this.scale, node.sn)
      });
    }
    this.context.restore()
  }
  drawNode(x: number, y: number, supernode: boolean) {

    if (supernode && this.drawsupernodes) {
      this.context.fillStyle = "#ff4081"
      this.context.strokeStyle = "#ff4081"
      this.context.beginPath()
      this.context.arc(x, y, (this.supernodeSizePixels), 0, 2 * Math.PI)
      this.context.fill()
    } else if (!supernode && this.drawnodes) {
      this.context.fillStyle = "hsl(120, 100%, 50%)"
      this.context.strokeStyle = "hsl(120, 100%, 50%)"
      this.context.beginPath()
      this.context.arc(x, y, this.nodeSizePixels, 0, 2 * Math.PI)
      this.context.fill()
    }
  }

  colorSections() {
    var intensity = 0
    for (var r = 0; r * this.blockSizePixels < this.canvasHeightPixels; r++) {
      for (var c = 0; c * this.blockSizePixels < this.canvasWidthPixels; c++) {
        intensity = Math.min(this.blockIntensites[r][c], 1)
        this.context.fillStyle = "hsla(260,100%,43%," + intensity.toString() + ")"
        this.context.fillRect(c * this.blockSizePixels + 1, r * this.blockSizePixels + 1, this.blockSizePixels - 2, this.blockSizePixels - 2)
      }
    }
  }



  onMouseWheelUp(event) {
    this.scale *= 1.2
    this.focusx = event.offsetX
    this.focusy = event.offsetY
    this.updateView()
  }
  onMouseWheelDown(event) {
    if (this.scale > 1) {
      this.scale /= 1.2
      this.focusx = event.offsetX
      this.focusy = event.offsetY
      this.updateView()
    } else {
      this.scale = 1
    }

  }
  mouseDown: boolean = false
  handPtr: string = "default-cursor"
  startX: number
  startY: number

  onMouseDown(event) {
    this.mouseDown = true
    this.startX = event.offsetX
    this.startY = event.offsetY
  }
  onMouseUp(event) {
    this.mouseDown = false
    this.handPtr = "default-cursor"
  }
  onMouseMove(event) {
    if (this.mouseDown) {
      this.handPtr = "grab-cursor"
      var dx = event.offsetX - this.startX
      var dy = event.offsetY - this.startY
      if (Math.abs(dx) > 2 || Math.abs(dy) > 2) {
        this.trx += dx
        this.try += dy
        this.updateView()
        this.startX = event.offsetX
        this.startY = event.offsetY
      }
    }
    else{
      if (this.paused) {
        this.nodes.SavedNodes().forEach((node, key, nodes) => {
          if (this.mouseOver(event,node)){
            // console.log(node)
          }
        });
      } else {
        
      }
    }
  }
  mouseOver(event,node):boolean{
    var margin=this.nodeSizePixels;
    if(node.sn){
      var margin=this.supernodeSizePixels;
    }
    //console.log(event,node,node.long * this.scale - event.offsetX+this.trx,this.try-node.lat * this.scale -event.offsetY)
    if((Math.abs(node.long * this.scale - event.offsetX+this.trx) < margin) && (Math.abs(this.try-node.lat * this.scale -event.offsetY) < margin)){
        return true;
    }

    return false;
  }
  onMouseOut(event) {
    this.mouseDown = false
    this.handPtr = "default-cursor"
  }
  onMouseIn(event) {
    if (event.button == 0 && event.buttons == 1) {
      this.mouseDown = true
      this.handPtr = "grab-cursor"
    }

  }

}
