from fastapi import APIRouter

router = APIRouter()
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
async def construct_building():
    return "Construct building"