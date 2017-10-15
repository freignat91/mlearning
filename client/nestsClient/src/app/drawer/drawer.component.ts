
import { ViewChild, Component, Directive, ElementRef, HostListener, Input, Renderer  } from '@angular/core';
import { HttpService } from '../services/http.service'
import { SessionService } from '../services/session.service'

const httpRetryDelay = 200
const httpRetryNumber = 3

@Component({
  selector: 'app-drawer',
  template: '<canvas #drawer style="border:1px solid #d3d3d3;margin:10px"></canvas>'
})
export class DrawerComponent {
  @ViewChild('drawer') canvas;
  private ctx: any;
  timer : any
  visionSize = 8
  display = true


  constructor(private httpService : HttpService, private sessionService : SessionService) {
    sessionService.onRedraw.subscribe(
      data => {
        this.getData()
      }
    )
    sessionService.onStart.subscribe(
      data => {
        this.start()
      }
    )
    sessionService.onStop.subscribe(
      data => {
        clearInterval(this.timer)
      }
    )
  }

  ngOnInit() {
    const canvasElement = this.canvas.nativeElement;
    this.httpService.getGlobalInfo().subscribe(
      data => {
        this.sessionService.xmin = data.xmin
        this.sessionService.ymin = data.ymin
        this.sessionService.xmax = data.xmax
        this.sessionService.ymax = data.ymax
        this.sessionService.width = data.xmax-data.xmin
        this.sessionService.height = data.ymax-data.ymin
        this.visionSize = data.ndir
        this.sessionService.selected = data.selectedAnt
        canvasElement.width = this.sessionService.width
        canvasElement.height = this.sessionService.height
        this.ctx = canvasElement.getContext('2d');
        this.ctx.scale(1,1)
        this.ctx.translate(0.5, 0.5)
        this.start()
        //console.log(data)
      },
      error => {
        console.log(error)
      }
    )
  }

  start() {
    clearInterval(this.timer)
    this.timer = setInterval(this.getData.bind(this), 100);
  }

  getData() {
    this.httpService.getData().subscribe(
      data => {
        this.sessionService.data = data
        this.draw()
      },
      error => {
        console.log(error)
      }
    )
  }

  draw() {
    if (!this.display) {
      return
    }
    //console.log(this.sessionService.data)
    const ctx = this.ctx
    ctx.lineWidth = 1;
    ctx.clearRect(-1, -1, this.sessionService.width+1, this.sessionService.height+1);

    for (let obj of this.sessionService.data.foods) {
      //console.log(obj)
      if (obj.x!=0 && obj.y!=0) {
        ctx.beginPath();
        ctx.fillStyle="green";
        ctx.fillRect(obj.x, obj.y, 3, 3)
      }
    }
    let id = 0
    for (let nest of this.sessionService.data.nests) {
      let col = "blue"
      if (id != 0) {
        col="red"
      }
      id++
      //console.log(col)
      for (let obj of nest.ants) {
        //console.log(obj)
        if (obj.life>0) {
          let angle = (Math.PI*2*obj.direction)/this.visionSize
          ctx.beginPath();
          if (obj.type == 0) {
            ctx.lineWidth = 1
          } else {
            ctx.lineWidth = 2
          }
          ctx.strokeStyle = col
          ctx.moveTo(this.getx(obj.x), this.gety(obj.y));
          ctx.lineTo(this.getx(obj.x+Math.sin(angle)*3), this.gety(obj.y+Math.cos(angle)*3))
          ctx.stroke();

          if (this.sessionService.displayContact && obj.contact) {
            ctx.beginPath();
            ctx.stokeStyle="black"
            ctx.lineWidth = 1
            ctx.arc(this.getx(obj.x), this.gety(obj.y), this.getl(7), 0, 2*Math.PI, false);
            ctx.stroke();
          }
          if (id == this.sessionService.nestSelected && obj.id == this.sessionService.selected) {
            ctx.beginPath();
            ctx.stokeStyle=col
            ctx.lineWidth = 1
            ctx.arc(this.getx(obj.x), this.gety(obj.y), this.getl(30), 0, 2*Math.PI, false);
            ctx.stroke();
          }
          if (this.sessionService.displayFight && obj.fight) {
            ctx.beginPath();
            ctx.stokeStyle="orange"
            ctx.lineWidth = 1
            ctx.arc(this.getx(obj.x), this.gety(obj.y), this.getl(7), 0, 2*Math.PI, false);
            ctx.stroke();
          }
        }
      }
      //console.log(this.sessionService.data.pheromones)
      for (let obj of nest.pheromones) {
        //console.log(obj)
        if (obj.level>0) {
          ctx.beginPath();
          ctx.fillStyle="black";
          //ctx.strokeStyle="black";
          //ctx.moveTo(obj.x, obj.y)
          //ctx.lineTo(obj.x, obj.y)
          ctx.fillRect(obj.x, obj.y, 1, 1)
        }
      }
    }
  }

  getx(x : number) : number {
    return (x-this.sessionService.xmin)*this.sessionService.width/this.sessionService.xmax
  }

  gety(y : number) : number {
    return (y-this.sessionService.ymin)*this.sessionService.height/this.sessionService.ymax
  }

  getl(l : number) : number {
    return l*this.sessionService.width/(this.sessionService.xmax-this.sessionService.xmin)
  }

}
