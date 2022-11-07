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
                			
                sh "sudo docker build . -t ${acrname}/gobuild:latest"
               withCredentials([usernamePassword(credentialsId: 'acrauth1', passwordVariable: 'acrpwd', usernameVariable: 'acruser')]) {
                    sh "sudo docker login ${acrname} -u ${acruser} -p ${acrpwd}"
                }
              
                sh "sudo docker push ${acrname}/gobuild:latest "
            }
        }
        
        stage('Deploy'){
            steps{
              withCredentials([usernamePassword(credentialsId: 'spauth', passwordVariable: 'PASSWORD_VAR', usernameVariable: 'USERNAME_VAR')])
		          {
			           sh ' az login --service-principal -u $USERNAME_VAR -p $PASSWORD_VAR -t $AZURE_TENANT_ID'
				  sh ' az aks get-credentials --resource-group Demo --name akscluster '
                                  sh ' kubectl apply -f goapp.yaml '
			  }
            }
        }
    }
}

def getVersion(){
    def commitHash = sh label: '', returnStdout: true, script: 'git rev-parse --short HEAD'
    return commitHash
}
