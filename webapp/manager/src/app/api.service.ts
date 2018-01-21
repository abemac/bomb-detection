import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http'
import {NODEDATA} from './types'
@Injectable()
export class ApiService {

  constructor(private http:HttpClient) { }


  getNodes():Promise<NODEDATA[]>{
    return this.http.get('http://localhost:8080/GetNodes').toPromise().then( resp =>{
       return resp['nodes'] as NODEDATA[]
     }).catch(err=>{
       return Promise.reject(err.message || err);
     })
   }

}
