val h2Version: String by project
val exposedVersion: String by project
val mapStructVersion: String by project

dependencies {
    implementation(project(":core"))
    compileOnly("org.jetbrains.kotlinx:kotlinx-serialization-json:1.7.3")
    implementation("org.jetbrains.exposed:exposed-core:$exposedVersion")
    implementation("org.jetbrains.exposed:exposed-jdbc:$exposedVersion")
    implementation("com.h2database:h2:$h2Version")
}
