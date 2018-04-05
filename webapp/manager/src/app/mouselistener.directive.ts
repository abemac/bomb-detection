import { Directive, Output, HostListener, EventEmitter } from '@angular/core';

@Directive({ selector: '[mouseListener]' })

export class MouseListenerDirective {
  @Output() mouseWheelUp = new EventEmitter();
  @Output() mouseWheelDown = new EventEmitter();
  @Output() mouseDown = new EventEmitter()
  @Output() mouseUp = new EventEmitter()
  @Output() mouseMove = new EventEmitter()
  @Output() mouseOut = new EventEmitter()
  @Output() mouseIn = new EventEmitter()
  @Output() mouseClick = new EventEmitter()

  @HostListener('mousewheel', ['$event']) onMouseWheelChrome(event: any) {
    this.mouseWheelFunc(event);
  }

  @HostListener('DOMMouseScroll', ['$event']) onMouseWheelFirefox(event: any) {
    this.mouseWheelFunc(event);
  }

  @HostListener('onmousewheel', ['$event']) onMouseWheelIE(event: any) {
    this.mouseWheelFunc(event);
  }

  @HostListener('mousedown',['$event']) onMouseDown(event:any){
    this.mouseDownFunc(event)
  }
  @HostListener('mouseup',['$event']) onMouseUp(event:any){
    this.mouseUpFunc(event)
  }
  @HostListener('mousemove',['$event']) onMouseMove(event:any){
    this.mouseMoveFunc(event)
  }
  @HostListener('mouseout',['$event']) onMouseOut(event:any){
    this.mouseOutFunc(event)
  }
  @HostListener('mouseenter',['$event']) onMouseIn(event:any){
    this.mouseInFunc(event)
  }
  @HostListener('click',['$event']) onMouseClick(event:any){
    this.onMouseClickFunc(event)
  }
  onMouseClickFunc(event:any){
    this.mouseClick.emit(event)
  }
  
  mouseWheelFunc(event: any) {
    var event = window.event || event; // old IE support
    var delta = Math.max(-1, Math.min(1, (event.wheelDelta || -event.detail)));
    if(delta > 0) {
        this.mouseWheelUp.emit(event);
    } else if(delta < 0) {
        this.mouseWheelDown.emit(event);
    }
    // for IE
    event.returnValue = false;
    // for Chrome and Firefox
    if(event.preventDefault) {
        event.preventDefault();
    }
  }

  mouseDownFunc(event:any){
    this.mouseDown.emit(event)
  } 
  mouseUpFunc(event:any){
    this.mouseUp.emit(event)
  } 
  mouseMoveFunc(event:any){
    this.mouseMove.emit(event)
  }
  mouseOutFunc(event:any){
    this.mouseOut.emit(event)
  }
  mouseInFunc(event:any){
    this.mouseIn.emit(event)
  }
}
