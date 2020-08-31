from flask import Flask
from flask_restful import Api, Resource
from resources.worker import AddWorker, DeleteWorker
from resources.image import AddImage

app = Flask(__name__)
api = Api(app)

api.add_resource(AddWorker, '/workers')
api.add_resource(AddImage, '/images')
api.add_resource(DeleteWorker, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)