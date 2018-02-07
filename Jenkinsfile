pipeline {
  agent any
  stages {
    stage('Build') {
      agent {
        docker {
          image 'golang:1.8'
        }
        
      }
      steps {
        sh 'go build ./...'
      }
    }
  }
}