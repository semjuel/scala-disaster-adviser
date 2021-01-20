name := "CalendarService"

version := "1.0"

scalaVersion := "2.13.0"

lazy val akkaVersion = "2.6.11"
lazy val postgresVersion = "42.2.2"

resolvers += Resolver.jcenterRepo

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-persistence" % akkaVersion,
  "com.typesafe.akka" %% "akka-persistence-query" % akkaVersion,

  // JDBC with PostgreSQL
  "org.postgresql" % "postgresql" % postgresVersion,
  "com.github.dnvriend" %% "akka-persistence-jdbc" % "3.5.3",
)
