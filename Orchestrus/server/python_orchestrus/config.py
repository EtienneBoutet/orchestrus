import os

class Config(object):
  DB_MODULE_URL = os.environ.get("DB_MODULE_URL", "127.0.0.1:5000/api/")