from flask import Flask
from functools import partial
from flask_restful import Api
from db import db
from flask_migrate import Migrate
from resources.worker import GetWorker, PostWorker
from flask_rebar import Rebar
from resources.image import AddImage, RemoveImage
from models.image import Image
from models.worker import Worker

# Config and app starter related 
app = Flask(__name__)

# TODO - Better config and clean up this file
app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql+psycopg2://postgres:postgres@127.0.0.1:5432/orchestrus"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

REBAR = Rebar()
REGISTRY = REBAR.create_handler_registry(prefix="/api")
REBAR.init_app(app)


migrate = Migrate(app, db)
db.init_app(app)
# make columns non-nullable by default, most of them should be
db.Column = partial(db.Column, nullable=False)

api = Api(app)

# Endpoints
api.add_resource(GetWorker, '/workers')
api.add_resource(PostWorker, '/workers')

api.add_resource(AddImage, '/workers/<string:host>/images')
api.add_resource(RemoveImage, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)
