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
  xmin = 0
  ymin = 0
  xmax = 800
  ymax = 500
  height = 500
  width = 800
  onRedraw = new Subject();
  onStart = new Subject();
  onStop = new Subject();
  mode = "select"
  displayContact = false
  displayFight = false
  foodRenew = true
  panicMode = true


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
