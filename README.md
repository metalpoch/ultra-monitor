# olt-blueprint

Application capable of displaying and carrying out traffic coming from OLT

## Getting Started

To initialize the application as a developer, follow these steps:

### 1. Clone the Repository

First, clone the repository to your local machine:

```bash
git clone https://github.com/metalpoch/olt-blueprint.git
```

### 2. Start Docker Containers

Run the following command in your terminal to start the Docker containers:

```bash
make container-run
```

### 3. Run the Auth Module (for example)

Run the following command in your terminal to start the Docker containers:

```bash
make dev-auth
```

This will start the authentication service in development mode.

#### Configuration

The application uses a configuration file for development purposes. You can find the credentials for the Docker images in the config.develop.json file. Important: This file contains credentials intended for testing and development only. Do not use these credentials in a production environment.
