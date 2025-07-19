# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')

### K8s Config ###

k8s_yaml('./deploy/development/k8s/secrets.yaml')

k8s_yaml('./deploy/development/k8s/app-config.yaml')

# PostgreSQL Database
k8s_yaml('./deploy/development/k8s/postgres-deployment.yaml')
k8s_resource('postgres', labels="databases")

# RabbitMQ Message Broker
k8s_yaml('./deploy/development/k8s/rabbitmq-deployment.yaml')
k8s_resource('rabbitmq', port_forwards=15672, labels="messaging")

### End of K8s Config ###
### API Gateway ###

# Generate Swagger docs before compiling
local_resource(
  'swagger-docs-generate',
  'make swagger-docs',
  deps=['./services/api-gateway/*.go', './services/api-gateway/grpc_clients/*.go'], 
  labels="compiles")

gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './deploy/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway/*.go', './services/api-gateway/grpc_clients/*.go', './services/api-gateway/docs/*.go', './pkg'], 
  resource_deps=['swagger-docs-generate'],
  labels="compiles")


docker_build_with_restart(
  'bolt-app/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./deploy/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './pkg',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./pkg', '/app/pkg'),
  ],
)

k8s_yaml('./deploy/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile'], labels="services")

             
### End of API Gateway ###
### Trip Service ###

trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/main.go'
if os.name == 'nt':
 trip_compile_cmd = './deploy/development/docker/trip-build.bat'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './pkg'], labels="compiles")

docker_build_with_restart(
  'bolt-app/trip-service',
  '.',
  entrypoint=['/app/build/trip-service'],
  dockerfile='./deploy/development/docker/trip-service.Dockerfile',
  only=[
    './build/trip-service',
    './pkg',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./pkg', '/app/pkg'),
  ],
)

k8s_yaml('./deploy/development/k8s/trip-service-deployment.yaml')
k8s_resource('trip-service', resource_deps=['trip-service-compile', 'postgres', 'rabbitmq'], labels="services")

### End of Trip Service ###

### Driver Service ###

driver_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/driver-service ./services/driver-service'
if os.name == 'nt':
 driver_compile_cmd = './deploy/development/docker/driver-build.bat'

local_resource(
  'driver-service-compile',
  driver_compile_cmd,
  deps=['./services/driver-service', './pkg'], labels="compiles")

docker_build_with_restart(
  'bolt-app/driver-service',
  '.',
  entrypoint=['/app/build/driver-service'],
  dockerfile='./deploy/development/docker/driver-service.Dockerfile',
  only=[
    './build/driver-service',
    './pkg',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./pkg', '/app/pkg'),
  ],
)

k8s_yaml('./deploy/development/k8s/driver-service-deployment.yaml')
k8s_resource('driver-service', resource_deps=['driver-service-compile', 'rabbitmq'], labels="services")

### End of Driver Service ###

### Payment Service ###

payment_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/payment-service ./services/payment-service/cmd/main.go'
if os.name == 'nt':
  payment_compile_cmd = './deploy/development/docker/payment-build.bat'

local_resource(
  'payment-service-compile',
  payment_compile_cmd,
  deps=['./services/payment-service', './pkg'], labels="compiles")

docker_build_with_restart(
  'bolt-app/payment-service',
  '.',
  entrypoint=['/app/build/payment-service'],
  dockerfile='./deploy/development/docker/payment-service.Dockerfile',
  only=[
    './build/payment-service',
    './pkg',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./pkg', '/app/pkg'),
  ],
)

k8s_yaml('./deploy/development/k8s/payment-service-deployment.yaml')
k8s_resource('payment-service', resource_deps=['payment-service-compile', 'rabbitmq'], labels="services")

### End of Payment Service ###

### Jaeger ###
k8s_yaml('./deploy/development/k8s/jaeger.yaml')
k8s_resource('jaeger', port_forwards=['16686:16686', '14268:14268'], labels="tooling")
### End of Jaeger ###

### Swagger UI ###
k8s_yaml('./deploy/development/k8s/swagger-ui-deployment.yaml')
k8s_resource('swagger-ui', port_forwards=8082, resource_deps=['api-gateway'], labels="tooling")
### End of Swagger UI ###