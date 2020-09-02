from flask import Flask
from db import db
from rebar import REBAR
from flask_migrate import Migrate
from models.image import Image
from models.worker import Worker
from api import workers, images

# Config and app starter related 
app = Flask(__name__)

# TODO - Better config and clean up this file
app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql+psycopg2://postgres:postgres@127.0.0.1:5432/orchestrus"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False

migrate = Migrate(app, db)

migrate.init_app(app)
REBAR.init_app(app)
db.init_app(app)

if __name__ == '__main__':
  app.run(debug=True)