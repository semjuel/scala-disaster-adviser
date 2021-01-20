package database

import model.User
import slick.jdbc.PostgresProfile.api._
import slick.lifted.ProvenShape

trait UsersTable {
  class Users(tag: Tag) extends Table[User](tag, Some("slick"), "users") {
    def id: Rep[Int] = column[Int]("id", O.PrimaryKey)

    def name: Rep[String] = column[String]("name")

    def token: Rep[String] = column[String]("token")

    def * : ProvenShape[User] = (id, name, token) <> (User.tupled, User.unapply)
  }

  val users = TableQuery[Users]
}
