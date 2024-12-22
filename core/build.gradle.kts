private val apacheCommonsValidatorVersion: String by project
private val apacheCommonsIoVersion: String by project
private val ipAddressVersion: String by project

dependencies {
    implementation("commons-validator:commons-validator:$apacheCommonsValidatorVersion")
    implementation("commons-io:commons-io:$apacheCommonsIoVersion")
    implementation("com.github.seancfoley:ipaddress:$ipAddressVersion")
}
