from fastapi import FastAPI, Depends

from fastapi.security import OAuth2PasswordBearer

from database.db import Base, engine
from features.authentication.routes import get_current_user
from features.buildings import routes as buildings_routes
from features.users import routes as users_routes

import features.users.models
import features.buildings.models
oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/token")
app = FastAPI()
Base.metadata.create_all(bind=engine)


# TODO: Add map and lan / lat to entries
app.include_router(buildings_routes.router, prefix="/buildings", tags=["buildings"], dependencies=[Depends(get_current_user)])
app.include_router(users_routes.router)