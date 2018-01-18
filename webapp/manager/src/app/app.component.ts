import { Component,ViewChild,ElementRef,AfterViewInit } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'app';
    @ViewChild('canvas') c :ElementRef;
    private context: CanvasRenderingContext2D;
    private canvas: HTMLCanvasElement;
    ngAfterViewInit() {
      this.canvas=(this.c.nativeElement as HTMLCanvasElement)
      this.context = this.canvas.getContext('2d');
      this.drawGrid()
      this.drawNode(50,50)
      this.colorSection(3,3,20)
      this.colorSection(5,2,40)
      
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

    colorSection(r:number,c:number,intensity:number){
      this.context.fillStyle="hsla(260,100%,43%,"+intensity+"%)"
      this.context.fillRect(r*50+1,c*50+1,48,48)

    }
    getRandomColor() {
      //see https://martin.ankerl.com/2009/12/09/how-to-create-random-colors-programmatically/
    }
}

