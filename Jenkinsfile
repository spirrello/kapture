#!groovy


@Library('visibilityLibs')
import com.liaison.jenkins.visibility.Utilities;
import com.liaison.jenkins.common.kubernetes.*


def deployments = new Deployments()
def k8sDocker = new Docker()
def utils = new Utilities()
def kubectl = new Kubectl()
def dockerRegistry = "registry-ci.at4d.liacloud.com"
def dockerRegistryCredential = "registry-ci"
// def dockerImageName = "devops/kcapture"
def dockerImageName = "devops/"
def k8sDeployName = "kcapture"
def serviceList = ["kcapture-api","kcapture-node"]

//Golang version
goVersion = "1.12.7"


def buildCommand(command) {

        sh """
            docker run --rm -v "$WORKSPACE":/usr/src/kcapture -w /usr/src/kcapture -e CGO_ENABLED=0 golang:${goVersion} ${command}

        """
}


def goBuild(serviceList) {

        serviceList.each {

            sh """
                docker build -t $it:$env.VERSION -f $it-docker .
                docker images
            """
        }
}

/*
Function for handling validation and deployment of the kcapture pods.  This helps with managing the deployment
and daemonset.
*/
def k8sProcessing(kubectl, deployments, dockerImageVer, dockerImageName, serviceList, process) {

    serviceList.each {
        //Liaison at4d-c4 Dev/QA
        k8sfile = "k8s/$it" + ".yaml"
        deploymentAt4dC4 = deployments.create(
            name: "$it",
            version: "${dockerImageVer}",
            description: "$it",
            dockerImageName: "${dockerImageName}$it",   // Without registry!
            dockerImageTag: "${dockerImageVer}",
            yamlFile: "${k8sfile}",   // optional, defaults to 'K8sfile.yaml'
            gitUrl: env.GIT_URL,        // optional, defaults to env.GIT_URL
            gitCommit: env.GIT_COMMIT,  // optional, defaults to env.GIT_COMMIT
            gitRef: env.VERSION,        // optional, defaults to env.GIT_COMMIT
            kubectl: kubectl,
            namespace: Namespace.KUBE_SYSTEM,
            clusters: [ Cluster.AT4D_C4 ]
        )


        switch(process) {
            case "validate":
                kubectl.validate(deploymentAt4dC4, Namespace.KUBE_SYSTEM, Cluster.AT4D_C4)
            case "deploy":
                kubectl.deploy(deploymentAt4dC4, Namespace.KUBE_SYSTEM, Cluster.AT4D_C4)
                kubectl.rolloutStatus(deploymentAt4dC4, Namespace.KUBE_SYSTEM, Cluster.AT4D_C4)
            default:
                println("No valid options submitted.")
        }
    }

}

node {

    stage('Checkout') {
        checkout scm
        env.VERSION = utils.runSh("awk '/^## \\[([0-9])/{ print (substr(\$2, 2, length(\$2) - 2));exit; }' CHANGELOG.md")
        env.GIT_COMMIT = utils.runSh('git rev-parse HEAD')
        env.GIT_URL = utils.runSh("git config remote.origin.url | sed -e 's/\\(.git\\)*\$//g' ")
        env.REPO_NAME = utils.runSh("basename -s .git ${env.GIT_URL}")
        env.RELEASE_NOTES = utils.runSh("awk '/## \\[${env.VERSION}\\]/{flag=1;next}/## \\[/{flag=0}flag' CHANGELOG.md")
        currentBuild.displayName  = "${env.VERSION}-${env.BUILD_NUMBER}"
        dockerImageVer = env.VERSION

        stash includes: 'k8s/*.yaml', name: 'k8syaml'
    }


    stage('Build/Publish Docker Images') {
        milestone(300)

        serviceList.each {
            image = docker.build("${dockerRegistry}/${dockerImageName}$it:${env.VERSION}", "-f $it-docker .")
            withDockerRegistry(url: "https://${dockerRegistry}", credentialsId: dockerRegistryCredential) {
                image.push()
            }
        }
    }

}


node( Cluster.AT4D_C4.deployAgent() ) {
    unstash 'k8syaml'

    stage('validate on at4d-c4') {
            milestone(500)

            k8sProcessing(kubectl, deployments, dockerImageVer, dockerImageName, serviceList, "deploy")
    }

    if("master" == env.BRANCH_NAME) {

        stage('deploy to at4d-c4') {
            milestone(600)

            k8sProcessing(kubectl, deployments, dockerImageVer, dockerImageName, serviceList, "deploy")

        }
    }
}
