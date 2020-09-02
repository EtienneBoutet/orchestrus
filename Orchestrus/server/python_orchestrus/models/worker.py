class Worker:
  def __init__(self, ip, active, images=[]):
    self.ip = ip
    self.active = active
    self.images = images

  def json(self):
    return {'ip': self.ip, 'active': self.active, 'images': self.images}