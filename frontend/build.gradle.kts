import com.github.gradle.node.npm.task.NpmTask

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

    val npmBuild = create<NpmTask>("npmBuild") {
        dependsOn(npmSetup, npmInstall)
        args = listOf("run", "build")
    }

    val npmClean = create<NpmTask>("npmClean") {
        dependsOn(npmSetup)
        args = listOf("run", "clean")
    }

    jar {
        dependsOn(npmLint, npmBuild)

        from(layout.buildDirectory.dir("frontend")) {
            into("/frontend/")
        }
    }

    clean {
        dependsOn(npmClean)
    }
}
