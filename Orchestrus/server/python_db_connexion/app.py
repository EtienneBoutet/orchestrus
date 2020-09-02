from flask import Flask
from db import db
from rebar import REBAR
from flask_migrate import Migrate
from models.image import Image
from models.worker import Worker
from api import workers, images
from config import Config

# Config and app starter related 
app = Flask(__name__)

app.config.from_object(Config)

migrate = Migrate(app, db)

migrate.init_app(app)
REBAR.init_app(app)
db.init_app(app)

if __name__ == '__main__':
  app.run(debug=True)