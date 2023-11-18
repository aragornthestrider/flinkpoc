docker pull bitnami/flink:1.17.1 

FLINK_PROPERTIES="jobmanager.rpc.address: jobmanager"

docker run -d --name=jobmanager -p 6123:6123 -p 8081:8081 \
    -e FLINK_MODE=jobmanager \
    -e FLINK_CFG_REST_BIND__ADDRESS=0.0.0.0 \
    bitnami/flink:1.17.1

docker run -d --name=taskmanager -p 6121:6121 -p 6122:6122 \
    -e FLINK_MODE=taskmanager \
    -e FLINK_JOB_MANAGER_RPC_ADDRESS=jobmanager \
    bitnami/flink:1.17.1




services:
  jobmanager:
    image: docker.io/bitnami/flink:1
    ports:
      - 6123:6123
      - 8081:8081
    environment:
      - FLINK_MODE=jobmanager
      - FLINK_CFG_REST_BIND__ADDRESS=0.0.0.0
  taskmanager:
    image: docker.io/bitnami/flink:1
    ports:
      - 6121:6121
      - 6122:6122
    environment:
      - FLINK_MODE=taskmanager
      - FLINK_JOB_MANAGER_RPC_ADDRESS=jobmanager