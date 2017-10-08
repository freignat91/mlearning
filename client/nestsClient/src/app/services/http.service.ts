import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Subject } from 'rxjs/Subject';
//import { Observable } from 'RxJS/Rx';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/retryWhen';
import 'rxjs/add/operator/scan';
import 'rxjs/add/operator/delay';



const httpRetryDelay = 200
const httpRetryNumber = 3

@Injectable()
export class HttpService {
  private token = ""
  onHttpError = new Subject();
  //addr = "http://localhost:3001/api/v1"
  addr : string


  constructor(private http : Http) {
    this.addr = "http://"+window.location.host+"/api/v1"
    console.log(this.addr)
  }

  getData() {
    return this.httpGet("/data")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  start() {
    return this.httpGet("/start")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  stop() {
    return this.httpGet("/stop")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  isStarted() {
    return this.httpGet("/isStarted")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  nextTime() {
    return this.httpGet("/nextTime")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  exportSample() {
    return this.httpGet("/exportAntSample")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  setSleep(value) {
    return this.httpGet("/setSleep/"+value)
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  setSelected(value) {
    return this.httpGet("/setSelected/"+value)
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  getGlobalInfo() {
    return this.httpGet("/globalInfo")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  getInfo() {
    return this.httpGet("/info")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  restart() {
    return this.httpGet("/restart")
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  addFoods(x : number, y : number) {
    return this.httpPost("/addFoods", {x: x, y: y})
      .map((res : Response) => {
        return res.json()
      }
    )
  }

  foodRenew(renew : boolean) {
    return this.httpPost("/foodRenew", { ret: renew })
      .map((res : Response) => {
        return res.json()
      }
    )
  }
//--------------------------------------------------------------------------------------
// http core functions
//--------------------------------------------------------------------------------------

  private formatedDate(daten : number) : string {
    let date = new Date(daten * 1000)
    let num = ""+date.getDate()
    if (date.getDate()<10) {
      num='0'+num
    }
    let month = ""+(date.getMonth()+1)
    if (date.getMonth()+1<10) {
      month = '0'+month
    }
    return date.getFullYear()  + "-" +
    month + "-" +
    num + " " +
    date.getHours() + ":" +
    date.getMinutes();
  }

  private setHeaders() {
    var headers = new Headers
    headers.set('Authorization', this.token)
    return headers
  }

  setToken(token : string) {
    this.token = token
  }

  httpGet(url : string) : Observable<any> {
    return this.http.get(this.addr+url)
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpDelete(url : string) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.delete(this.addr+url, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpPost(url : string, data : any) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.post(this.addr+url, data, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

  httpPut(url : string, data : any) : Observable<any> {
    let headers = this.setHeaders()
    return this.http.put(this.addr+url, data, { headers: this.setHeaders() })
      .retryWhen(e => e.scan<number>((errorCount, err) => {
        console.log("retry: "+(errorCount+1))
        if (errorCount >= httpRetryNumber-1) {
            throw err;
        }
        return errorCount + 1;
      }, 0).delay(httpRetryDelay)
    )
  }

}
