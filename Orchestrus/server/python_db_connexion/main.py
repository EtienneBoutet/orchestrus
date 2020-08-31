from flask_sqlalchemy import SQLAlchemy
from flask import Flask
from flask_restful import Api, Resource
from flask_migrate import Migrate
# from resources.worker import Worker, DeleteWorker
# from resources.image import Image

app = Flask(__name__)
app.config['SQLALCHEMY_DATABASE_URI'] = "postgresql+psycopg2://postgres:postgres@127.0.0.1:5432/orchestrus"
app.config['SQLALCHEMY_TRACK_MODIFICATIONS'] = False
db = SQLAlchemy(app)
migrate = Migrate(app, db)

api = Api(app)

# from app import db ======= afin d'acceder à la db dans les autres modules

# api.add_resource(Worker, '/workers')
# api.add_resource(Image, '/images')
# api.add_resource(DeleteWorker, '/workers/<string:host>/images/<string:id>')

if __name__ == '__main__':
	app.run(debug=True)
