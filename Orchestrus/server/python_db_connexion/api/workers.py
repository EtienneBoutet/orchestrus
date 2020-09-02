import flask_rebar
from flask_rebar import errors
from schemas.worker import WorkerSchema, CreateWorkerSchema
from db import db
from models.worker import Worker
from rebar import registry

@registry.handles(
  rule="/workers",
  method="GET",
  response_body_schema={
    200: WorkerSchema(many=True)
  }
)
def workers_list():
  """Return all workers"""
  return Worker.query.all()

@registry.handles(
  rule="/workers",
  method="POST",
  request_body_schema=CreateWorkerSchema(),
  response_body_schema={
    200: WorkerSchema()
  }
)
def create_worker():
  body = flask_rebar.get_validated_body()

  ip = body["ip"]

  if not ip:
    raise errors.UnprocessableEntity("Please add an IP.") 

  worker = Worker.query.filter_by(ip=ip).first()
  
  if worker is not None:
    raise errors.UnprocessableEntity("This worker already exists.")

  worker = Worker(ip, body["active"])

  db.session.add(worker)
  db.session.commit()

  return worker

# TODO - Remove a worker