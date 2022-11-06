pipeline{
    agent any
    environment {
      DOCKER_TAG = getVersion()
	  aksrg = "Demo"
	  akscls = "akscluster"
	  acrname = "aksmyrigistry.azurecr.io"
	  acr_username = "aksmyrigistry"
	  AZURETENANTID = "4ce76999-714a-427c-a127-9ec695a4a4e3"
	  AZURESUBSCRIPTIONID = "082de385-1e0e-4fd1-8476-2438ff46ed7e"
    }
    stages{
        stage('init'){
            steps{
                checkout scm
            }
        }
        
        stage('Build'){
            steps{
                sh "mvn clean package"
				
                sh "sudo docker build . -t ${acrname}/helloapp:latest"
               withCredentials([usernamePassword(credentialsId: 'acrauth', passwordVariable: 'acrpwd', usernameVariable: 'acruser')]) {
                    sh "sudo docker login ${acrname} -u ${acruser} -p ${acrpwd}"
                }
              
                sh "sudo docker push ${acrname}/helloapp:latest "
            }
        }
        
        stage('Deploy'){
            steps{
              withCredentials([usernamePassword(credentialsId: '9ee0b9dc-8213-4b07-87bc-7a1276b19349', passwordVariable: 'PASSWORD_VAR', usernameVariable: 'USERNAME_VAR')])
		          {
			           sh ' az login --service-principal -u $USERNAME_VAR -p $PASSWORD_VAR -t $AZURE_TENANT_ID'
                 sh 'mvn package azure-webapp:deploy  -Dazure.client=${USERNAME_VAR} -Dazure.key=${PASSWORD_VAR}'
			  }
            }
        }
    }
}

def getVersion(){
    def commitHash = sh label: '', returnStdout: true, script: 'git rev-parse --short HEAD'
    return commitHash
}
