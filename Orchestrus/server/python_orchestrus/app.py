from flask import Flask
from flask_restful import Api, Resource
from resources.worker import AddWorker
from resources.image import AddImage, DeleteImage

app = Flask(__name__)
api = Api(app)

api.add_resource(AddWorker, '/workers')
api.add_resource(AddImage, '/images')
api.add_resource(DeleteImage, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)