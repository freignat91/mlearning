
import { ViewChild, Component, Directive, ElementRef, HostListener, Input, Renderer  } from '@angular/core';
import { HttpService } from './services/http.service'
import { SessionService } from './services/session.service'
import { Info } from './models/info.model'

const httpRetryDelay = 200
const httpRetryNumber = 3

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = "title"
  messageError = ""
  isStarted = false
  messageStartStop="start"
  graphPanelHeight = 500
  graphPanelWidth = 500
  speed = 1
  logLevel = 1
  timer : any
  info : any

  constructor(private httpService : HttpService, private sessionService : SessionService) {
    sessionService.onStart.subscribe(
      data => {
        this.start()
      }
    )
    sessionService.onStop.subscribe(
      data => {
        this.stop()
      }
    )
    this.httpService.getGlobalInfo().subscribe(
      data => {
        this.speed = data.waiter
      },
      error => {
        console.log(error)
      }
    )
  }


  startStop() {
    if (this.messageStartStop == "start") {
      this.sessionService.start()
    } else {
      this.sessionService.stop()
    }
  }


  start() {
      this.messageStartStop = "stop"
      console.log("starting")
      this.httpService.start().subscribe(
        data => {
          clearInterval(this.timer)
          this.getInfo()
          this.timer = setInterval(this.getInfo.bind(this), 2000);
          console.log("started")
          this.sessionService.started = true
        },
        error => {
          console.log(error)
        }
      )
    }

   stop() {
     this.messageStartStop = "start"
     console.log("stoping")
     this.httpService.stop().subscribe(
       data => {
         console.log("stopped")
         clearInterval(this.timer)
         this.sessionService.started = false
          this.nextTime()
       },
       error => {
         console.log(error)
       }
     )
   }

   nextTime() {
     this.httpService.nextTime().subscribe(
       data => {
         this.getInfo()
         this.sessionService.redraw()
       },
       error => {
         console.log(error)
       }
     )
   }

   restart() {
     this.sessionService.stop()
     this.httpService.restart().subscribe(
       data => {
        this.sessionService.start()
         console.log("restarted")
       },
       error => {
         console.log(error)
       }
     )
   }

   getInfo() {
     this.httpService.getInfo().subscribe(
       data => {
         //console.log(data)
         this.info = data
       },
       error => {
         console.log(error)
       }
     )
   }

   exportSample() {
     this.httpService.exportSample().subscribe(
       data => {
         console.log(data)
       },
       error => {
         console.log(error)
       }
     )
   }

   setSleep(value) {
     this.speed = this.speed * value
     if (this.speed <0) {
       this.speed = 0
     }
     this.httpService.setSleep(this.speed).subscribe(
       data => {
         console.log(data)
       },
       error => {
         console.log(error)
       }
     )
   }

   select(id) {
     this.httpService.setSelected(id).subscribe(
       data => {
         this.sessionService.selected = id
         this.nextTime()
         //console.log(data)
       },
       error => {
         console.log(error)
       }
     )
   }

   tmp() {
   }

   selectItem(evt : MouseEvent) {
     let x = evt.clientX - 7
     let y = evt.clientY- 60
     let xr = x * this.sessionService.xmax / this.sessionService.width + this.sessionService.xmin
     let yr = y * this.sessionService.ymax / this.sessionService.height + this.sessionService.ymin
     let selectedAnt = null
     let distm = -1
     for (let ant of this.sessionService.data) {
       let dist = (ant.x - xr)*(ant.x - xr)+(ant.y - yr)*(ant.y - yr)
       if (distm==-1 || dist<distm) {
         distm = dist
         selectedAnt = ant
       }
     }
     //console.log(x +","+y+","+xr+","+yr)
     //console.log(selectedAnt)
     if (selectedAnt != null) {
       this.select(selectedAnt.id)
     }
   }

}