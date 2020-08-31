from flask_restful import Resource, reqparse
from models.worker import Worker
import requests

class Worker(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('ip', type=str, required=True, help="The IP is required.")

  def post(self):
    data = Worker.parser.parse_args()
    worker = Worker(data['ip'], False, [])

    # Check if worker is active
    response = requests.get("http://" + worker.ip + ":5000/")
    if response.status_code == '204':
      worker.active = True

    return worker.json(), 201

class DeleteWorker(Resource):
  def delete(self, host, id):
    # Verify that the host is reachable
    response = requests.get("http://" + host + ":5000/")
    if response.status_code != '204':
      return "The host is unreachable.", 404

    # Send image deletion to host worker
    response = requests.delete("http://" + host + "/images/" + id)
    if response.status_code != '204':
      return "Error from host worker", 500

    # TODO - Remove image from database 