
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session

from database.db import SessionLocal
from features.buildings.models import DefBuildings, Buildings
from features.buildings.shemas import ConstructBuildingRequest

router = APIRouter()

def get_db():
    db =SessionLocal()
    try:
        yield db
    finally:
        db.close()

db_dependency = Annotated[Session, Depends(get_db)]

@router.get("/list")
async def get_all_buildings():
    return "Get all buildings"

@router.get("/buildings/details/{building_id}")
async def get_building(building_id: int):
    return "Get building details"

@router.get("/constructionlist")
async def get_all_possible_constructions():
    return "Get all possible constructions"

@router.get("/productions/{building_id}")
async def get_productions(building_id: int):
    return "Get productions"

@router.post("/construct")
async def construct_building(body: ConstructBuildingRequest, db: db_dependency):

    if db.query(DefBuildings).filter(DefBuildings.id == body.def_id).first() is None:
        raise HTTPException(status_code=404, detail=f"Could not find building with def_id of {body.def_id}")

    building = Buildings(
        user_id = 1,
        def_id = body.def_id,
        name = body.display_nane
    )

    db.add(building)
    db.commit()
    return