from flask_restful import Resource, reqparse
from models.worker import WorkerModel
import requests

class Worker(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('ip', type=str, required=True, help="The IP is required.")

  def post(self):
    data = Worker.parser.parse_args()
    worker = WorkerModel(data['ip'], False, [])

    # Check if worker is active
    response = requests.get("http://" + worker.ip + ":5000/")
    if response.status_code == '204':
      worker.active = True

    return worker.json(), 201