from flask_rebar import RequestSchema
from marshmallow import fields, Schema

class ImageSchema(Schema):
  """Schema for getting an image's information"""
  id = fields.Integer(required=True)
  name = fields.String(required=True)
  port = fields.Raw(required=True)
  worker_ip = fields.Integer(required=True)

class CreateImageSchema(Schema):
  """Schema for adding an image"""
  id = fields.Integer(required=True)
  name = fields.String(required=True)
  port = fields.Raw(required=True)