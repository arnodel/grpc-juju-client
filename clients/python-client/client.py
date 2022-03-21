import argparse
import grpc
from juju.client.application.v1.application_pb2_grpc import ApplicationServiceStub
from juju.client.application.v1.application_pb2 import DeployRequest, RemoveRequest, ResponseLineType


parser = argparse.ArgumentParser()
parser.add_argument('--remove', action='store_true')
args = parser.parse_args()

channel = grpc.insecure_channel('localhost:8080')
stub = ApplicationServiceStub(channel)

if args.remove:
    resp = stub.Remove(RemoveRequest(application_name="postgresql"))
else:
    resp = stub.Deploy(DeployRequest(artifact_name="postgresql"))

for line in resp:
    line_type = ResponseLineType.Name(line.type)
    print(line_type, line.content)
