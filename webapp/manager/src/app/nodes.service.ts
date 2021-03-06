import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http'
import {NODEDATA} from './types'

@Injectable()
export class NodesService {

  constructor(private http:HttpClient) { }

  public nodesBuffer: Map<number,NODEDATA>[] = [new Map<number,NODEDATA>(),new Map<number,NODEDATA>(),new Map<number,NODEDATA>()]
  public savedFrame:Map<number,NODEDATA>=new Map<number,NODEDATA>();

  updateNodeData(numTicks:number): Promise<any> {
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
      if(numTicks>0){
        this.nodesBuffer[1].forEach((node,key,nodes)=>{
          node.dlat=(this.nodesBuffer[0].get(key).lat -node.lat)/(numTicks)
          node.dlong=(this.nodesBuffer[0].get(key).long -node.long)/(numTicks)
          
          if(!node.sn && (Math.abs(node.dlat*numTicks) > 90 || Math.abs(node.dlong*numTicks)>180)){//in case node wraps around, don't animate it
            node.dlat=0
            node.dlong=0
          }
        });
    }else{
      this.nodesBuffer[0].forEach((node,key,nodes)=>{
        node.latcp=node.lat;
        node.longcp=node.long;
      })
    }
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
   SavedNodes():Map<number,NODEDATA> {
    return this.savedFrame;
  }

   recalc(numTicks:number){
    this.nodesBuffer[1].forEach((node,key,nodes)=>{
      node.dlat=(this.nodesBuffer[0].get(key).lat -node.lat)/(numTicks)
      node.dlong=(this.nodesBuffer[0].get(key).long -node.long)/(numTicks)
      
      if(!node.sn && (Math.abs(node.dlat*numTicks) > 90 || Math.abs(node.dlong*numTicks)>180)){//in case node wraps around, don't animate it
        node.dlat=0
        node.dlong=0
      }
      
    });

    this.nodesBuffer[2].forEach((node,key,nodes)=>{
      node.dlat=(this.nodesBuffer[1].get(key).lat -node.lat)/(numTicks)
      node.dlong=(this.nodesBuffer[1].get(key).long -node.long)/(numTicks)
      
      if(!node.sn && (Math.abs(node.dlat*numTicks) > 90 || Math.abs(node.dlong*numTicks)>180)){//in case node wraps around, don't animate it
        node.dlat=0
        node.dlong=0
      }
    });

   }

   saveCurrentFrame(){
     this.savedFrame=new Map<number,NODEDATA>(this.nodesBuffer[2])
     this.savedFrame.clear();
     this.nodesBuffer[2].forEach((node,key,nodes) => {
       var copy=new NODEDATA();
       copy.id =node.id;
       copy.lat=node.lat;
       copy.long=node.long;
       copy.sn=node.sn;
       this.savedFrame.set(key,copy);
     });
   }

   reset(){
     this.nodesBuffer = [new Map<number,NODEDATA>(),new Map<number,NODEDATA>(),new Map<number,NODEDATA>()]
     this.savedFrame=new Map<number,NODEDATA>();
   }
}
