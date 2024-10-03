# Explicação Detalhada das Camadas e Interações no Projeto

O projeto **FastFood Golang** foi estruturado seguindo os princípios da **Arquitetura Limpa** e **Arquitetura Hexagonal**, com o objetivo de criar um sistema modular, escalável e de fácil manutenção. Abaixo, detalhamos cada camada, suas responsabilidades, como elas interagem entre si e referências aos arquivos relevantes.

## Camadas da Arquitetura

1. [Camada de Apresentação (Presentation Layer)](#1-camada-de-apresentação-presentation-layer)
2. [Camada de Aplicação (Application Layer)](#2-camada-de-aplicação-application-layer)
3. [Camada de Domínio (Domain Layer)](#3-camada-de-domínio-domain-layer)
4. [Camada de Infraestrutura (Infrastructure Layer)](#4-camada-de-infraestrutura-infrastructure-layer)

---

### 1. Camada de Apresentação (Presentation Layer)

**Responsabilidade:**

- Gerenciar a comunicação com o mundo externo (clientes, navegadores, etc.).
- Processar solicitações HTTP, validar entradas e formatar respostas.
- Delegar a lógica de negócios para a camada de aplicação.

**Interação:**

- Recebe solicitações HTTP dos clientes.
- Invoca serviços da camada de aplicação para processar a lógica de negócios.
- Retorna respostas HTTP apropriadas aos clientes.

**Arquivos Relevantes:**

- **Manipuladores HTTP (Handlers):**
    - `internal/adapter/http/handler/user_handler.go`: Contém a implementação do `UserHandler`, que processa as solicitações relacionadas a usuários.
- **Roteador HTTP:**
    - `internal/adapter/http/router.go`: Configura as rotas da API e registra os manipuladores correspondentes.

**Detalhes e Decisões:**

- **Uso do Gin Framework:** Escolhemos o Gin devido à sua eficiência e facilidade de uso para criar APIs RESTful.
- **Responsabilidade Limitada dos Handlers:** Os manipuladores focam em processar solicitações e respostas, delegando a lógica de negócios para a camada de aplicação.

---

### 2. Camada de Aplicação (Application Layer)

**Responsabilidade:**

- Implementar a lógica de negócios da aplicação.
- Orquestrar operações entre a camada de apresentação e a camada de domínio.
- Garantir que as regras de negócios sejam aplicadas corretamente.

**Interação:**

- Recebe chamadas dos manipuladores (handlers) da camada de apresentação.
- Utiliza entidades e interfaces da camada de domínio para processar dados.
- Chama repositórios através das interfaces definidas na camada de domínio.

**Arquivos Relevantes:**

- **Serviços de Aplicação:**
    - `internal/application/service/user_service.go`: Implementa a lógica de negócios relacionada a usuários, como criação, leitura, atualização e exclusão.
- **Data Transfer Objects (DTOs):**
    - `internal/application/dto/user_dto.go`: Define estruturas para transferência de dados entre camadas, evitando expor diretamente as entidades do domínio.

**Detalhes e Decisões:**

- **Isolamento da Lógica de Negócios:** Centralizamos a lógica aqui para facilitar testes e manutenções futuras.
- **Uso de DTOs:** Facilita a validação e transformação de dados entre as camadas, promovendo a segurança e integridade dos dados.

---

### 3. Camada de Domínio (Domain Layer)

**Responsabilidade:**

- Representar o núcleo da aplicação com as regras de negócio fundamentais.
- Definir entidades e interfaces que modelam conceitos do domínio.
- Permanecer independente de detalhes de implementação ou frameworks externos.

**Interação:**

- As entidades e interfaces são usadas pela camada de aplicação para processar dados.
- As interfaces de repositório definidas aqui são implementadas pela camada de infraestrutura.

**Arquivos Relevantes:**

- **Entidades de Domínio:**
    - `internal/domain/entity/user.go`: Define a estrutura da entidade `User`, representando um usuário no sistema.
- **Interfaces de Repositório:**
    - `internal/domain/repository.go`: Declara a interface `UserRepository`, especificando os métodos que devem ser implementados para manipulação de usuários.

**Detalhes e Decisões:**

- **Independência Tecnológica:** Ao não depender de frameworks ou pacotes externos, a camada de domínio permanece flexível e adaptável a mudanças.
- **Definição de Interfaces:** As interfaces permitem que diferentes implementações sejam usadas sem alterar a lógica de negócios, facilitando testes e trocas de tecnologia.

---

### 4. Camada de Infraestrutura (Infrastructure Layer)

**Responsabilidade:**

- Fornecer implementações concretas para as interfaces definidas na camada de domínio.
- Lidar com detalhes técnicos, como acesso a bancos de dados, serviços externos, etc.

**Interação:**

- Implementa as interfaces de repositório, interagindo diretamente com o banco de dados.
- É chamada pela camada de aplicação através das interfaces do domínio.

**Arquivos Relevantes:**

- **Implementação dos Repositórios:**
    - `internal/adapter/repository/user_repository.go`: Implementa `UserRepository`, realizando operações de banco de dados para a entidade `User`.
- **Configurações da Aplicação:**
    - `internal/config/config.go`: Carrega e gerencia as configurações do aplicativo, como variáveis de ambiente e strings de conexão.
- **Injeção de Dependências:**
    - `internal/di/container.go`: Configura o container de injeção de dependências usando o Uber Dig, registrando todos os provedores necessários.
- **Conexão com o Banco de Dados:**
    - `internal/di/database.go`: Estabelece a conexão com o PostgreSQL utilizando `pgxpool`.

**Detalhes e Decisões:**

- **Uso do `pgxpool`:** Optamos pelo driver `pgx` para melhor performance e recursos avançados na interação com o PostgreSQL.
- **Injeção de Dependências com Uber Dig:** Facilita o gerenciamento de dependências complexas e aumenta a testabilidade da aplicação.
- **Centralização das Configurações:** Permite alterar facilmente parâmetros de configuração sem modificar o código-fonte.

---

## Como as Camadas Interagem

### Fluxo de uma Operação: "Get All Users"

1. **Solicitação HTTP:**
    - O cliente envia uma requisição HTTP GET para `/api/v1/users`.

2. **Camada de Apresentação:**
    - **Router (`router.go`):** Direciona a requisição para o `UserHandler`.
    - **Handler (`user_handler.go`):**
        - Recebe a requisição.
        - Realiza validações iniciais, se necessário.
        - Chama o método `GetAllUsers()` do `UserService`.

3. **Camada de Aplicação:**
    - **Service (`user_service.go`):**
        - Aplica regras de negócio (por exemplo, filtragem, ordenação).
        - Chama o método `GetAll(ctx)` do repositório de usuários através da interface `UserRepository`.

4. **Camada de Domínio:**
    - **Interface (`repository.go`):**
        - Define o contrato para `GetAll(ctx)` que deve ser implementado.
    - **Entidade (`user.go`):**
        - Estrutura de dados que representa um usuário.

5. **Camada de Infraestrutura:**
    - **Repositório (`user_repository.go`):**
        - Implementa `GetAll(ctx)`, executando uma consulta SQL no banco de dados.
        - Utiliza o pool de conexões `pgxpool` para interação com o PostgreSQL.
    - **Banco de Dados:**
        - Executa a consulta e retorna os resultados para o repositório.

6. **Resposta:**
    - **Repositório:** Retorna a lista de usuários para o serviço.
    - **Serviço:** Pode aplicar transformações adicionais nos dados.
    - **Handler:** Formata a resposta (JSON) e envia de volta ao cliente com o status HTTP adequado.

### Interações Chave:

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

### Uso de Interfaces e Abstrações

- **Flexibilidade:** Permite trocar implementações (por exemplo, substituir o banco de dados) sem alterar as camadas superiores.
- **Isolamento da Lógica de Negócios:** A lógica de negócios não depende de detalhes de infraestrutura, seguindo o princípio da inversão de dependência.

### Gerenciamento Centralizado de Configurações

- **Segurança:** Variáveis sensíveis são gerenciadas em um único lugar, facilitando a proteção de dados.
- **Facilidade de Configuração:** Alterações nos ambientes (desenvolvimento, teste, produção) são simplificadas.

### Escolha de Tecnologias

- **Golang:** Escolhido pela performance, simplicidade e forte suporte à concorrência.
- **PostgreSQL com `pgxpool`:** Proporciona uma conexão eficiente e recursos avançados para interagir com o banco de dados.
- **Docker Compose:** Utilizado para orquestrar serviços, garantindo consistência entre ambientes e facilitando a implantação.

---

## Referências aos Arquivos por Camada

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