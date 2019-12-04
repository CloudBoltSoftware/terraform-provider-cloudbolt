pipeline {
    agent none
    parameters {
        // Allow user to select the branch to run against
        string(name: 'GIT_BRANCH', defaultValue: 'master')
        // The UUID of the credentials used to fetch the git repos
        string(name: 'GIT_CREDS_GUID', defaultValue: '')
        // Allow users to select the version of Terraform to test against
        string(name: 'TERRAFORM_VERSION', defaultValue: '0.11.14')
        // The Git repository for the Terraform Provider
        string(name: 'GIT_REPO_URL', defaultValue: 'https://github.com/CloudBoltSoftware/terraform-provider-cloudbolt.git')
        // The Git repository for the Go SDK
        string(name: 'SDK_REPO_URL', defaultValue: 'https://github.com/CloudBoltSoftware/cloudbolt-go-sdk.git')
        // The CloduBolt host we will be running the tests against
        string(name: 'CLOUDBOLT_HOST', defaultValue: '127.0.0.1')
        // The credentials used to authenticate with that host
        string(name: 'CLOUDBOLT_USERNAME', defaultValue: 'TerraformUser01')
        string(name: 'CLOUDBOLT_PASSWORD', defaultValue: 'TerraformPassword01')
        // Host and Port can be overridden, although they probably won't need to be
        string(name: 'CLOUDBOLT_PORT', defaultValue: '443')
        string(name: 'CLOUDBOLT_PROTOCOL', defaultValue: 'https')
    }
    environment {
        GO_SDK_DIR = "./cloudbolt-go-sdk"
        TERRAFORM_PROVIDER_DIR = "./terraform-provider-cloudbolt"
        TERRAFORM_PROVIDER_BIN_NAME = "terraform-provider-cloudbolt"
        TEST_DIR = "./terraform-provider-cloudbolt/examples/order-blueprint"
    }
    stages {
        stage('Build') {
            agent {
                // Building happens in the Alpine Golang container
                // We don't _need_ alpine, but the Terraform container is also alpine,
                // and we good to keep them as close as possible
                docker { image 'golang:alpine' }
            }
            environment {
                // Set the Go environment variables to be relative to the workspace directory
                GOPATH = "${env.WORKSPACE}/go"
                GOCACHE = "${env.WORKSPACE}/go/.cache"
            }
            steps {
                // Clone the CloudBolt Go SDK to build against
                dir("${env.GO_SDK_DIR}") {
                    git credentialsId: "${params.GIT_CREDS_GUID}", url: "${params.SDK_REPO_URL}", poll: false
                }
                // Clone the Terraform Provider
                dir("${env.TERRAFORM_PROVIDER_DIR}") {
                    git credentialsId: "${params.GIT_CREDS_GUID}", url: "${params.GIT_REPO_URL}", branch: "${params.GIT_BRANCH}", poll: false

                    // Override the version of the Go SDK in `go.mod` to be the local checkout
                    sh "go mod edit -replace github.com/cloudboltsoftware/cloudbolt-go-sdk/cbclient=../cloudbolt-go-sdk/cbclient"
                    sh "go mod tidy"
                    sh "go mod verify"
                    // Build the Terraform Provider
                    sh "go build -o ${env.TERRAFORM_PROVIDER_BIN_NAME}"
                }
            }
        }
        stage('Test: Order Blueprint') {
            agent {
                docker {
                    // Run the provider in the Terraform container
                    image "hashicorp/terraform:${params.TERRAFORM_VERSION}"
                    // Override the default Docker Entrypoint with an empty one
                    // so we can treat this like a normal run environment
                    // Otherwise we every `sh` step must be a `terraform` sub-command
                    args '--entrypoint='
                }
            }
            environment {
                // Pass the parameters though to Terraform as environment variables
                TF_VAR_CB_HOST          = "${params.CLOUDBOLT_HOST}"
                TF_VAR_CB_PROTOCOL      = "${params.CLOUDBOLT_PROTOCOL}"
                TF_VAR_CB_PORT          = "${params.CLOUDBOLT_PORT}"
                TF_VAR_CB_USERNAME      = "${params.CLOUDBOLT_USERNAME}"
                TF_VAR_CB_PASSWORD      = "${params.CLOUDBOLT_PASSWORD}"
                // Turn on Verbose logging for Terraform
                TF_LOG =  "1"
            }
            steps {
                // Copy the Terraform Provider into the test directory
                sh "cp -f ${env.TERRAFORM_PROVIDER_DIR}/${env.TERRAFORM_PROVIDER_BIN_NAME} ${env.TEST_DIR}"

                // Run tests in terraform-provider-cloudbolt/examples/order-blueprint
                dir("${env.TEST_DIR}") {
                    // Init
                    sh 'terraform init'
                    // Plan
                    sh 'terraform plan'
                    // Apply twice
                    // Applying the second time should be a no-op
                    sh 'terraform apply -auto-approve'
                    sh 'terraform apply -auto-approve'
                    // Show
                    sh 'terraform show'
                    // Destroy twice
                    // Destroying the second time should be a no-op
                    sh 'terraform destroy -auto-approve'
                    sh 'terraform destroy -auto-approve'
                }
            }
        }
    }
}