### High-level Architecture
This diagram shows the overall structure of your application's layers and their interactions.

Explanation:

- Handlers / Controllers (A) in the Presentation Layer call the Services (B) in the Application Layer.
- Services (B) interact with Entities and Repository Interfaces (C) in the Domain Layer.
- The Domain Layer (C) defines the Repository Interfaces, which are implemented by the Repositories (D) in the Infrastructure Layer.
- The Infrastructure Layer (D) communicates back to the Domain Layer (C) by fulfilling the repository contracts.

```mermaid
flowchart TD
    %% Define the layers
    subgraph Presentation Layer
        A[Handlers / Controllers]
    end

    subgraph Application Layer
        B[Services]
    end

    subgraph Domain Layer
        C[Entities &<br>Repository Interfaces]
    end

    subgraph Infrastructure Layer
        D[Repositories /<br>Database Access]
    end

    %% Define the interactions
    A --> B
    B --> C
    C <--> D
```

### Sequence for "Get all Users"

```mermaid
sequenceDiagram
    participant Client
    participant Handler as UserHandler
    participant Service as UserService
    participant Repository as UserRepository
    participant Database as PostgreSQL

    Client->>+Handler: HTTP GET /api/v1/users
    Handler->>+Service: GetAllUsers()
    Service->>+Repository: GetAll(ctx)
    Repository->>+Database: Execute SQL Query
    Database-->>-Repository: Query Results
    Repository-->>-Service: List of Users
    Service-->>-Handler: List of Users
    Handler-->>-Client: HTTP Response with Users
```

### Diagram with Dependecy Injection

```mermaid
flowchart LR
    subgraph DI_Container
        Config[Config Loader]
        DBPool["Database Connection<br>(pgxpool.Pool)"]
        RepoImpl[UserRepository<br>Implementation]
        Service[UserService]
        Handler[UserHandler]
        Router[HTTP Router]
    end

    Config --> DBPool
    Config --> RepoImpl
    DBPool --> RepoImpl
    RepoImpl --> Service
    Service --> Handler
    Handler --> Router
    Router -->|Registers| Routes

%% External Dependencies
    subgraph External
        Database[(PostgreSQL Database)]
        Client[HTTP Client]
    end

    RepoImpl --> Database
    Client -->|Sends Request| Router

```

### Package Structure Diagram

```mermaid
flowchart TB
    cmd[cmd/your-app]
    di[internal/di]
    config[internal/config]
    handler[internal/adapter/http/handler]
    router[internal/adapter/http/router]
    service[internal/application/service]
    entity[internal/domain/entity]
    repositoryIntf[internal/domain/repository.go]
    repositoryImpl[internal/adapter/repository]

    cmd --> di
    cmd --> router
    di --> config
    di --> handler
    di --> service
    di --> repositoryImpl
    handler --> service
    service --> repositoryIntf
    repositoryImpl --> repositoryIntf
    repositoryImpl --> entity
    service --> entity
    handler --> entity
```

### Hexagonal Architecture Diagram

```mermaid
flowchart TB
    %% Core Domain
    subgraph Domain
        Direction[Business Logic]
        EntityModels[Entities]
        Ports["Ports (Interfaces)"]
    end

    %% Adapters
    subgraph Adapters
        Inbound[Inbound Adapters]
        Outbound[Outbound Adapters]
    end

    %% External Systems
    Client[Clients]
    Database[(Database)]

    %% Interactions
    Client --> Inbound
    Inbound -->|Calls| Ports
    Ports -->|Calls| Outbound
    Outbound --> Database

```

### Clean Architexture Diagrama

```
TO-DO
```
