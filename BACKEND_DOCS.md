# Documentación del Backend de Ultra Monitor

Esta documentación proporciona una guía detallada sobre la arquitectura, estructura y funcionamiento del backend del sistema Ultra Monitor.

## 1. Introducción

Ultra Monitor es un sistema de monitoreo y gestión de red diseñado para supervisar el tráfico, el estado de los dispositivos (OLTs, ONTs) y generar reportes detallados. El backend está construido en **Go (Golang)**, priorizando el rendimiento y la concurrencia.

### Tecnologías Principales
- **Lenguaje**: Go
- **Framework Web**: [Fiber v3](https://gofiber.io/)
- **Base de Datos Relacional**: PostgreSQL (usando `sqlx`)
- **Base de Datos NoSQL**: MongoDB
- **Caché**: Redis
- **Métricas**: Prometheus
- **Protocolos**: HTTP/REST, SNMP

## 2. Arquitectura del Sistema

El proyecto sigue una **Clean Architecture** (Arquitectura Limpia) o Arquitectura en Capas, separando claramente las responsabilidades.

### Diagrama de Arquitectura

```mermaid
graph TD
    Client[Cliente Web / CLI] --> Router[Router / Middleware]
    
    subgraph "Capa de Presentación (Transport Layer)"
        Router --> Handler[Handlers / Controllers]
    end
    
    subgraph "Capa de Negocio (Domain Layer)"
        Handler --> UseCase[Use Cases / Services]
        UseCase --> Entity[Entities / Models]
    end
    
    subgraph "Capa de Datos (Data Layer)"
        UseCase --> Repository[Repositories]
        Repository --> DB[(PostgreSQL)]
        Repository --> Mongo[(MongoDB)]
        Repository --> Redis[(Redis Cache)]
        Repository --> Prom[Prometheus Client]
    end
```

### Descripción de Capas

1.  **Handler (Controladores)**: Maneja las peticiones HTTP, valida la entrada y llama a los casos de uso apropiados. No contiene lógica de negocio compleja.
2.  **UseCase (Casos de Uso)**: Contiene la lógica de negocio pura. Orquesta el flujo de datos entre las entidades y los repositorios.
3.  **Repository (Repositorios)**: Abstrae el acceso a datos. Se encarga de interactuar directamente con la base de datos, caché o servicios externos como Prometheus.
4.  **Entity (Entidades)**: Define los modelos de dominio y estructuras de datos fundamentales.

## 3. Estructura del Proyecto

La estructura de directorios sigue las convenciones estándar de Go:

| Directorio | Descripción |
| :--- | :--- |
| `cmd/` | Puntos de entrada de la aplicación. Contiene `main.go`. |
| `entity/` | Definiciones de estructuras de datos (structs) que representan los objetos del dominio. |
| `handler/` | Controladores que manejan las peticiones HTTP y respuestas. |
| `routes/` | Configuración de rutas y agrupación de endpoints. |
| `usecase/` | Lógica de negocio de la aplicación. |
| `repository/` | Implementación del acceso a datos (SQL, Mongo, Redis). |
| `middleware/` | Middleware para logging, autenticación, rate limiting, etc. |
| `internal/` | Código privado de la aplicación (configuración de DB, utilidades, clientes externos). |
| `pkg/` | Código que podría ser reutilizable en otros proyectos (si aplica). |

## 4. Flujo de Trabajo y Datos

A continuación se muestra el flujo típico de una petición para obtener datos de tráfico:

```mermaid
sequenceDiagram
    participant Client as Cliente
    participant Router as Fiber Router
    participant Handler as Traffic Handler
    participant UseCase as Traffic UseCase
    participant Repo as Traffic Repository
    participant DB as PostgreSQL/Redis

    Client->>Router: GET /api/traffic/total
    Router->>Handler: GetNationalTraffic()
    Handler->>UseCase: GetNationalTraffic(dates)
    UseCase->>Repo: GetNationalTraffic(dates)
    Repo->>DB: Query SQL / Get Cache
    DB-->>Repo: Datos
    Repo-->>UseCase: Structs de Dominio
    UseCase-->>Handler: Respuesta Procesada
    Handler-->>Client: JSON Response
```

## 5. Modos de Ejecución (CLI)

El backend no solo funciona como un servidor web, sino que también incluye herramientas de línea de comandos para tareas de mantenimiento y recolección de datos. El archivo `cmd/main.go` gestiona estos modos.

```mermaid
graph TD
    Start(Ejecutar Binario) --> CheckArgs{Argumentos?}
    
    CheckArgs -->|server| Server[Iniciar Servidor Web API]
    CheckArgs -->|scan| Scan[Escanear Dispositivos Prometheus]
    CheckArgs -->|traffic| Traffic[Actualizar Resumen de Tráfico]
    CheckArgs -->|traffic-ont| TrafficONT[Monitoreo Distribuido de ONTs]
    CheckArgs -->|interface-bandwidth| Bandwidth[Actualizar Ancho de Banda Interfaces]
    
    Server --> Fiber[Fiber App Listen]
    Scan --> PromScan[Prometheus Device Scan] --> DBUpdate[Actualizar DB]
    Traffic --> CalcTraffic[Calcular Tráfico Diario] --> DBUpdate
    TrafficONT --> MonitorLoop[Bucle de Monitoreo por ONT]
    Bandwidth --> MongoUpdate[Actualizar MongoDB]
```

### Comandos Disponibles

-   **`server`**: Inicia el servidor API REST.
    ```bash
    ./app server
    ```
-   **`scan`**: Escanea la red en busca de dispositivos a través de Prometheus y actualiza el inventario.
    ```bash
    ./app scan
    ```
-   **`traffic <YYYYMMDD>`**: Procesa y resume el tráfico para una fecha específica.
    ```bash
    ./app traffic 20241125
    ```
-   **`traffic-ont <community>`**: Inicia un proceso de larga duración para monitorear el tráfico de ONTs individuales de forma distribuida.
    ```bash
    ./app traffic-ont public
    ```
-   **`interface-bandwidth`**: Actualiza las estadísticas de ancho de banda de las interfaces.
    ```bash
    ./app interface-bandwidth
    ```

## 6. Componentes Clave

### Entidades (`entity/`)
Definen la estructura de los datos. Ejemplo: `Traffic`, `User`, `Ont`, `Prometheus`.

### Repositorios (`repository/`)
Encapsulan las consultas SQL y comandos de MongoDB.
-   **Patrón**: Cada repositorio tiene métodos específicos para su dominio (e.g., `traffic.repository.go` maneja todo lo relacionado con tablas de tráfico).
-   **Tecnologías**: `sqlx` para PostgreSQL, driver oficial para MongoDB.

### Casos de Uso (`usecase/`)
Aquí reside la inteligencia del sistema.
-   **Validación**: Verifica reglas de negocio antes de persistir datos.
-   **Orquestación**: Puede llamar a múltiples repositorios (e.g., obtener datos de Redis, si no están, ir a SQL).

### Handlers (`handler/`)
La capa de entrada HTTP.
-   **Parámetros**: Parsea query params y body.
-   **Respuesta**: Formatea la respuesta en JSON estándar.

## 7. Integraciones

### Prometheus
Se utiliza para la recolección de métricas de red. El cliente interno (`internal/prometheus`) permite consultar el estado de dispositivos y métricas de interfaces.

### MongoDB
MongoDB se utiliza exclusivamente como una fuente de datos externa para la sincronización. Otra plataforma almacena métricas de ancho de banda en esta base de datos, y Ultra Monitor realiza consultas periódicas para extraer estos datos y actualizar su base de datos principal (PostgreSQL). No se utiliza para persistencia directa de la lógica de negocio del backend.

### Redis
Utilizado para caché de consultas frecuentes y almacenamiento temporal de sesiones o estados.

## 8. Configuración

El sistema se configura mediante variables de entorno. Asegúrese de definir las siguientes variables en su archivo `.env` o en el entorno de despliegue:

| Variable | Descripción | Ejemplo |
| :--- | :--- | :--- |
| `PORT` | Puerto donde escuchará el servidor. | `3000` |
| `ENVIROMENT` | Entorno de ejecución (`development` o `production`). | `development` |
| `POSTGRES_URI` | Cadena de conexión a PostgreSQL. | `postgres://user:pass@localhost:5432/db` |
| `MONGODB_URI` | URI de conexión a MongoDB. | `mongodb://localhost:27017` |
| `REDIS_URI` | URI de conexión a Redis. | `redis://localhost:6379` |
| `PROMETHEUS_URL` | URL del servidor Prometheus. | `http://localhost:9090` |
| `AUTH_SECRET_KEY` | Clave secreta para firmar tokens JWT. | `supersecretkey` |
| `REPORTS_DIRECTORY` | Directorio donde se guardarán los reportes generados. | `./reports` |
| `CORS_ALLOW_ORIGIN` | (Prod) Orígenes permitidos para CORS. | `https://mi-app.com` |
| `WEB_APP_DIRECTORY` | (Prod) Directorio de los archivos estáticos del frontend. | `./web/dist` |

