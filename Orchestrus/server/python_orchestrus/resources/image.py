from flask_restful import Resource, reqparse
from models.image import ImageModel
import requests

# TODO - Better error handling 

class Image(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('host', type=str, required=True, help="The host is required.")
  parser.add_argument('name', type=str, required=True, help="The image name is required.")
  parser.add_argument('port', type=dict, required=True, help="The image port is required.")

  def post(self):
    data = Image.parser.parse_args()
    image = ImageModel(data['host'], None, data['name'], data['port'])

    # Verify that the host is reachable
    response = requests.get("http://" + image.host + ":5000/")
    if response.status_code != '204':
      return "The host is unreachable.", 404

    # Send image information to host worker
    response = requests.post("http://" + image.host + ":5000/images", json=image.json(), timeout=1.5)
    if response.status_code == '201':
      # TODO - Add the ID to the image
      return image.json(), 201
    else:
      return "Error from host worker", 500