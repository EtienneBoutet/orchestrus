import flask_rebar
from flask_rebar import errors
from schemas.image import ImageSchema, CreateImageSchema
from db import db
from models.image import Image
from models.worker import Worker
from rebar import registry

@registry.handles(
  rule="/workers/<string:ip>/images",
  method="POST",
  request_body_schema=CreateImageSchema(),
  response_body_schema={
    200: ImageSchema()
  }
)
def add_image(ip):
  body = flask_rebar.get_validated_body()

  # Verify that worker exists
  worker = Worker.query.filter_by(ip=ip).first()
  if worker is not None:
    raise errors.UnprocessableEntity("This worker does not exists.")

  # Fields handling
  id = body["id"]
  name = body["name"]

  if not id:
    raise errors.UnprocessableEntity("Please add an ID to the image.")

  if not name:
    raise errors.UnprocessableEntity("Please add a name to the image.")

  image = Image.query.filter_by(worker_ip=ip, id=id).first()
  if image is not None:
    raise errors.UnprocessableEntity("This worker's image already exists.")
  
  image = Image(ip, id, name, body["port"])
  
  db.session.add(image)
  db.session.commit()

  return image