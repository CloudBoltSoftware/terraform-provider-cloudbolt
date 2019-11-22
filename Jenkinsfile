pipeline {
    agent none
    stages {
        stage('Cleanup') {
            agent any
            steps {
                sh 'rm -rf /var/jenkins_home/go/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt*'
                sh 'rm -rf /var/jenkins_home/go/src/github.com/cloudboltsoftware/cloudbolt-go-sdk*'
            }
        }
        stage('Build-provider') {
            agent {
                docker { image 'golang' }
            }
            steps {
                withEnv(['GOPATH=/var/jenkins_home/go', 'GOCACHE=/var/jenkins_home/go/.cache']) {
                    git credentialsId: '843338a5-615e-47be-a267-6e7c9b266843', url: 'https://github.com/CloudBoltSoftware/terraform-provider-cloudbolt.git'
                    sh 'cp -r . $GOPATH/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt'
                    git credentialsId: '843338a5-615e-47be-a267-6e7c9b266843', url: 'https://github.com/CloudBoltSoftware/cloudbolt-go-sdk.git'
                    sh 'cp -r . $GOPATH/src/github.com/cloudboltsoftware/cloudbolt-go-sdk'
                    sh 'echo "replace github.com/cloudboltsoftware/cloudbolt-go-sdk => ../cloudbolt-go-sdk" >> "$GOPATH/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt/go.mod"'
                    dir("$GOPATH/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt") {
                        sh 'go build -o terraform-provider-cloudbolt'
                        sh 'mv terraform-provider-cloudbolt /var/jenkins_home/.terraform.d/plugins'
                    }
                }
            }
        }
        stage('Download-terraform-0-11-14') {
            agent any
            when { expression { !fileExists('./terraform_0.11.14_linux_amd64.zip') } }
            steps {
                sh 'wget https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip'
            }
        }
        stage('Unzip-terraform') {
            agent any
            when { expression { !fileExists('./terraform') } }
            steps {
                sh '/usr/bin/unzip terraform_0.11.14_linux_amd64.zip'
            }
        }
        stage('Test') {
            agent any
            steps {
                sh 'cp terraform /var/jenkins_home/go/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt/examples/order-blueprint/'
                dir('/var/jenkins_home/go/src/github.com/cloudboltsoftware/terraform-provider-cloudbolt/examples/order-blueprint/') {
                    withEnv(['TF_VAR_CB_HOST=docker.for.mac.localhost']) {
                        sh 'cp terraform.tfvars.dist terraform.tfvars'
                        sh './terraform init'
                        sh './terraform plan -var "CB_HOST=$TF_VAR_CB_HOST"'
                        sh './terraform apply -auto-approve -var "CB_HOST=$TF_VAR_CB_HOST"'
                        sh './terraform apply -auto-approve -var "CB_HOST=$TF_VAR_CB_HOST"'
                        sh './terraform show'
                        sh './terraform destroy -auto-approve -var "CB_HOST=$TF_VAR_CB_HOST"'
                        sh './terraform destroy -auto-approve -var "CB_HOST=$TF_VAR_CB_HOST"'
                    }
                }
            }            
        }
    }
}

