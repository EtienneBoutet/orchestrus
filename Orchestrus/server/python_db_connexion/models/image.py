import json
from sqlalchemy.dialects.postgresql import JSON
from db import db

class Image(db.Model):
  __tablename__ = "images"
  id = db.Column(db.Text, primary_key=True, autoincrement=False)
  host = db.Column(db.Text, primary_key=True, autoincrement=False)
  name = db.Column(db.Text, nullable=False)
  port = db.Column(JSON)  
  worker_ip = db.Column(db.Text, db.ForeignKey('worker.ip'), nullable=False)

  # def __init__(self, host, id, name, port):
  #   self.host = host
  #   self.id = id
  #   self.name = name
  #   self.port = port

  # def json(self):
  #   return {'host': self.host, 'id': self.id, 'name': self.name, 'port': self.port}