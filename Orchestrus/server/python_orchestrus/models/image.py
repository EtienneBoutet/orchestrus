class ImageModel:
  def __init__(self, host, id, name, port):
    self.host = host
    self.id = id
    self.name = name
    self.port = port

  def json(self):
    return {'host': self.host, 'id': self.id, 'name': self.name, 'port': self.port}