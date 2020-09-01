import flask_rebar
from flask_rebar import errors
from schemas.worker import WorkerSchema
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
  return Worker.query.all()