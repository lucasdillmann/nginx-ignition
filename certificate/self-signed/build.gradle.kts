val bouncyCastleVersion: String by project

dependencies {
    implementation(project(":core"))
    implementation(project(":certificate-commons"))
    implementation("org.bouncycastle:bcpkix-jdk18on:$bouncyCastleVersion")
}
