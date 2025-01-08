# Configurando um Cluster Kubernetes com Minikube

Este guia fornece instruções para configurar um cluster Kubernetes local utilizando o Minikube e para implantar os recursos de um aplicativo (banco de dados e API).

## Pré-requisitos

Antes de começar, certifique-se de ter instalado os seguintes componentes:

- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

## Passos para Configuração

### 1. Inicializar o Minikube

1. **Inicie um cluster Minikube** com o seguinte comando:
   ```bash
   minikube start
   ```
   Esse comando inicializa um cluster Kubernetes local utilizando o Minikube, criando um ambiente controlado para testes e desenvolvimento.

2. **Verifique o status do cluster:**
   ```bash
   minikube status
   ```
   Este comando garante que o cluster foi iniciado corretamente e está pronto para uso.

### 2. Configurar os Recursos do Banco de Dados

1. **Aplique o Persistent Volume Claim (PVC) para o banco de dados:**
   ```bash
   kubectl apply -f k8s/pvc-db.yml
   ```
   O PVC garante que o banco de dados tenha armazenamento persistente, mesmo que o pod seja reiniciado.

2. **Implante o Deployment do banco de dados:**
   ```bash
   kubectl apply -f k8s/deployment-db.yml
   ```
   O Deployment gerencia a criação e atualização dos pods do banco de dados, garantindo a alta disponibilidade e a replicação, se configurada.

3. **Configure o Service para o banco de dados:**
   ```bash
   kubectl apply -f k8s/service-db.yml
   ```
   O Service expõe o banco de dados internamente no cluster, permitindo que outros serviços, como a API, se conectem a ele.

### 3. Configurar os Recursos da API

1. **Crie os segredos necessários para a API:**
   ```bash
   kubectl apply -f k8s/secrets-api.yml
   ```
   Os segredos armazenam informações sensíveis, como credenciais de acesso ao banco de dados, de forma segura no cluster.

2. **Implante o Deployment da API:**
   ```bash
   kubectl apply -f k8s/deployment-api.yml
   ```
   Semelhante ao banco de dados, o Deployment da API gerencia os pods que executam o serviço da aplicação.

3. **Configure o Service para a API:**
   ```bash
   kubectl apply -f k8s/service-api.yml
   ```
   Este Service expõe a API dentro do cluster e/ou para o ambiente externo, dependendo da configuração.

4. **Configure o Horizontal Pod Autoscaler (HPA) para a API:**
   ```bash
   kubectl apply -f k8s/hpa-api.yml
   ```
   O HPA ajusta automaticamente o número de réplicas do pod da API com base na carga de trabalho, garantindo escalabilidade e alta disponibilidade.

### 4. Acessar o Aplicativo

1. **Liste os serviços disponíveis no cluster:**
   ```bash
   kubectl get svc
   ```
   Esse comando exibe os serviços configurados, permitindo verificar se os serviços do banco de dados e da API estão ativos.

2. **Obtenha a URL para acessar o serviço da API:**
   ```bash
   minikube service go-app-service --url
   ```
   Este comando expõe o serviço da API localmente, retornando uma URL que pode ser acessada pelo navegador ou ferramentas de teste.

## Debug e Resolução de Problemas

- **Para verificar os logs de um pod específico:**
  ```bash
  kubectl logs <nome-do-pod>
  ```
  Use este comando para identificar possíveis erros ou mensagens de log dos pods.

- **Para acessar um pod interativamente:**
  ```bash
  kubectl exec -it <nome-do-pod> -- /bin/bash
  ```
  Este comando permite executar comandos diretamente dentro do pod, útil para depuração manual.

- **Para inspecionar os eventos do cluster:**
  ```bash
  kubectl get events
  ```
  Este comando ajuda a monitorar eventos que podem indicar problemas ou mudanças no cluster.