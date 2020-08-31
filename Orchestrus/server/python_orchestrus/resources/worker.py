from flask_restful import Resource, reqparse
from models.worker import WorkerModel
import requests

class AddWorker(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('ip', type=str, required=True, help="The IP is required.")

  def post(self):
    data = AddWorker.parser.parse_args()
    worker = WorkerModel(data['ip'], False, [])

    # Verify that the host is reachable
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