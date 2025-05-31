# 📐 Ultra Monitor 
Ultra Monitor is a platform for traffic monitoring and analysis in OLTs (Optical Line Terminals) using SNMP (Simple Network Management Protocol). It supports traffic analysis by OLT, port, client, or geographical location.

## 🗂️ Structure
- **cmd/server/**: Main backend (REST API, authentication, reports, traffic, etc.) in Go.
- **cmd/cli/**: CLI tool for administration, scanning, and technical operations.
- **controller/**, **handler/**, **routes/**, **usecase/**: Shared domain logic, controllers, and routes.
- **client/**: React frontend (coming soon).
- **smart/**: AI module in Python (coming soon).

## ⚙️ App workflow
```mermaid
graph LR
ONT11[(ONT)] --> OLT1[(OLT)]
ONT12[(ONT)] --> OLT1[(OLT)]

ONT2n[(ONT)] --> OLT2[(OLT n)]

OLT1[(OLT)] -- SNMP --> M(Measurement*)

OLT2[(OLT)] -- SNMP --> M(Measurement*)

M --> D[(Database)]
M --> S(Server*) 
M --> R[(Cache)]

D --> S
R --> S
S --> UI(Client)
```

This diagram shows how measurement data from various OLTs is collected and stored in a database, then used to generate traffic and reports, which are accessible to clients through an authentication system.


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
git clone https://github.com/metalpoch/ultra-monitor.git
```

```sql
-- create the database
CREATE DATABASE ultramonitor;
```

Define the following environment variables
```bash
POSTGRES_URI="postgresql://user:password@host:5432/ultramonitor"
AUTH_SECRET_KEY="secret-super-secure!"
REDIS_URI="redis://user:password@host:6379"
```

Execute docker-compose for run server or cli

## 📦 Contributing

Want to contribute? Awesome! Please check the contribution guidelines and open a Pull Request.

---


## ⚖️ License
[MIT]

[MIT]:https://github.com/metalpoch/olt-blueprint/blob/main/LICENSE
[repository]: https://github.com/metalpoch/olt-blueprint/
[Go]:https://go.dev/
[PostgreSQL]:https://www.postgresql.org/
[Python]:https://www.python.org/
[Ollama]:https://ollama.com/
[Astro]:https://astro.build/
[Docker]:https://www.docker.com/
