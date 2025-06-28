from fastapi import FastAPI
from features.buildings import routes as buildings_routes
app = FastAPI()


app.include_router(buildings_routes.router, prefix="/buildings", tags=["buildings"])