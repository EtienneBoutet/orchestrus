from flask import Flask
from flask_restful import Api, Resource
from resources.worker import Worker, DeleteWorker
from resources.image import Image

app = Flask(__name__)
api = Api(app)

api.add_resource(Worker, '/workers')
api.add_resource(Image, '/images')
api.add_resource(DeleteWorker, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)