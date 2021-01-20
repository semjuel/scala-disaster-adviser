package actors

import akka.actor.Actor
import akka.event.{Logging, LoggingAdapter}

import scala.concurrent.{ExecutionContext, ExecutionContextExecutor}

object SyncCalendarActor {
  case object Sync
}

class SyncCalendarActor extends Actor {
  val log: LoggingAdapter = Logging(context.system, this)

  def receive: Receive = {
    case SyncCalendarActor.Sync => runSync()
    case _ => log.info("received unknown message")
  }

  def runSync(): Unit = {
    implicit val ec: ExecutionContextExecutor = ExecutionContext.global




    print("Sync run....")
  }
}
