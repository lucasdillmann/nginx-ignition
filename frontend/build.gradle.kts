import com.github.gradle.node.npm.task.NpmTask
import com.github.gradle.node.npm.task.NpxTask

plugins {
    id("com.github.node-gradle.node") version "7.1.0"
}

node {
    download = true
    version = "22.11.0"
}

sonar {
    properties {
        property("sonar.sources", "src")
    }
}

tasks {
    val ideMode = System.getProperty("nginx-ignition.dev-mode") == "true"

    val eslint = create<NpxTask>("eslint") {
        dependsOn(npmSetup, npmInstall)
        command = "npx"
        args = listOf("eslint", "--max-warnings", "0", "src")
    }

    val prettierCheck = create<NpxTask>("prettierCheck") {
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

    check {
        if (!ideMode)
            dependsOn(eslint, prettierCheck)
    }

    jar {
        if (!ideMode)
            dependsOn(npmBuild)

        from(layout.buildDirectory.dir("frontend")) {
            into("/nginx-ignition/frontend/")
        }
    }

    clean {
        if (!ideMode)
            dependsOn(npmClean)
    }
}
