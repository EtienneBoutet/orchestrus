from flask_rebar import Rebar

REBAR = Rebar()
registry = REBAR.create_handler_registry(prefix="/api")
