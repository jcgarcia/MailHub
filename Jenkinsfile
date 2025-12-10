pipeline {
    agent any

    environment {
        REGISTRY = 'ghcr.io/jcgarcia'
        IMAGE_NAME = 'mailhub-admin'
        BUILD_NUMBER_TAG = "${env.BUILD_NUMBER}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                script {
                    env.GIT_COMMIT_SHORT = sh(
                        script: "git rev-parse --short HEAD",
                        returnStdout: true
                    ).trim()
                    env.IMAGE_TAG = "${BUILD_NUMBER_TAG}-${GIT_COMMIT_SHORT}"
                    echo "🚀 Building MailHub Admin"
                    echo "========================="
                    echo "📦 Image: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                    sh 'ls -la'
                }
            }
        }

        stage('Docker Build') {
            steps {
                script {
                    echo "🐳 Building Docker image (multi-stage build compiles Go)..."
                    sh """
                        docker build -t ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG} \
                                     -t ${REGISTRY}/${IMAGE_NAME}:latest .
                    """
                    echo "✅ Docker image built: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
                }
            }
        }

        stage('Docker Push') {
            steps {
                script {
                    echo "📤 Pushing Docker image to registry..."
                    withCredentials([usernamePassword(credentialsId: 'github-credentials', usernameVariable: 'GH_USER', passwordVariable: 'GH_TOKEN')]) {
                        sh '''
                            echo "$GH_TOKEN" | docker login ghcr.io -u "$GH_USER" --password-stdin
                        '''
                    }
                    sh """
                        docker push ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}
                        docker push ${REGISTRY}/${IMAGE_NAME}:latest
                    """
                    echo "✅ Docker image pushed"
                }
            }
        }

        stage('Deploy to K8s') {
            steps {
                withCredentials([file(credentialsId: 'oci-kubeconfig', variable: 'KUBECONFIG')]) {
                    script {
                        echo "☸️ Deploying to Kubernetes..."
                        sh '''
                            # Create namespace if not exists
                            kubectl get namespace mailhub 2>/dev/null || kubectl create namespace mailhub

                            # Apply K8s manifests
                            kubectl apply -f k8s/

                            # Update deployment with new image
                            kubectl -n mailhub set image deployment/mailhub-admin \
                                mailhub-admin=${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}

                            # Wait for rollout
                            kubectl -n mailhub rollout status deployment/mailhub-admin --timeout=120s
                        '''
                        echo "✅ Deployment successful"
                    }
                }
            }
        }
    }

    post {
        success {
            echo "🎉 Pipeline successful: ${REGISTRY}/${IMAGE_NAME}:${IMAGE_TAG}"
        }
        failure {
            echo '❌ Pipeline failed'
        }
    }
}
