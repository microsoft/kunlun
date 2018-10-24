pipeline {
    agent { docker { image 'golang:1.10.3' } }
    stages {
        stage('unit test') {
            steps {
                sh "curl https://raw.githubusercontent.com/xplaceholder/test-infra/draft/scripts/unit_test.sh -o unit_test.sh"
                sh "chmod +x ./unit_test.sh"
                sh "./unit_test.sh"
            }
        }
    }
    options {
        buildDiscarder(logRotator(numToKeepStr:'5'))
    }
}
