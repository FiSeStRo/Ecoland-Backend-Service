from fastapi import FastAPI, Depends

from fastapi.security import OAuth2PasswordBearer

from database.db import Base, engine
from features.authentication.service import get_current_user

from features.buildings import routes as buildings_routes
from features.users import routes as users_routes
from features.authentication import routes as auth_routes

import features.users.models
import features.buildings.models

app = FastAPI()
Base.metadata.create_all(bind=engine)


# TODO: Add map and lan / lat to entries
app.include_router(buildings_routes.router, prefix="/buildings", tags=["buildings"], dependencies=[Depends(get_current_user)])
#no auth for testing
# app.include_router(buildings_routes.router, prefix="/buildings", tags=["buildings"])
app.include_router(users_routes.router)
app.include_router(auth_routes.router)