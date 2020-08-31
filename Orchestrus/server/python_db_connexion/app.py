from flask import Flask
from flask_restful import Api, Resource
from db import db
from flask_migrate import Migrate
# from resources.worker import Worker, DeleteWorker
# from resources.image import Image
from models.image import Image
from models.worker import Worker

app = Flask(__name__)

app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql+psycopg2://postgres:postgres@127.0.0.1:5432/orchestrus"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

migrate = Migrate(app, db)

api = Api(app)

# api.add_resource(Worker, '/workers')
# api.add_resource(Image, '/images')
# api.add_resource(DeleteWorker, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)
