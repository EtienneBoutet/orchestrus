from flask import Flask
from rebar import rebar
from config import Config
from api import workers, images

app = Flask(__name__)

app.config.from_object(Config)

rebar.init_app(app)

if __name__ == '__main__':
  app.run(debug=True, port=5001)