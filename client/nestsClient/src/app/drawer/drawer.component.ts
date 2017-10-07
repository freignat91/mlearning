
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
        clearInterval(this.timer)
        this.timer = setInterval(this.getData.bind(this), 100);
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
    canvasElement.width = this.sessionService.ww
    canvasElement.height = this.sessionService.hh
    this.ctx = canvasElement.getContext('2d');
    this.ctx.scale(1,1)
    this.ctx.translate(0.5, 0.5)
    this.httpService.getGlobalInfo().subscribe(
      data => {
        this.sessionService.xmin = data.xmin
        this.sessionService.ymin = data.ymin
        this.sessionService.xmax = data.xmax
        this.sessionService.ymax = data.ymax
        this.visionSize = data.ndir
        this.sessionService.selected = data.selectedAnt
        console.log(data)
      },
      error => {
        console.log(error)
      }
    )
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
    ctx.clearRect(-1,-1,this.sessionService.width+1, this.sessionService.height+1);

    for (let obj of this.sessionService.data.ants) {
      //console.log(obj)
      let angle = (Math.PI*2*obj.direction)/this.visionSize

      ctx.beginPath();
      if (obj.contact) {
        ctx.strokeStyle="red";
      } else {
        ctx.strokeStyle="black";
      }
      if (obj.id == this.sessionService.selected) {
        ctx.arc(this.getx(obj.x), this.gety(obj.y), this.getl(30), 0, 2*Math.PI, false);
      }
      ctx.moveTo(this.getx(obj.x), this.gety(obj.y));
      ctx.lineTo(this.getx(obj.x+Math.sin(angle)*3), this.gety(obj.y+Math.cos(angle)*3))
      ctx.stroke();
      if (obj.contact) {
        ctx.beginPath();
        ctx.strokeStyle="red";
        ctx.arc(this.getx(obj.x), this.gety(obj.y), this.getl(15), 0, 2*Math.PI, false);
        ctx.stroke();
      }
    }
    for (let obj of this.sessionService.data.foods) {
      //console.log(obj)
      ctx.beginPath();
      ctx.fillStyle="red";
      ctx.fillRect(obj.x, obj.y, 3, 3)
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
