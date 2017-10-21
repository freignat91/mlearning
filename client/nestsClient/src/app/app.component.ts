
import { ViewChild, Component, Directive, ElementRef, HostListener, Input, Renderer  } from '@angular/core';
import { HttpService } from './services/http.service'
import { SessionService } from './services/session.service'

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
  messageStartStop="Start"
  speed = 1
  logLevel = 1
  timer : any
  info : any
  tindex = 0

  constructor(private httpService : HttpService, public sessionService : SessionService) {
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

  ngOnInit() {
    this.start()
  }

  startStop() {
    if (this.messageStartStop == "Start") {
      this.sessionService.start()
    } else {
      this.sessionService.stop()
    }
  }


  start() {
    this.messageStartStop = "Stop"
    console.log("starting")
    this.httpService.start().subscribe(
      data => {
        clearInterval(this.timer)
        this.getInfo()
        this.timer = setInterval(this.getInfo.bind(this), 3000);
        console.log("started")
        this.sessionService.started = true
      },
      error => {
        console.log(error)
      }
    )
  }

   stop() {
     this.messageStartStop = "Start"
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

   restart(nb) {
     this.sessionService.stop()
     this.httpService.restart(nb).subscribe(
       data => {
        this.httpService.getGlobalInfo().subscribe(
          info => {
            this.sessionService.nestsInfo = info.nests
            this.sessionService.start()
             console.log("restarted")
           },
           error => {
             console.log(error)
           }
         )
       }
     )
   }

   getInfo() {
     this.httpService.getInfo().subscribe(
       data => {
         //console.log(data)
         this.info = data
         if (!this.info.selectedInfo) {
           this.info.selectedInfo = { gRate:0 }
         }
         this.sessionService.foodRenew = data.foodRenew
         this.sessionService.panicMode = data.panicMode
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
     //console.log(value+": "+this.speed)
     this.speed = this.speed * value
     if (this.speed <1) {
       this.speed = 0
     }
     if (value == 2 && this.speed == 0) {
       this.speed = 1
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

   select(nestId, antId) {
     this.httpService.setSelected(nestId, antId).subscribe(
       data => {
         this.sessionService.nestSelected = nestId
         this.sessionService.selected = antId
         //this.sessionService.redraw()
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

   clickEvent(evt : MouseEvent) {
     if (this.sessionService.mode == "select") {
       this.selectItem(evt)
       return
     } else if (this.sessionService.mode == "setfoodGroup") {
       this.setFoodGroup(evt)
     } else if (this.sessionService.mode == "pan") {
       this.pan(evt)
     }
   }

  pan(evt : MouseEvent) {
    let x = evt.clientX - 20
    let y = evt.clientY + 60
    this.sessionService.fpanx = this.sessionService.getInvx(x)
    this.sessionService.fpany = this.sessionService.getInvx(y)
    this.sessionService.setZoom(1)
    this.sessionService.redraw()
  }

   setFoodGroup(evt : MouseEvent) {
     let x = evt.clientX - 20
     let y = evt.clientY + 60
     let xr = this.sessionService.getInvx(x)
     let yr = this.sessionService.getInvx(y)
     this.httpService.addFoods(xr, yr).subscribe(
       data => {
         console.log("foods added")
         this.nextTime()
       },
       error => {
         console.log(error)
       }
     )
   }

   selectItem(evt : MouseEvent) {
     let x = evt.clientX - 20
     let y = evt.clientY + 60
     let xr = this.sessionService.getInvx(x)
     let yr = this.sessionService.getInvx(y)
     let selectedAnt = undefined
     let selectedNest = 0
     let distm = 30*30
     let selectedNestId = 0
     let selectedAntId = 0
     let id=0
     for (let nest of this.sessionService.data.nests) {
       id++
       for (let ant of nest.ants) {
         let dist = (ant.x - xr)*(ant.x - xr)+(ant.y - yr)*(ant.y - yr)
         if (dist<distm) {
           distm = dist
           selectedAntId = ant.id
           selectedNestId = id
         }
       }
     }
     //console.log(selectedNestId+"-"+selectedAntId)
     this.select(selectedNestId, selectedAntId)
   }

   addFoods() {
     this.sessionService.mode = "setfoodGroup"
   }

   selectAnt() {
     this.sessionService.mode = "select"
   }

   zoomPan() {
     this.sessionService.mode = "pan"
   }


   foodRenew(evt) {
     this.httpService.foodRenew(evt.target.checked).subscribe(
       data => {
         console.log("food renew set to: "+evt.target.cheked)
       },
       error => {
         console.log(error)
       }
     )
   }

   panicMode(evt) {
     this.httpService.panicMode(!evt.target.checked).subscribe(
       data => {
         console.log("panic mode set to: "+evt.target.cheked)
       },
       error => {
         console.log(error)
       }
     )
   }

   clearFoodGroup() {
     this.httpService.clearFoodGroup().subscribe(
       data => {
         console.log("food groups cleared")
       },
       error => {
         console.log(error)
       }
     )
   }

   fightCircles(evt) {
     this.sessionService.displayFight = false
     if (evt.target.checked) {
       this.sessionService.displayFight = true
     }
   }

   displayGraph(evt) {
     this.sessionService.display = evt.target.checked
     if (!this.sessionService.display) {
       this.sessionService.clear()
     }
   }

   displayTableToggle() {
     if (this.sessionService.nestsInfo.length <= 2) {
       this.tindex = 0
       return
     }
     if (this.tindex == 0) {
       this.tindex = 2
     } else {
       this.tindex = 0
     }
   }
}
