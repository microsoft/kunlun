pipeline {
    agent { docker { image 'golang:1.10.3' } }
    stages {
        stage('unit test') {
            steps {
                sh "./scripts/jenkins_unit_test.sh"
            }
        }
    }
    options {
        buildDiscarder(logRotator(numToKeepStr:'5'))
    }
}
