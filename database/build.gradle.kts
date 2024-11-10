val h2Version: String by project
val exposedVersion: String by project
val mapStructVersion: String by project
val postgresDriverVersion: String by project
val mariaDbDriverVersion: String by project
val kotlinSerializationVersion: String by project

dependencies {
    implementation(project(":core"))
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json:$kotlinSerializationVersion")
    implementation("org.jetbrains.exposed:exposed-core:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-jdbc:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-java-time:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-json:$exposedVersion")
    implementation("com.h2database:h2:$h2Version")
    implementation("org.postgresql:postgresql:$postgresDriverVersion")
    implementation("org.mariadb.jdbc:mariadb-java-client:$mariaDbDriverVersion")
}
