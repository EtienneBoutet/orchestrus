from flask_sqlalchemy import SQLAlchemy
from functools import partial

db = SQLAlchemy()
# make columns non-nullable by default, most of them should be
db.Column = partial(db.Column, nullable=False)
