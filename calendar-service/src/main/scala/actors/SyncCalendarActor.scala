package actors

import akka.actor.Actor
import akka.event.{Logging, LoggingAdapter}
import slick.jdbc.PostgresProfile.api.Database

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

    val db = Database.forConfig("databaseUrl")
//    val repository = new UsersPostgresRepository(db)
//    val f = repository.getAll
//
//    //f.map()
//    f.onComplete{
//      case Success(value) => log.info("Users repository was successfully prepared")
//      case Failure(exception) => log.info("exception")
//    }

    print("Sync run....")
  }
}
