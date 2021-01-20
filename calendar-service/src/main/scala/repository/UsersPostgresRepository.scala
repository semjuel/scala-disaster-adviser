package repository

import database.UsersTable
import model.User
import slick.jdbc.PostgresProfile.api._
import scala.concurrent.Future

class UsersPostgresRepository(db: Database) extends UsersRepository with UsersTable {
  override def add(user: User): Future[Int] = {
    db.run(
      users += user
    )
  }

  override def update(user: User): Future[Int] = {
    db.run(
      users.filter(_.id === user.id)
        .map(u => (u.name, u.token))
        .update((user.name, user.token))
    )
  }

  override def deleteUser(userId: Int): Future[Int] = {
    db.run(
      users.filter(_.id === userId).delete
    )
  }

  override def getUser(userId: Int): Future[Seq[User]] = {
    db.run(
      users.filter(_.id === userId).result
    )
  }

  override def getAll: Future[Seq[User]] = {
    db.run(
      users.result
    )
  }

  override def prepareRepository(): Future[Unit] = {
    db.run(
      users.schema.createIfNotExists
    )
  }
}
