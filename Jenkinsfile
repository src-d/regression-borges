pipeline {
  agent {
    kubernetes {
      label 'regression-borges'
      inheritFrom 'performance'
      defaultContainer 'golang'
      containerTemplate {
        name 'golang'
        image 'golang:1.11-alpine'
        ttyEnabled true
        command 'cat'
      }
    }
  }
  environment {
    GOPATH = "/go"
    GO_IMPORT_PATH = "github.com/src-d/regression-borges"
    GO_IMPORT_FULL_PATH = "${env.GOPATH}/src/${env.GO_IMPORT_PATH}"
  }
  stages {
    stage('Prepare') {
      steps {
        sh 'apk add --no-cache git git-daemon make bash'
        sh script: """#!/bin/sh
          mkdir -p `dirname "${env.GO_IMPORT_FULL_PATH}"`
          ln -s "`pwd`" "${env.GO_IMPORT_FULL_PATH}"
        """
      }
    }
    stage('Build') {
      steps {
        sh script: """
          cd ${env.GO_IMPORT_FULL_PATH}
          go build ./cmd/regression/...
        """
      }
    }
    stage('Run') {
      steps {
        sh './regression --complexity=0 latest remote:master || true'
      }
    }
    stage('Plot') {
      steps {
        script {
          plotFiles = findFiles(glob: "plot_*.csv")
          plotFiles.each {
            echo "plot ${it.getName()}"
            sh "cat ${it.getName()}"
            plot(
              group: 'performance',
              csvFileName: it.getName(),
              title: it.getName(),
              numBuilds: '100',
              style: 'line',
              csvSeries: [[
                displayTableFlag: false,
                exclusionValues: '',
                file: it.getName(),
                inclusionFlag: 'OFF',
              ]]
            )
          }
        }
      }
    }
  }
}
