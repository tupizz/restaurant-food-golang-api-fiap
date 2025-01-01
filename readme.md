# Event Storming (DDD)

Event Storming é uma técnica colaborativa que ajuda a descobrir, compreender e mapear processos de negócios complexos através da identificação de eventos, comandos, agregados, políticas, participantes, entre outros elementos. É especialmente útil em projetos de DDD, pois facilita a comunicação entre equipes técnicas e de negócio.

## Fluxo DDD Miro

[Board DDD](https://miro.com/app/board/uXjVLeiK8DM=/)

## Linguagem Ubíqua

- Pedido
- Cliente
- Pagamento
- Preparação
- Entrega
- Status do Pedido
- Notificação


## Realização do pedido e pagamento
```mermaid
sequenceDiagram
    participant Cliente
    participant SistemaAutoatendimento as Sistema de Autoatendimento
    participant SistemaPagamento as Sistema de Pagamento
    participant Cozinha
    participant PainelCliente as Painel do Cliente

    Cliente->>SistemaAutoatendimento: Iniciar Pedido
    SistemaAutoatendimento->>Cliente: Exibir Menu de Produtos

    loop Seleção de Itens
        Cliente->>SistemaAutoatendimento: Adicionar Item ao Pedido
        SistemaAutoatendimento->>Cliente: Confirmar Adição
    end

    Cliente->>SistemaAutoatendimento: Confirmar Pedido
    SistemaAutoatendimento->>SistemaPagamento: Processar Pagamento via QRCode
    SistemaPagamento-->>Cliente: Exibir QRCode do Mercado Pago
    Cliente->>SistemaPagamento: Realizar Pagamento
    SistemaPagamento-->>SistemaAutoatendimento: Pagamento Confirmado

    SistemaAutoatendimento->>Cozinha: Enviar Pedido
    SistemaAutoatendimento->>PainelCliente: Atualizar Status para "Recebido"
```
## Preparação e Entrega do pedido

```mermaid
sequenceDiagram
    participant Cozinha
    participant PainelCliente as Painel do Cliente
    participant Cliente

    Cozinha->>Cozinha: Iniciar Preparação
    Cozinha->>PainelCliente: Atualizar Status para "Em Preparação"

    Cozinha->>Cozinha: Finalizar Preparação
    Cozinha->>PainelCliente: Atualizar Status para "Pronto"
    PainelCliente-->>Cliente: Notificar Pedido Pronto

    Cliente->>Cozinha: Retirar Pedido
    Cozinha->>PainelCliente: Atualizar Status para "Finalizado"
```

## Fluxo detalhado com diagrama de atividades

```mermaid
flowchart TD
    A[Iniciar Pedido] --> B{Cliente Identificado?}
    B -- Sim --> C[Registrar Cliente]
    B -- Não --> D[Prosseguir sem Identificação]
    C --> D
    D --> E[Selecionar Itens]
    E --> F[Confirmar Pedido]
    F --> G[Processar Pagamento]
    G --> H{Pagamento Aprovado?}
    H -- Sim --> I[Registrar Pedido]
    H -- Não --> J[Notificar Falha no Pagamento]
    I --> K["Atualizar Status para status Recebido"]
    K --> L[Enviar Pedido para Cozinha]

```


# Guia de instalação e execução do projeto

Este guia irá ajudá-lo a configurar e executar o projeto **FastFood Golang** em sua máquina, seja utilizando Docker ou rodando a aplicação diretamente. Siga as instruções abaixo para preparar o ambiente de desenvolvimento e executar a aplicação.

## Pré-requisitos

- **Git**: Para clonar o repositório.
- **Golang**: Versão 1.18 ou superior.
- **Docker** e **Docker Compose**: Se preferir executar a aplicação em contêineres.
- **Air**: Ferramenta para live reloading durante o desenvolvimento.
- **Golang-Migrate**: Para gerenciar migrações de banco de dados.
- **Swagger**: Para documentação da API.

---

## Configuração do ambiente

### 1. Clonar o repositório

Abra o terminal e clone o repositório para a sua máquina local:

```bash
git clone https://github.com/tupizz/restaurant-food-golang-api-fiap
cd restaurant-food-golang-api-fiap
```

### 2. Instalar dependências Go

Certifique-se de ter o Go instalado e configurado em sua máquina. Baixe as dependências do projeto:

```bash
go mod download
```

### 3. Instalar o Air para Live Reloading

Air é uma ferramenta que recompila e reinicia automaticamente a aplicação quando mudanças no código são detectadas nos arquivos mapeados.

#### Instalação

Ou, se preferir, instale via Go (confirme que o diretório `$GOPATH/bin` está no seu `PATH`):

```bash
go install github.com/air-verse/air@latest
```

### 4. Instalar o Golang-Migrate para migrações de banco de dados

Golang-Migrate é usado para gerenciar migrações do banco de dados.

#### Instalação

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Certifique-se de que o diretório `$GOPATH/bin` está no seu `PATH` para acessar o comando `migrate`.

---

#### Observações (passo 3 e 4):

1. Caso você use alguma ferramenta para gerenciar diferentes versões do Go, como o [ASDF](https://github.com/asdf-vm/asdf), você precisrá regerar os _shims_ para que os binários instalados diretamente com o `go install` estejam disponíveis.

```bash
# Para o cenário do ASDF, pode ser:
asdf reshim golang
```

2. Os seguintes comandos a seguir podem verificar a correta isntalação das ferramentas acima:

```bash
air -v
migrate -version
```

---

## Executando o projeto com Docker

### 1. Configurar variáveis de ambiente

Crie um arquivo `.env` na raiz do projeto com o seguinte conteúdo:

```env
DATABASE_URL=postgres://postgres:postgres@db:5432/yourdb?sslmode=disable
```

### 2. Construir e iniciar os serviços com Docker Compose

Execute o seguinte comando para construir as imagens e iniciar os contêineres:

```bash
docker-compose up --build
```

Isso irá:

- Construir a imagem Docker da aplicação Go.
- Iniciar o contêiner do banco de dados PostgreSQL.
- Executar as migrações do banco de dados.
- Iniciar o contêiner da aplicação Go com o Air para live reloading.

### 3. Acessar a aplicação

A aplicação estará disponível em `http://localhost:8080`.

#### Testar endpoints

- **Listar usuários:**

  ```bash
  curl http://localhost:8080/api/v1/users
  ```

- **Criar usuário:**

  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"João Silva", "email":"joao.silva@example.com", "age":30}' http://localhost:8080/api/v1/users/
  ```

---

## Executando o projeto sem Docker

### 1. Configurar o banco de dados PostgreSQL (2 formas):

#### 1.1 Rodando o banco de dados localmente

Instale o PostgreSQL em sua máquina e crie um banco de dados chamado `fiap_fast_food`.

Atualize a variável `DATABASE_URL` no arquivo `.env` para apontar para o seu banco de dados local:

```env
DATABASE_URL=postgres://postgres:suasenha@localhost:5432/fiap_fast_food?sslmode=disable
```

#### 1.2 Rodando o banco de dados em um container Docker

Também é possível utilizar o container para o postgres disponível no `docker-compose.yml`. Neste caso basta subir apenas este container e rodar somente o _app_ localmente:

**Caso prefira rodar o banco em sua máquina, isto é, sem uso de containers, desconsidere este passo**.

```bash
docker-compose up db --build
```

Atualize a variável `DATABASE_URL` no arquivo `.env` para apontar para o seu banco de dados do container (as credenciais podem ser vistas como variáveis de ambiente no `docker-compose.yml`):

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable
```

### 2. Executar migrações do banco de dados

Execute as migrações para criar as tabelas necessárias (subtituindo a DATABASE_URL abaixo):

```bash
migrate -database "${DATABASE_URL}" -path ./database/migrations up
```

Ou, se preferir, rode através do [make](https://www.gnu.org/software/make/).

```bash
make migrate-up
```

### 3. Iniciar a aplicação com Air

Inicie a aplicação usando o Air para habilitar o live reloading:

```bash
air
```

Ou, se preferir, rode através do [make](https://www.gnu.org/software/make/).

```bash
make run-air
```

**Observação:** Certifique-se de que o comando `air` está disponível no seu `PATH`. Se instalou o Air via Go, o binário estará em `$GOPATH/bin`.

### 4. Acessar a aplicação

A aplicação estará disponível em `http://localhost:8080`. Utilize os mesmos comandos mencionados anteriormente para testar os endpoints.

---

## Dicas e Solução de Problemas

- **Portas em Uso:** Verifique se as portas `8080` (aplicação) e `5432` (banco de dados) estão livres.
- **Variáveis de Ambiente:** Certifique-se de que o `DATABASE_URL` está corretamente configurado no arquivo `.env`.
- **Permissões de Arquivo:** Se encontrar problemas de permissão, ajuste as permissões dos arquivos e diretórios:

  ```bash
  chmod -R 755 ./fastfood-golang
  ```

- **Logs da aplicação:** Monitore os logs para identificar possíveis erros:

  ```bash
  docker-compose logs -f
  ```

- **Reinstalar dependências:** Se encontrar erros relacionados a dependências, execute:

  ```bash
  go mod tidy
  go mod download
  ```

---

## Solução de problemas com migrations

Para mais detalhes, consulte o arquivo [MIGRATION_GUIDE.md](./docs/migrations.md).

Sample error:
```
error: Dirty database version 7. Fix and force version.
```

Solution:
- Force the past version
- Update again

```bash
migrate -path ./database/migrations -database "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" force 6
migrate -path ./database/migrations -database "postgres://postgres:postgres@localhost:5432/fiap_fast_food?sslmode=disable" up
```

---

## Estrutura do projeto

- **`cmd/main.go`**: Ponto de entrada da aplicação.
- **`internal/`**: Código interno da aplicação.
    - **`adapter/`**: Adaptadores para handlers HTTP e implementações de repositórios.
    - **`application/`**: Lógica de negócio e serviços.
    - **`domain/`**: Entidades de negócio e interfaces de repositórios.
    - **`config/`**: Carregamento e gerenciamento de configurações.
    - **`di/`**: Configuração de injeção de dependências com o Uber Dig.
- **`database/`**: Configuração do banco de dados com o SQLC.
    - **`migrations/`**: Arquivos de migração do banco de dados.
    - **`queries/`**: Arquivos de queries do banco de dados.
    - **`sqlc`**: Arquivos gerados pelo SQLC.
- **`Dockerfile`**: Arquivo Docker para construir a imagem da aplicação.
- **`docker-compose.yml`**: Configuração do Docker Compose.
- **`go.mod` e `go.sum`**: Arquivos de módulos do Go.

---

# Explicação detalhada das camadas e interações no projeto

O projeto **FastFood Golang** foi estruturado seguindo os princípios da **Arquitetura Limpa** e **Arquitetura Hexagonal**, com o objetivo de criar um sistema modular, escalável e de fácil manutenção. Abaixo, detalhamos cada camada, suas responsabilidades, como elas interagem entre si e referências aos arquivos relevantes.

## Camadas da arquitetura

1. [Camada de apresentação (presentation layer)](#1-camada-de-apresentação-presentation-layer)
2. [Camada de aplicação (application layer)](#2-camada-de-aplicação-application-layer)
3. [Camada de domínio (domain layer)](#3-camada-de-domínio-domain-layer)
4. [Camada de infraestrutura (infrastructure layer)](#4-camada-de-infraestrutura-infrastructure-layer)

---

### 1. Camada de apresentação (presentation layer)

**Responsabilidade:**

- Gerenciar a comunicação com o mundo externo (clientes, navegadores, etc.).
- Processar solicitações HTTP, validar entradas e formatar respostas.
- Delegar a lógica de negócios para a camada de aplicação.

**Interação:**

- Recebe solicitações HTTP dos clientes.
- Invoca serviços da camada de aplicação para processar a lógica de negócios.
- Retorna respostas HTTP apropriadas aos clientes.

**Arquivos relevantes:**

- **Manipuladores HTTP (Handlers):**
    - `internal/adapter/http/handler/user_handler.go`: Contém a implementação do `UserHandler`, que processa as solicitações relacionadas a usuários.
- **Roteador HTTP:**
    - `internal/adapter/http/router.go`: Configura as rotas da API e registra os manipuladores correspondentes.

**Detalhes e decisões:**

- **Uso do Gin Framework:** Escolhemos o Gin devido à sua eficiência e facilidade de uso para criar APIs RESTful (foram cogitados outros frameworks, como o echo).
- **Responsabilidade Limitada dos Handlers:** Os manipuladores focam em processar solicitações e respostas, delegando a lógica de negócios para a camada de aplicação.

---

### 2. Camada de aplicação (application layer)

**Responsabilidade:**

- Implementar a lógica de negócios da aplicação.
- Orquestrar operações entre a camada de apresentação e a camada de domínio.
- Garantir que as regras de negócios sejam aplicadas corretamente.

**Interação:**

- Recebe chamadas dos manipuladores (handlers) da camada de apresentação.
- Utiliza entidades e interfaces da camada de domínio para processar dados.
- Chama repositórios através das interfaces definidas na camada de domínio.

**Arquivos relevantes:**

- **Serviços de Aplicação:**
    - `internal/application/service/user_service.go` (exemplo): Implementa a lógica de negócios relacionada a usuários, como criação, leitura, atualização e exclusão.
- **Data Transfer Objects (DTOs):**
    - `internal/application/dto/product_input.go` e `internal/application/dto/product_output.go` (exemplo): Define estruturas para transferência de dados entre camadas, evitando expor diretamente as entidades do domínio.

**Detalhes e decisões:**

- **Isolamento da Lógica de Negócios:** Centralizamos a lógica aqui para facilitar testes e manutenções futuras.
- **Uso de DTOs:** Facilita a validação e transformação de dados entre as camadas, promovendo a segurança e integridade dos dados.

---

### 3. Camada de domínio (domain layer)

**Responsabilidade:**

- Representar o núcleo da aplicação com as regras de negócio fundamentais.
- Definir entidades e interfaces que modelam conceitos do domínio.
- Permanecer independente de detalhes de implementação ou frameworks externos.

**Interação:**

- As entidades e interfaces são usadas pela camada de aplicação para processar dados.
- As interfaces de repositório definidas aqui são implementadas pela camada de infraestrutura.

**Arquivos relevantes:**

- **Entidades de Domínio:**
    - `internal/domain/entity/user.go`: Define a estrutura da entidade `User`, representando um usuário no sistema.
- **Interfaces de Repositório:**
    - `internal/domain/repository.go`: Declara a interface `UserRepository`, especificando os métodos que devem ser implementados para manipulação de usuários.

**Detalhes e decisões:**

- **Independência Tecnológica:** Ao não depender de frameworks ou pacotes externos, a camada de domínio permanece flexível e adaptável a mudanças.
- **Definição de Interfaces:** As interfaces permitem que diferentes implementações sejam usadas sem alterar a lógica de negócios, facilitando testes e trocas de tecnologia.

---

### 4. Camada de infraestrutura (infrastructure layer)

**Responsabilidade:**

- Fornecer implementações concretas para as interfaces definidas na camada de domínio.
- Lidar com detalhes técnicos, como acesso a bancos de dados, serviços externos, etc.

**Interação:**

- Implementa as interfaces de repositório, interagindo diretamente com o banco de dados.
- É chamada pela camada de aplicação através das interfaces do domínio.

**Arquivos relevantes:**

- **Implementação dos Repositórios:**
    - `internal/adapter/repository/user_repository.go`: Implementa `UserRepository`, realizando operações de banco de dados para a entidade `User`.
- **Configurações da Aplicação:**
    - `internal/config/config.go`: Carrega e gerencia as configurações do aplicativo, como variáveis de ambiente e strings de conexão.
- **Injeção de Dependências:**
    - `internal/di/container.go`: Configura o container de injeção de dependências usando o Uber Dig, registrando todos os provedores necessários.
- **Conexão com o Banco de Dados:**
    - `internal/di/database.go`: Estabelece a conexão com o PostgreSQL utilizando `pgxpool`, fazendo uso de uma _connection pool_ para evitar gargalos entre as _gorountines_ que competirão por acesso ao Banco de Dados.

**Detalhes e decisões:**

- **Uso do `pgxpool`:** Optamos pelo driver `pgx` para melhor performance e recursos avançados na interação com o PostgreSQL.
- **Injeção de Dependências com Uber Dig:** Facilita o gerenciamento de dependências complexas e aumenta a testabilidade da aplicação.
- **Centralização das Configurações:** Permite alterar facilmente parâmetros de configuração sem modificar o código-fonte.

---

## Como as camadas interagem

### Fluxo de uma operação: "Get All Users"

1. **Solicitação HTTP:**
    - O cliente envia uma requisição HTTP GET para `/api/v1/users`.

2. **Camada de apresentação:**
    - **Router (`router.go`):** Direciona a requisição para o `UserHandler`.
    - **Handler (`user_handler.go`):**
        - Recebe a requisição.
        - Realiza validações iniciais, se necessário.
        - Chama o método `GetAllUsers()` do `UserService`.

3. **Camada de aplicação:**
    - **Service (`user_service.go`):**
        - Aplica regras de negócio (por exemplo, filtragem, ordenação).
        - Chama o método `GetAll(ctx)` do repositório de usuários através da interface `UserRepository`.

4. **Camada de domínio:**
    - **Interface (`repository.go`):**
        - Define o contrato para `GetAll(ctx)` que deve ser implementado.
    - **Entidade (`user.go`):**
        - Estrutura de dados que representa um usuário.

5. **Camada de infraestrutura:**
    - **Repositório (`user_repository.go`):**
        - Implementa `GetAll(ctx)`, executando uma consulta SQL no banco de dados.
        - Utiliza o pool de conexões `pgxpool` para interação com o PostgreSQL.
    - **Banco de Dados:**
        - Executa a consulta e retorna os resultados para o repositório.

6. **Resposta:**
    - **Repositório:** Retorna a lista de usuários para o serviço.
    - **Serviço:** Pode aplicar transformações adicionais nos dados.
    - **Handler:** Formata a resposta (JSON) e envia de volta ao cliente com o status HTTP adequado.

### Interações chave:

- **Handlers ↔ Services:**
    - Os handlers invocam métodos dos serviços para processar a lógica de negócios.
- **Services ↔ Repositories:**
    - Os serviços chamam métodos dos repositórios através das interfaces definidas no domínio.
- **Repositories ↔ Database:**
    - Os repositórios interagem com o banco de dados, executando operações CRUD.

---

## Decisões de Arquitetura e Racionalização

### Separação de Responsabilidades

- **Clareza e Manutenibilidade:** Cada camada tem uma responsabilidade distinta, facilitando a compreensão e manutenção do código.
- **Colaboração entre Equipes:** Diferentes equipes ou desenvolvedores podem trabalhar em camadas específicas sem causar conflitos.

### Injeção de Dependências com Uber Dig

- **Gerenciamento Simplificado:** O Uber Dig permite resolver e injetar dependências de forma declarativa.
- **Testabilidade Aumentada:** Facilita a injeção de mocks ou stubs durante testes unitários.
- **Redução de Acoplamento:** Evita dependências rígidas entre componentes, promovendo um design mais flexível.

### Uso de interfaces e abstrações

- **Flexibilidade:** Permite trocar implementações (por exemplo, substituir o banco de dados) sem alterar as camadas superiores.
- **Isolamento da Lógica de Negócios:** A lógica de negócios não depende de detalhes de infraestrutura, seguindo o princípio da inversão de dependência.

### Gerenciamento centralizado de configurações

- **Segurança:** Variáveis sensíveis são gerenciadas em um único lugar, facilitando a proteção de dados.
- **Facilidade de Configuração:** Alterações nos ambientes (desenvolvimento, teste, produção) são simplificadas.

### Escolha de tecnologias

- **Golang:** Escolhido pela performance, simplicidade e forte suporte à concorrência.
- **PostgreSQL com `pgxpool`:** Proporciona uma conexão eficiente e recursos avançados para interagir com o banco de dados.
- **Docker Compose:** Utilizado para orquestrar serviços, garantindo consistência entre ambientes e facilitando a implantação.

---

## Referências aos arquivos por camada

- **Camada de Apresentação:**
    - `internal/adapter/http/handler/`
        - `user_handler.go`
    - `internal/adapter/http/router.go`
- **Camada de Aplicação:**
    - `internal/application/service/`
        - `user_service.go`
    - `internal/application/dto/`
        - `user_dto.go`
- **Camada de Domínio:**
    - `internal/domain/entity/`
        - `user.go`
    - `internal/domain/repository.go`
- **Camada de Infraestrutura:**
    - `internal/adapter/repository/`
        - `user_repository.go`
    - `internal/config/config.go`
    - `internal/di/container.go`
    - `internal/di/database.go`

---

## Conclusão

A arquitetura implementada no projeto **FastFood Golang** foi cuidadosamente planejada para promover boas práticas de desenvolvimento de software, como modularidade, separação de preocupações e independência tecnológica. Ao seguir os princípios da Arquitetura Limpa e Hexagonal, garantimos que a aplicação seja:

- **Escalável:** Facilmente expansível para adicionar novos recursos ou módulos.
- **Manutenível:** Simples de entender e modificar, reduzindo o tempo de desenvolvimento e custos.
- **Testável:** Com componentes desacoplados, os testes unitários e de integração são mais fáceis de implementar.
- **Flexível:** Capaz de se adaptar a mudanças tecnológicas sem grandes refatorações.

Esta arquitetura não apenas atende aos requisitos acadêmicos da disciplina de Arquitetura de Software, mas também prepara o terreno para projetos profissionais de alta qualidade.

---

**Nota Final:** A compreensão das interações entre as camadas e a racionalização por trás das decisões de arquitetura é fundamental para qualquer desenvolvedor que deseje contribuir para o projeto ou aplicar conceitos semelhantes em outros contextos. Esperamos que esta explicação detalhada facilite esse entendimento e sirva como um guia para futuras implementações.
