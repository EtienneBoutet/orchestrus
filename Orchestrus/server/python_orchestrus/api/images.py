import flask_rebar
from flask_rebar import errors
from rebar import registry
from schemas.image import ImageSchema, CreateImageSchema
from models.image import Image
import requests
from app import Config
import helpers

@registry.handles(
  rule="/workers/<string:host>/images/",
  method="POST",
  request_body_schema=CreateImageSchema(),
  response_body_schema={
    200: ImageSchema()
  }
)
def add_image(host):
  body = flask_rebar.get_validated_body()

  name = body["name"]
  if host is None:
    raise errors.UnprocessableEntity("You need to supply a name.")

  port = body["port"]
  if host is None:
    raise errors.UnprocessableEntity('You need to supply a port. i.e : {"8080" : "80"}')

  if not helpers.is_host_reachable("http://" + host + ":5000/"):
    raise errors.UnprocessableEntity("The host is unreachable.")

  image = Image(host, None, name, port)

  # Send image information to host worker
  response = requests.post("http://" + host + ":5000/images", json=image.json())
  if response.status_code != '201':
    raise errors.UnprocessableEntity("The image could not be started on host.")

  # We get the image ID
  image.id = response.text

  # Add to database
  response = requests.post("http://" + Config.DB_MODULE_URL + "/workers/" + host + "/images", json=image.json())
  if response.status_code != '200':
    raise errors.UnprocessableEntity("The image could not be added to the database.")
    # TODO - Stop image on worker if database fail

  return image

@registry.handles(
  rule="/workers/<string:host>/images/<string:id>",
  method="POST",
)
def delete_image(host, id):
  if not helpers.is_host_reachable("http://" + host + ":5000/"):
    raise errors.UnprocessableEntity("The host is unreachable.")

  # Remove image from host worker
  response = requests.delete("http://" + host + "/images/" + id)
  if response.status_code != '204':
    raise errors.UnprocessableEntity("The image could not be stopped on the host.")

  response = requests.delete("http://" + Config.DB_MODULE_URL + "/workers/" + host + "/images/" + id)
  if response.status_code != '200':
    raise errors.UnprocessableEntity("The image could not be removed from database.")