from flask import Flask
from flask_restful import Api, Resource
from resources.worker import Worker


app = Flask(__name__)
api = Api(app)

api.add_resource(Worker, '/workers')

if __name__ == '__main__':
  app.run(debug=True)