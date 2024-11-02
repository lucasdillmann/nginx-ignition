//package br.com.dillmann.nginxsidewheel.application.plugins
//
//import io.ktor.http.*
//import io.ktor.server.application.*
//import io.ktor.server.request.*
//import io.ktor.server.response.*
//import io.ktor.server.routing.*
//import java.sql.*
//import kotlinx.coroutines.*
//import org.jetbrains.exposed.sql.*
//
//fun Application.configureDatabases() {
//    val database = Database.connect(
//        url = "jdbc:h2:mem:test;DB_CLOSE_DELAY=-1",
//        user = "root",
//        driver = "org.h2.Driver",
//        password = "",
//    )
//    val dbConnection: Connection = connectToPostgres(embedded = true)
//    val cityService = CityService(dbConnection)
//
//    routing {
//
//        // Create city
//        post("/cities") {
//            val city = call.receive<City>()
//            val id = cityService.create(city)
//            call.respond(HttpStatusCode.Created, id)
//        }
//
//        // Read city
//        get("/cities/{id}") {
//            val id = call.parameters["id"]?.toInt() ?: throw IllegalArgumentException("Invalid ID")
//            try {
//                val city = cityService.read(id)
//                call.respond(HttpStatusCode.OK, city)
//            } catch (e: Exception) {
//                call.respond(HttpStatusCode.NotFound)
//            }
//        }
//
//        // Update city
//        put("/cities/{id}") {
//            val id = call.parameters["id"]?.toInt() ?: throw IllegalArgumentException("Invalid ID")
//            val user = call.receive<City>()
//            cityService.update(id, user)
//            call.respond(HttpStatusCode.OK)
//        }
//
//        // Delete city
//        delete("/cities/{id}") {
//            val id = call.parameters["id"]?.toInt() ?: throw IllegalArgumentException("Invalid ID")
//            cityService.delete(id)
//            call.respond(HttpStatusCode.OK)
//        }
//    }
//    routing {
//    }
//}
