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

@NonCPS
def goBuild(serviceList) {

        //loop through the services directory to compile and build the docker images

       serviceList.each {

            // cd $WORKSPACE/services/$it
            //     docker run --rm -v "$WORKSPACE":/usr/src/kcapture -w /usr/src/kcapture -e CGO_ENABLED=0 golang:${goVersion} go build -o $it
            //     ldd $it | grep 'not a dynamic executable'

            //     cd ../

            sh """
                ls -ltr $it


            """
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

        stash includes: 'k8sfile.yaml', name: 'k8syaml'
    }

    // stage('Test') {
    //     milestone(100)

    //     buildCommand("go vet .")

    // }

    stage('Build') {
        milestone(200)

        //buildCommand("go build .")
        goBuild(serviceList)
        // sh """
        //     ls -ltr $WORKSPACE

        //     ldd kcapture | grep 'not a dynamic executable'

        // """

    }

    // stage('Build Docker image') {
    //     milestone(300)

    //     image = docker.build("${dockerRegistry}/${dockerImageName}:${env.VERSION}")
    // }

    // stage('Publish Docker image') {
    //     milestone(400)

    //         withDockerRegistry(url: "https://${dockerRegistry}", credentialsId: dockerRegistryCredential) {
    //         image.push()
    //     }
    // }
}


// node( Cluster.AT4D_C4.deployAgent() ) {
//     unstash 'k8syaml'

//     //Liaison at4d-c4 Dev/QA
//         deploymentAt4dC4 = deployments.create(
//             name: "${k8sDeployName}",
//             version: "${dockerImageVer}",
//             description: "${k8sDeployName}",
//             dockerImageName: "${dockerImageName}",   // Without registry!
//             dockerImageTag: "${dockerImageVer}",
//             yamlFile: 'k8sfile.yaml',   // optional, defaults to 'K8sfile.yaml'
//             gitUrl: env.GIT_URL,        // optional, defaults to env.GIT_URL
//             gitCommit: env.GIT_COMMIT,  // optional, defaults to env.GIT_COMMIT
//             gitRef: env.VERSION,        // optional, defaults to env.GIT_COMMIT
//             kubectl: kubectl,
//             clusters: [ Cluster.AT4D_C4 ]
// )

//     stage('validate on at4d-c4') {
//             milestone(500)

//             kubectl.validate(deploymentAt4dC4, Cluster.AT4D_C4)
//     }

//     if("master" == env.BRANCH_NAME) {

//         stage('deploy to at4d-c4') {
//             milestone(600)

//             kubectl.deploy(deploymentAt4dC4, Cluster.AT4D_C4)
//             kubectl.rolloutStatus(deploymentAt4dC4, Cluster.AT4D_C4)
//         }
//     }
// }
