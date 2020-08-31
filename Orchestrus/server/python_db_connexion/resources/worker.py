from flask_restful import Resource, reqparse
from models.worker import Worker
from db import db

class GetWorker(Resource):
  def get(self):
    try:
      return [worker.json() for worker in Worker.query.all()], 200
    except Exception:
      return "Could not fetch workers from the DB.", 400

class PostWorker(Resource):
  parser = reqparse.RequestParser()
  parser.add_argument('ip', type=str, required=True, help="The ip is required.")
  parser.add_argument('active', type=bool, required=True, help="The active status is required.")

  def post(self):
    data = PostWorker.parser.parse_args()
    worker = Worker(data['ip'], data['active'])

    try:
      db.session.add(worker)
      db.session.commit()

      return worker.json(), 201
    except Exception as e:
      return "Could not add this worker to the DB.", 400