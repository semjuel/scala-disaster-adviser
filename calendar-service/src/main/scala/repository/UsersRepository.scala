package repository

import model.User
import scala.concurrent.Future

trait UsersRepository {
  def add(user: User): Future[Int]
  def update(user: User): Future[Int]
  def deleteUser(userId: Int): Future[Int]
  def getUser(userId: Int): Future[Seq[User]]
  def getAll: Future[Seq[User]]

  def prepareRepository(): Future[Unit]
}
