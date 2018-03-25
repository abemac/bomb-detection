import { Injectable,Output,EventEmitter } from '@angular/core';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';

@Injectable()
export class TabControlService {
  @Output() tabindex: EventEmitter<number> = new EventEmitter();
  constructor() { }
  setIndex(index: number ){
    this.tabindex.emit(index)
  }

}
