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

    stage ("test_run") {
      steps {
        container('docker') {
          // Create test pod and sleep untill ready and then run tests
          sh """
            kubectl get po
            docerk ps -a
          """
        }
      }
    }




    // Push production image if above tests passes
    // stage ("Push") {
    //   steps {
    //     container('docker') {
    //       script {
    //         // Build image only if the above stages succeeds
    //         sh """
    //           docker push krishbharath/mkdocs_image
    //         """
    //       }
    //     }
    //   }
    // }


  }
}