
import os

class Config(object):
  POSTGRES_URL = os.environ.get("POSTGRES_URL", "127.0.0.1:5432")
  POSTGRES_USER = os.environ.get("POSTGRES_USER", "postgres")
  POSTGRES_PW = os.environ.get("POSTGRES_PW", "postgres")
  POSTGRES_DB = os.environ.get("POSTGRES_DB", "orchestrus")
  DB_URL = 'postgresql+psycopg2://{user}:{password}@{url}/{db}'.format(user=POSTGRES_USER,password=POSTGRES_PW,url=POSTGRES_URL,db=POSTGRES_DB)
  SQLALCHEMY_DATABASE_URI = DB_URL
  SQLALCHEMY_TRACK_MODIFICATIONS = False