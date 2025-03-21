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

Para mais detalhes, consulte o arquivo [INFRA.md](./../infra.md).

**Detalhes e decisões:**

- **Uso do `pgxpool`:** Optamos pelo driver `pgx` para melhor performance e recursos avançados na interação com o PostgreSQL.
- **Injeção de Dependências com Uber Dig:** Facilita o gerenciamento de dependências complexas e aumenta a testabilidade da aplicação.
- **Centralização das Configurações:** Permite alterar facilmente parâmetros de configuração sem modificar o código-fonte.

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

## Conclusão

A arquitetura implementada no projeto **FastFood Golang** foi cuidadosamente planejada para promover boas práticas de desenvolvimento de software, como modularidade, separação de preocupações e independência tecnológica. Ao seguir os princípios da Arquitetura Limpa e Hexagonal, garantimos que a aplicação seja:

- **Escalável:** Facilmente expansível para adicionar novos recursos ou módulos.
- **Manutenível:** Simples de entender e modificar, reduzindo o tempo de desenvolvimento e custos.
- **Testável:** Com componentes desacoplados, os testes unitários e de integração são mais fáceis de implementar.
- **Flexível:** Capaz de se adaptar a mudanças tecnológicas sem grandes refatorações.

Esta arquitetura não apenas atende aos requisitos acadêmicos da disciplina de Arquitetura de Software, mas também prepara o terreno para projetos profissionais de alta qualidade.

---

**Nota Final:** A compreensão das interações entre as camadas e a racionalização por trás das decisões de arquitetura é fundamental para qualquer desenvolvedor que deseje contribuir para o projeto ou aplicar conceitos semelhantes em outros contextos. Esperamos que esta explicação detalhada facilite esse entendimento e sirva como um guia para futuras implementações.
