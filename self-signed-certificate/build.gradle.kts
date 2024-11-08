val bouncyCastleVersion: String by project

dependencies {
    implementation(project(":core"))
    implementation("org.bouncycastle:bcpkix-jdk18on:$bouncyCastleVersion")
}
