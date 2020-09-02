from flask_rebar import RequestSchema
from marshmallow import fields, Schema
from schemas.image import ImageSchema

class WorkerSchema(Schema):
  """Schema for getting a worker's information"""
  ip = fields.String(required=True)
  active = fields.Boolean(required=True)
  images = fields.Nested(ImageSchema, many=True, required=True)

class CreateWorkerSchema(Schema):
  """Schema for creating a worker"""
  ip = fields.String(required=True)
  active = fields.Boolean(required=True)
  images = fields.List(fields.Nested(ImageSchema))
