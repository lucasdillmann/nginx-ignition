val acme4jVersion: String by project

dependencies {
    implementation(project(":core"))
    implementation("org.shredzone.acme4j:acme4j:$acme4jVersion")
}
