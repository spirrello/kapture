#!groovy


@Library('visibilityLibs')
import com.liaison.jenkins.visibility.Utilities;
// import com.liaison.jenkins.common.testreport.TestResultsUploader
// import com.liaison.jenkins.common.sonarqube.QualityGate
import com.liaison.jenkins.common.kubernetes.*
// import com.liaison.jenkins.common.e2etest.*
// import com.liaison.jenkins.common.servicenow.ServiceNow
// import com.liaison.jenkins.common.slack.*
// import com.liaison.jenkins.common.releasenote.ReleaseNote

def deployments = new Deployments()
def k8sDocker = new Docker()
def utils = new Utilities()
def kubectl = new Kubectl()
def dockerRegistry = "registry-ci.at4d.liacloud.com"
def dockerRegistryCredential = "registry-ci"
def dockerImageName = "devops/kcapture"

//Golang version
goVersion = "1.12.7"


def buildCommand(command) {

        sh """
            docker run --rm -v "$WORKSPACE":/usr/src/kcapture -w /usr/src/kcapture golang:${goVersion} ${command}

        """
}

timestamps {
    node {

        stage('Checkout') {
            checkout scm
            env.VERSION = utils.runSh("awk '/^## \\[([0-9])/{ print (substr(\$2, 2, length(\$2) - 2));exit; }' CHANGELOG.md")
            env.GIT_COMMIT = utils.runSh('git rev-parse HEAD')
            env.GIT_URL = utils.runSh("git config remote.origin.url | sed -e 's/\\(.git\\)*\$//g' ")
            env.REPO_NAME = utils.runSh("basename -s .git ${env.GIT_URL}")
            env.RELEASE_NOTES = utils.runSh("awk '/## \\[${env.VERSION}\\]/{flag=1;next}/## \\[/{flag=0}flag' CHANGELOG.md")
            currentBuild.displayName  = env.VERSION

            // deployment = deployments.create(
            //     name: 'File Upload',
            //     version: env.VERSION,
            //     description: 'File Upload',
            //     dockerImageName: dockerImageName,
            //     dockerImageTag: env.VERSION,
            //     yamlFile: 'K8sfile.yaml',   // optional, defaults to 'K8sfile.yaml'
            //     gitUrl: env.GIT_URL,        // optional, defaults to env.GIT_URL
            //     gitCommit: env.GIT_COMMIT,   // optional, defaults to env.GIT_COMMIT
            //     kubectl: kubectl
            // )
        }

        stage('Test') {

            buildCommand("go vet .")

        }

        stage('Build') {

            buildCommand("go build .")

        }


        stage('Build Docker image') {
            // k8sDocker.build(imageName: dockerImageName);
            milestone label: 'Docker image built', ordinal: 100
            image = docker.build("${dockerRegistry}/${dockerImageName}:${env.VERSION}")
        }

        stage('Publish Docker image') {
             withDockerRegistry(url: "https://${dockerRegistry}", credentialsId: dockerRegistryCredential) {
                image.push()
            }
        }
        // Exit the node
        //return
    }

        if("master" == env.BRANCH_NAME) {
            stage('Deploy to at4d-c3') {

            echo "Not there yet."
                // try {
                //   deployments.deploy(
                //           deployment: deployment,
                //           kubectl: kubectl,
                //           serviceNow: null,
                //           namespace: Namespace.DEVELOPMENT,
                //           rollingUpdate: true     // optional, defaults to true
                //   )

                // } catch (err) {
                //   currentBuild.result = "FAILURE";
                //   error "${err}"
                // }
                //   }

           }
        }
}