import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { HttpService } from './http.service';


@Injectable()
export class SessionService {
  started = false
  selected = 0
  data : any
  globalInfo : any
  xmin = 0
  ymin = 0
  xmax = 500
  ymax = 500
  ww = 500
  hh = 500
  height = 500
  width = 500
  onRedraw = new Subject();
  onStart = new Subject();
  onStop = new Subject();


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

}
