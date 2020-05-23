IMAGE="ImageName"
pipeline {
  agent {
    // Using kubernetes cloud to create jenkins agent
    kubernetes {
      // Pod spec for agent. Agent image created using jenkins-agent.Dockerfile
      yaml """
spec:
  containers:
  - image: "krishbharath/jenkins-slave:triad"
    imagePullPolicy: Always
    name: "docker"
    command:
      - cat
    tty: true
    volumeMounts:
    - mountPath: /var/run/docker.sock
      name: docker-socket
  restartPolicy: "Never"
  securityContext: {}
  volumes:
  - name: docker-socket
    hostPath:
      path: /var/run/docker.sock
      """
    }
  }
  stages {
    stage ("Checkout") {
      steps {
        container('docker') {
          // chekout mkdocs repo
          checkout scm

          script {
            def gitCommitTag = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
            IMAGE = "docker.io/krishbharath/$JOB_NAME:$gitCommitTag"
          }
        }
      }
    }

    stage ("Build") {
      steps {
        container('docker') {
          script {
            withCredentials([usernamePassword(credentialsId: 'docker_hub_creds', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
              sh """
                docker login --username=$DOCKER_USERNAME --password=$DOCKER_PASSWORD
                docker build -t $IMAGE .
                docker push $IMAGE
              """
            }
          }
        }
      }
    }

  }
}