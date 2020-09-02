import flask_rebar
from flask_rebar import errors
from rebar import registry
from schemas.worker import CreateWorkerSchema, WorkerSchema
from models.worker import Worker
import requests
from app import Config
import helpers

@registry.handles(
  rule="/workers",
  method="POST",
  request_body_schema=CreateWorkerSchema(),
  response_body_schema={
    200: WorkerSchema()
  }
)
def add_worker():
  body = flask_rebar.get_validated_body()

  ip = body["ip"]

  if ip is None:
    raise errors.UnprocessableEntity("You need to supply an IP address.")
  
  isActive = helpers.is_host_reachable("http://" + ip + ":5000/")

  worker = Worker(ip, isActive)

  # Add to database
  response = requests.delete("http://" + Config.DB_MODULE_URL + "/workers", json=worker.json())
  if response.status_code != '200':
    raise errors.UnprocessableEntity("The image could not be added to the database.")

  return worker

@registry.handles(
  rule="/workers",
  method="GET",
  response_body_schema={
    200: WorkerSchema(many=True)
  }
)
def workers_list():
  response = requests.get("http://" + Config.DB_MODULE_URL + "/workers")
  return response.text