import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { HttpService } from './http.service';


@Injectable()
export class SessionService {
  started = false
  nestSelected = 0
  selected = 0
  data : any
  globalInfo : any
  info : any
  panelHeight = 500
  panelWidth = 800
  fx0 = 0
  fy0 = 0
  fwidth = 800
  fheight = 500
  xmin = 0
  ymin = 0
  xmax = 800
  ymax = 500
  height = 500
  width = 800
  fpanx = (this.xmax + this.xmin) / 2
  fpany = (this.ymax + this.ymin) / 2
  zoom = 1
  coefZoom = 1
  onRedraw = new Subject();
  onStart = new Subject();
  onStop = new Subject();
  onClear = new Subject();
  mode = "select"
  displayFight = false
  foodRenew = true
  panicMode = true
  display = true
  nestsInfo = []
  nestColors = ["blue", "red", "black", "green"]


  constructor(private httpService : HttpService) {
  }

  public redraw() {
    this.onRedraw.next()
  }

  public start() {
    this.onStart.next()
  }

  public stop() {
    this.onStop.next()
  }

  public clear() {
    this.onClear.next()
  }

  public setZoom(val : number) {
    this.zoom = this.zoom * val
    this.fwidth = this.width * this.zoom
    this.fheight = this.fwidth * this.panelHeight / this.panelWidth
    this.coefZoom = this.panelWidth / this.fwidth
    this.fx0 = this.fpanx - this.fwidth/2*this.coefZoom
    this.fy0 = this.fpany - this.fheight/2*this.coefZoom
    this.redraw()
  }

  public getInvx(x : number) {
    return (x - this.panelWidth/2)/this.coefZoom+this.fx0
  }

  public getInvy(y : number) {
    return (y - this.panelHeight/2)/this.coefZoom+this.fy0
  }

  public getx(x : number) : number {
    return (x-this.fx0)*this.coefZoom+this.panelWidth/2
  }

  public gety(y : number) : number {
    return (y-this.fy0)*this.coefZoom+this.panelHeight/2
  }

  public getl(l : number) : number {
    return l*this.coefZoom
  }

  getInvl(l : number) : number {
    return l/this.coefZoom
  }
}
