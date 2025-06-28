from fastapi import FastAPI

from database.db import Base, engine
from features.buildings import routes as buildings_routes
from features.users import routes as users_routes

import features.users.models
import features.buildings.models

app = FastAPI()
Base.metadata.create_all(bind=engine)



app.include_router(buildings_routes.router, prefix="/buildings", tags=["buildings"])
app.include_router(users_routes.router)