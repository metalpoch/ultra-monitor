# 📐 OLT Blueprint 
Application for traffic monitoring in OLTs (Optical Line Terminals) using SNMP (Simple Network Management Protocol) for traffic analysis by OLT, port, clients, or geographical location.

## ⚙️ App workflow
```mermaid
graph LR
X[(OLT 1)] -- SNMP --> M(Measurement*)
Y[(OLT 2)] -- SNMP --> M(Measurement*)
Z[(OLT n)] -- SNMP --> M(Measurement*)
M --> D[(Database)]
D --> S(Smart*) 
D --> R(Report*)
D --> T(Traffic*)
T --> A(Auth*)
R --> A(Auth*)
S --> A(Auth*)
A --> C(Client)
```
This diagram shows how measurement data from various OLTs is collected and stored in a database, then used to generate traffic and reports, which are accessible to clients through an authentication system. Below are the components and their interactions:

1.  **OLT (Optical Line Terminal)**:
    -   **OLT's**: Represent different OLTs sending measurement data.
    -   **SNMP Protocol**: Measurement data is transmitted from the OLTs to the measurement node using the SNMP (Simple Network Management Protocol).
2.  **Measurement**: Node representing the module that receives and centralizes the measurements sent by the OLTs.
3.  **Database**: Node where the module received measurements are stored for further processing.
4.  **Data Processing**:
    -   **Traffic**: Node representing the module the traffic generated from the stored data.
    -   **Report**: Node representing the module the generation of reports based on the stored data.
    - **Smart**: Node representing the module that processing of intelligent data from stored information.
5.  **Authentication**:
    - **Auth**: Node representing the module that controls access to both the traffic and the reports, ensuring that only authenticated users can access the information.
6.  **Client**:  Node representing the end client who accesses the authenticated data.

## 🛰️ Technologies used

- <img src="https://go.dev/favicon.ico" width=30 /> [Go]: concurrent, compiled, efficient, simple and scalable language
- <img src="https://www.python.org/favicon.ico" width=30 /> [Python]: interpreted language, versatile, easy, powerful and popular
- <img src="https://astro.build/favicon.svg" width=30 /> [Astro]: framework for creating fast and efficient websites
- <img src="https://www.postgresql.org/favicon.ico" width=30 /> [PostgreSQL]: robust, flexible and free relational database
- <img src="https://ollama.com/public/icon-32x32.png" width=30 /> [Ollama]: platform for running and customizing language models
- <img src="https://www.docker.com/favicon.ico" width=30 /> [Docker]: is an efficient and portable application container platform

## 🚀 Getting Started

To initialize the application as a developer, follow these steps:
```bash
# clone this repository
git clone https://github.com/metalpoch/olt-blueprint.git
```

```sql
-- create the database
CREATE DATABASE olt;
```
create the configuration file called `config.json` in the root of the project
```json
{
  "db_uri": "postgresql://xxx:yyy@zzz:5432/olt",
  "cache_uri": "redis://xxx:yyy@zzz:6379",
  "secret_key": "supersecrettoken",
  "telegram_bot_token_id": "xx:yyy-zzz",
  "telegram_chat_id": "xxx"
}
```

run any module
```bash
# run auth module
make dev-auth

# run traffic module
make dev-traffic

# run report module
make dev-report

# run traffic module
make dev-traffic

# run the CLI measurement module
CONFIG_JSON=./config.json go run ./measurement/cmd/main.go
```
## Licencia
[MIT]

[MIT]:https://github.com/metalpoch/olt-blueprint/blob/main/LICENSE
[repository]: https://github.com/metalpoch/olt-blueprint/
[Go]:https://go.dev/
[PostgreSQL]:https://www.postgresql.org/
[Python]:https://www.python.org/
[Ollama]:https://ollama.com/
[Astro]:https://astro.build
[Docker]:https://www.docker.com/
