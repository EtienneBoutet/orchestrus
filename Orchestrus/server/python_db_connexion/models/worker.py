from db import db

class Worker(db.Model):
  ip = db.Column(db.Text, primary_key=True, autoincrement=False)
  active = db.Column(db.Boolean, nullable=False)
  images = db.relationship('Image', backref='worker', lazy=True)

  def __init__(self, ip, active, images=[]):
    self.ip = ip
    self.active = active
    self.images = images

  def json(self):
    return {'ip': self.ip, 'active': self.active, 'images': self.images}