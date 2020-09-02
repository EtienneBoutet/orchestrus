class Image:
  def __init__(self, worker_ip, id, name, port):
    self.worker_ip = worker_ip
    self.id = id
    self.name = name
    self.port = port

  def json(self):
    return {'worker_ip': self.worker_ip, 'id': self.id, 'name': self.name, 'port': self.port}