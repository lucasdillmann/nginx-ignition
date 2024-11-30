import com.github.gradle.node.npm.task.NpmTask
import com.github.gradle.node.npm.task.NpxTask

plugins {
    id("com.github.node-gradle.node") version "7.1.0"
}

node {
    download = true
    version = "22.11.0"
}

tasks {
    val npmLint = create<NpmTask>("npmLint") {
        dependsOn(npmSetup, npmInstall)
        args = listOf("run", "lint")
    }

    val prettierCheck = create<NpxTask>("npxPrettierCheck") {
        dependsOn(npmSetup, npmInstall)
        command = "npx"
        args = listOf("prettier", "--check", "src")
    }

    val npmBuild = create<NpmTask>("npmBuild") {
        dependsOn(npmSetup, npmInstall)
        args = listOf("run", "build")
    }

    val npmClean = create<NpmTask>("npmClean") {
        dependsOn(npmSetup)
        args = listOf("run", "clean")
    }

    jar {
        dependsOn(npmLint, prettierCheck, npmBuild)

        from(layout.buildDirectory.dir("frontend")) {
            into("/nginx-ignition/frontend/")
        }
    }

    clean {
        dependsOn(npmClean)
    }
}
