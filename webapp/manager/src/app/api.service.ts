import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http'
import {NODEDATA} from './types'

@Injectable()
export class ApiService {

  constructor(private http:HttpClient) { }

  public nodesBuffer: Map<number,NODEDATA>[] = [new Map<number,NODEDATA>(),new Map<number,NODEDATA>(),new Map<number,NODEDATA>()]
  
  updateNodeData(): Promise<any> {
    return this.http.get('/GetNodes').toPromise<any>().then( resp =>{
      this.nodesBuffer[0]= new Map<number,NODEDATA>()
      for (let node of resp['nodes']){
            this.nodesBuffer[0].set(node.id,node)
      }
      if(this.nodesBuffer[2].size==0){
        this.nodesBuffer[2]=new Map<number,NODEDATA>(this.nodesBuffer[0])
      }
      if(this.nodesBuffer[1].size==0){
        this.nodesBuffer[1]=new Map<number,NODEDATA>(this.nodesBuffer[0])
      }
      this.nodesBuffer[1].forEach((node,key,nodes)=>{
        node.dlat=(this.nodesBuffer[0].get(key).lat -node.lat)/(25.0)
        node.dlong=(this.nodesBuffer[0].get(key).long -node.long)/(25.0)
        
        if(Math.abs(node.dlat*25) > 90 || Math.abs(node.dlong*25)>180){//in case node wraps around, don't animate it
          node.dlat=0
          node.dlong=0
        }
      })
     }).catch(err=>{
       return Promise.reject(err.message || err);
     })
   }

   shiftBuffer(){
    this.nodesBuffer[2]=this.nodesBuffer[1]
    this.nodesBuffer[1]=this.nodesBuffer[0]
   }

   Nodes():Map<number,NODEDATA> {
     return this.nodesBuffer[2]
   }
}
