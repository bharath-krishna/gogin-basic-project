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
                docker build -t $JOB_NAME:$BUILD_NUMBER .
                docker tag $JOB_NAME:$BUILD_NUMBER docker.io/krishbharath/$JOB_NAME:$BUILD_NUMBER
                docker push docker.io/krishbharath/$JOB_NAME:$BUILD_NUMBER
              """
            }
          }
        }
      }
    }

  }
}