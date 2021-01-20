package app

import actors.SyncCalendarActor
import akka.actor.{ActorSystem, Props}

import scala.concurrent.duration.{Duration, DurationInt}
import scala.concurrent.{ExecutionContext, ExecutionContextExecutor}

object ScheduleApp extends App {
  val system: ActorSystem = ActorSystem()

  val syncActor = system.actorOf(Props(classOf[SyncCalendarActor]))

  implicit val ec: ExecutionContextExecutor = ExecutionContext.global

  val cancellable =
    system.scheduler.scheduleWithFixedDelay(Duration.Zero, 10.seconds, syncActor, SyncCalendarActor.Sync)
}
