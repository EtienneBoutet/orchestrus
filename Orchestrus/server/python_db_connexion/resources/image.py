from flask_restful import Resource, reqparse
from models.image import Image
from models.worker import Worker
from db import db

class AddImage(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('id', type=str, required=True, help="The host is required.")
  parser.add_argument('name', type=str, required=True, help="The image name is required.")
  parser.add_argument('port', type=dict, required=True, help="The image port is required.")

  def post(self, host):
    data = AddImage.parser.parse_args()
    image = Image(host, data['id'], data['name'], data['port'])

    try:
      db.session.add(image)
      db.session.commit()

      return image.json(), 201
    except Exception as e:
      print(e)
      return "Could not add this image to the DB", 400

class RemoveImage(Resource):
  def delete(self, host, id):
    try:
      worker = Worker.query.filter_by(ip=host).first()
      
      related_image = None
      for image in worker.images:
        if image.id == id:
          related_image = image

      db.session.delete(related_image)
      db.session.commit()

      return related_image.json(), 204
    except Exception as e:
      print(e)
      return "Could not delete image from the DB", 400