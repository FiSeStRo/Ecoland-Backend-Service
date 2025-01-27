// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

Project eco_land {
  database_type: 'MariaDb'
}


Table users {
  id int [pk, increment]
  username varchar [unique, not null]
  password varchar [not null]
  email varchar [not null]
  role int
  time_created timestamp
  time_last_activity timestamp
}

Table user_resources {
  user_id int [pk, ref: - users.id]
  money decimal [not null, default: 100000.00] //default value comes from inital sign up
  prestige int  
}

Table buildings {
  id int [pk, increment] 
  user_id int [not null, ref: > users.id]
  def_id int [not null, ref: > def_buildings.id]
  name varchar
  time_build timestamp [not null]
    
}

//Definiton of Products
Table def_products {
  id int [pk]
  base_value int [not null]
  token_name sting
}

//Definiton Production
Table def_production {
  id int [pk]
  token_name varchar
  cost decimal
  base_duration int
}

//Definiton Recipe of any production
Table def_rel_production_product {
  production_id int [ref: > def_production.id]
  product_id int [ref: > def_products.id]
  is_input bool
  amount int 
} 

//Defintion of Buildings
Table def_buildings {
  id int [pk]
  token_name varchar
  base_construction_cost decimal
  base_construction_time int
  }

Table def_rel_building_production {
  building_id int [not null,ref: - def_buildings.id]
  production_id int [not null,ref: < def_production.id]
}

//Storage
Table rel_building_product{
 building_id int [ref:> buildings.id]
 product_id int [ref: > def_products.id]
 amount int
 capacity int [not null, default: 500] //deault comes from constant config
}

//Production order
Table rel_buildng_def_production {
  id int [pk]
  building_id int [ref: > buildings.id]
  production_id int [ref: - def_production.id]
  time_sart timestamp [not null]
  time_end timestamp [not null]
  cycles int [not null]
  is_completed bool [default: false] //for later for production backlog implementation
}




Ref: "def_rel_building_production"."building_id" < "def_rel_building_production"."production_id"
