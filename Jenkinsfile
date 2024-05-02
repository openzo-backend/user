pipeline {
    agent any
    def app
stages{

    stage('Clone repository') {
       steps{

        checkout scm
       }      

    }

    stage('Build image') {
       steps{

       app = docker.build "tanush128/user:${env.BUILD_NUMBER}"
       }
    }

    stage('Test image') {
  

        app.inside {
            sh 'echo "Tests passed"'
        }
    }

    stage('Push image') {
        steps{

        docker.withRegistry('https://registry.hub.docker.com', 'dockerhub') {
            app.push("${env.BUILD_NUMBER}")
        }
        }
    }
    
    stage('Trigger ManifestUpdate') {
        steps{

                echo "triggering updatemanifestjob"
                build job: 'updatemanifest', parameters: [string(name: 'DOCKERTAG', value: env.BUILD_NUMBER)]
        }
        }
}
}