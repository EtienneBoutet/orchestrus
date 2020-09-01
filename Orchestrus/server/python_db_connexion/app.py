from flask import Flask
from db import db
from rebar import REBAR
from flask_migrate import Migrate
from models.image import Image
from models.worker import Worker
from api import workers

# Config and app starter related 
app = Flask(__name__)

# TODO - Better config and clean up this file
app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql+psycopg2://postgres:postgres@127.0.0.1:5432/orchestrus"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

REBAR.init_app(app)

# TODO - Fix pour utiliser les schemas pour les migrations et non les modèles
migrate = Migrate(app, db)
db.init_app(app)

# Endpoints
# api.add_resource(AddImage, '/workers/<string:host>/images')
# api.add_resource(RemoveImage, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
  app.run(debug=True)
