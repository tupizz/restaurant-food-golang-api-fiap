# Clean Architecture no Projeto

Este projeto segue os princípios da **Clean Architecture**, garantindo modularidade, testabilidade e desacoplamento. Abaixo está uma explicação detalhada da estrutura de pastas e como cada uma delas se relaciona com os conceitos da Clean Architecture.

---

## Estrutura de Diretórios

```plaintext
/
├── cmd/
├── database/
│   ├── migrations/
│   ├── queries/
|   └── sqlc/
├── docs/
├── internal/
│   ├── adapters/
│   │   ├── db/
│   │   │   └── repository/
│   │   ├── gateways/
│   │   │   └── payment/
│   │   └── http/
│   │       |── handler/
|   |       └── middleware/
|   ├── config/
│   ├── core/
│   │   ├── domain/
│   │   │   ├── entities/
│   │   │   ├── error/
│   │   │   └── validator/
│   │   ├── usecase/
│   │   │   ├── dto/
|   |   |   ├── mappers/
│   │   │   └── ports/
│   ├── di/
|   └── shared/
```

---

## Detalhes de Cada Diretório

### 1. **cmd/**

Contém o ponto de entrada da aplicação, como arquivos `main.go` para inicializar o projeto.
- Configura dependências e inicializa servidores HTTP ou outros serviços.

### 2. **database/**

Contém scripts SQL e outros arquivos relacionados ao banco de dados:

- **migrations/**: Scripts para migrações do banco de dados (ex.: criar tabelas, adicionar colunas).
- **queries/**: Arquivos SQL usados para consultas específicas no banco de dados.
- **sqlc/**: Arquivos autogerados pelo pacote `sqlc` usados para consultas específicas no banco de dados.

### 3. **docs/**

Armazena documentação do projeto, como arquivos gerados automaticamente para a API (e.g., Swagger) ou manuais escritos.

### 4. **internal/**

O núcleo da aplicação está dentro de `internal`, dividido em várias subpastas:

#### a) **adapters/**

Responsável por conectar o núcleo da aplicação com frameworks, bibliotecas externas e tecnologias específicas. Contém as implementações específicas de entrada e saída.

- **db/repository/**: Implementações específicas dos repositórios que interagem com o banco de dados. Esses repositórios são responsáveis por persistir e recuperar dados.

- **gateways/payment/**: Gateways são responsáveis por integrar o sistema com recursos externos, encapsulando a lógica de comunicação e protegendo o domínio de dependências externas, neste caso foram retratadas as comunicações com Gateways de pagmento.

- **http/handler/**: Contém os manipuladores HTTP (endpoints da API). Eles recebem requisições, validam dados e chamam os casos de uso.

- **http/middleware/**: Contém os manipuladores HTTP que rodam antes dos handlers (um bo exemplo seria uma camada de autenticação). Eles recebem requisições, validam dados e permitem que as chamadas cheguem aos handlers enriquecidas.

#### b) **config/**

Configurações da aplicação, como variáveis de ambiente, configuração de banco de dados, etc.

#### c) **core/**

O núcleo central da aplicação, seguindo os conceitos de Clean Architecture:

- **domain/**:
  - **entities/**: Define as entidades principais do domínio. Essas entidades representam as regras de negócio fundamentais e geralmente são objetos simples.
  - **error/**: Contém erros de domínio que representam situações específicas que podem ocorrer nas regras de negócio.
  - **validator/**: Validações específicas para as entidades ou regras do domínio.

- **usecase/**:
  - **dto/**: Define os Data Transfer Objects (DTOs), usados para transferir dados entre as camadas (ex.: entrada e saída de dados nos casos de uso).
  - **mappers/**: Define as conversões/traduções entre as DTOs e entidades.
  - **ports/**: Define as interfaces dos casos de uso e repositórios, permitindo que o núcleo da aplicação seja desacoplado das implementações externas.

#### d) **di/**

Contém a configuração de injeção de dependências. Aqui, as implementações concretas são conectadas às interfaces do núcleo, configurando o grafo de dependências da aplicação.

#### e) **shared/**

Código compartilhado comum entre os pacotes, como um simples loop, ou algo assim, mas que não afeta em nada as entidades conehcidas pelo app.

---

## Fluxo de Dados na Arquitetura

1. **Requisição HTTP**:
   - Chega no `handler` dentro de `adapters/http/handler`. O handler valida os dados da requisição e chama os métodos dos casos de uso.

2. **Casos de Uso**:
   - O caso de uso, localizado em `core/usecase`, implementa a lógica de negócio e chama as portas definidas para acessar repositórios ou outras dependências.

3. **Repositórios**:
   - As portas do repositório em `core/usecase/ports` são implementadas em `adapters/db/repository`. Essas implementações interagem com o banco de dados ou outro sistema de persistência.

4. **Resposta**:
   - O caso de uso retorna os dados processados ao handler, que formata a resposta e a envia ao cliente.

---

## Princípios Seguidos

1. **Desacoplamento**:
   - As camadas não têm dependência direta entre si, mas comunicam-se por meio de interfaces.

2. **Independência do Framework**:
   - O núcleo da aplicação não depende de frameworks. Por exemplo, a lógica de negócio não sabe que está sendo chamada por um servidor HTTP.

3. **Testabilidade**:
   - Cada camada pode ser testada isoladamente, graças ao uso de interfaces e injeção de dependências.

4. **Centralidade do Domínio**:
   - As regras de negócio estão isoladas no núcleo (`core/domain`) e não dependem de implementações externas.

---

## Exemplos de Uso

### Criar um Cliente

1. Requisição chega no endpoint HTTP definido no `handler`.
2. O `handler` chama o método correspondente do caso de uso (`usecase.ClientUseCase`).
3. O caso de uso:
   - Valida os dados.
   - Cria a entidade do domínio.
   - Salva no banco de dados via repositório.
4. A resposta é enviada de volta ao cliente.

---

## Conclusão

Essa estrutura foi projetada para oferecer flexibilidade e organização, permitindo que a aplicação evolua de maneira sustentável. A Clean Architecture garante que mudanças no banco de dados, frameworks ou tecnologias externas não afetem o núcleo da aplicação.
