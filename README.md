# Distributed Manufacturing Pipeline Simulation System

## Problem Statement
Design and implement a **distributed manufacturing pipeline simulation system** that runs on **Kubernetes**, providing **real-time control and monitoring** capabilities through a **web interface**. The system follows modern **microservices architecture principles**, **async processing**, and **cloud-native design patterns**.

## Project Overview
This project simulates a real-world **manufacturing pipeline with multiple stages**. It is designed to be **scalable, resilient**, and capable of **real-time monitoring**. Different **microservices** represent various **production stages** and are deployed using a **Kubernetes-based cloud environment**.

### **Key Features**
✅ Real-time control and monitoring
✅ Asynchronous processing for efficiency
✅ Microservices-based architecture
✅ Cloud-native deployment with Kubernetes
✅ Integration with **RabbitMQ/NATS** for messaging

## **Architecture**
The system follows a **microservices architecture**, where each component is modular and serves a specific purpose.

### **Frontend (React + Material UI)**
- **Authentication & Authorization:** User registration, login, and role-based access.
- **Pipeline Management:**
  - Create pipelines with configurable stages.
  - View all created pipelines with status updates.
  - Execute and cancel pipelines.
- **Real-Time Monitoring:**
  - WebSocket integration for live pipeline status updates.
  - Dashboard visualizations.

### **Backend (Go)**
- **User Authentication:** Uses **Supabase authentication** with **JWT-based session management**.
- **Pipeline Execution:**
  - Execution model managed by **ParallelPipelineOrchestrator**.
  - Supports **parallel and sequential execution**.
  - **Rollback support** for failure handling.
- **API Services:**
  - **REST API** for frontend interactions.
  - **gRPC API** for CLI (`democtl`).
- **Database Interaction:** Uses **Supabase (PostgreSQL)** for pipeline, stage, and execution storage.
- **Messaging System:** Uses **RabbitMQ/NATS** for async processing.

### **WebSockets (Real-Time Updates)**
- Maintains active client connections.
- Broadcasts events when a stage status changes.
- Eliminates polling for better performance.

## **Technology Stack**
- **Frontend:** React + Material UI + WebSockets
- **Backend:** Go + Gin (REST API) + gRPC (CLI)
- **Database:** Supabase (PostgreSQL)
- **Messaging:** RabbitMQ/NATS (async processing)
- **Orchestration:** Kubernetes + Docker

## **Pipeline Execution Flow**
### **Creating a Pipeline**
1. User submits a pipeline creation request.
2. Backend validates input and stores metadata in Supabase.
3. User receives a confirmation with pipeline details.

### **Executing a Pipeline**
1. User selects a pipeline and starts execution.
2. Backend retrieves the pipeline configuration.
3. **ParallelPipelineOrchestrator** manages execution.
4. Each stage updates its status (**Pending → Running → Completed**).
5. WebSockets push real-time updates to the frontend.

## **Deployment & Scaling**
- **Kubernetes-Based Deployment**
  - Backend & Frontend deployed as separate microservices.
  - **Autoscaling** backend instances based on load.
  - **Dockerized** containers for isolation.
- **Messaging System (RabbitMQ/NATS)** ensures event-driven execution.

## **How to Run the Project**
### **Method 1: Running Locally**
#### **Step 1: Install Required Tools**
Ensure you have the following installed:
- [Golang v1.21+](https://go.dev/doc/install)
- [Node.js & npm](https://nodejs.org/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop) (with Kubernetes enabled)

#### **Step 2: Clone the Repository**
```sh
git clone https://github.com/sarika-p9/my-pipeline-project.git
cd my-pipeline-project
```

#### **Step 3: Create a `.env` File**
Create a `.env` file in the project root with the following:
```sh
SUPABASE_URL=https://<Project_ID>.supabase.co
SUPABASE_KEY=xxxxxxxxxxxxxxxxxxxx
POSTGRES_DSN=postgresql://postgres.<Project_ID>:<password>@aws-0-ap-south-1.pooler.supabase.com:6543/postgres
```

#### **Step 4: Load Environment Variables**
Uncomment these lines in the following files to enable local execution:
```go
// if err := godotenv.Load(); err != nil {
//     log.Println("Warning: No .env file found. Proceeding with existing environment variables.")
// }
```
**Files to update:**
- `cmd/democtl/cmd/start.go`
- `cmd/grpc_server/main.go`
- `cmd/main_server/main.go`
- `internal/infrastructure/database.go`
- `internal/infrastructure/supabase_client.go`

#### **Step 5: Start Messaging Services**
```sh
nats-server -DV
kubectl apply -f rabbitmq-deployment.yaml
```

#### **Step 6: Start Backend & Frontend**
```sh
# Start Backend
go run cmd/main_server/main.go

# Start Frontend
cd frontend
npm install  # First-time setup only
npm start
```

#### **Step 7: Access the Application**
Open **[http://localhost:3000](http://localhost:3000)** in your browser.

---

### **Method 2: Running with Docker Images**
#### **Step 1: Install Docker Desktop**
Ensure [Docker Desktop](https://www.docker.com/products/docker-desktop) is installed and running with Kubernetes enabled.

#### **Step 2: Pull & Run Docker Images**
```sh
docker run sarikapt9/my-backend
docker run sarikapt9/my-frontend
```

#### **Step 3: Create Kubernetes Secrets**
```sh
kubectl create secret generic supabase-postgres-secret \
  --from-literal=SUPABASE_URL="https://<Project_ID>.supabase.co" \
  --from-literal=SUPABASE_KEY="xxxxxxxxxxxxxxxxxxxx" \
  --from-literal=POSTGRES_DSN="postgresql://postgres.<Project_ID>:<password>@aws-0-ap-south-1.pooler.supabase.com:6543/postgres"
```

#### **Step 4: Deploy Application to Kubernetes**
```sh
kubectl apply -f backend-deployment.yaml
kubectl apply -f frontend-deployment.yaml
kubectl apply -f rabbitmq-deployment.yaml
kubectl apply -f nats-deployment.yaml
```

#### **Step 5: Access the Application**
Once deployed, access the frontend at **[http://localhost:30001](http://localhost:30001)**.

## **Using democtl CLI**
```sh
# Register a new user
./democtl register --email="xyz@abc.com" --password="xxxxx"

# Login
./democtl login --email="xyz@abc.com" --password="xxxxx"

# Create a pipeline
./democtl pipeline create --user="xxxxx" --stages=3 --pipeline-name="TestPipeline" --stage-names="a,b,c"

# Start pipeline execution
./democtl pipeline start --pipeline-id="xxxxx" --user-id="xxxxx" --input="{}"

# Get pipeline status
./democtl pipeline status --pipeline-id="xxxxx"
```

## **Conclusion**
The **Distributed Manufacturing Pipeline Simulation System** provides a scalable, real-time solution for managing **manufacturing pipelines**. It leverages modern **cloud-native** and **microservices** principles, ensuring high performance and efficiency.
