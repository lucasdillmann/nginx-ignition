private val apacheCommonsValidatorVersion: String by project
private val apacheCommonsIoVersion: String by project

dependencies {
    implementation("commons-validator:commons-validator:$apacheCommonsValidatorVersion")
    implementation("commons-io:commons-io:$apacheCommonsIoVersion")
}
