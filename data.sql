CREATE TABLE "users_role" (
 "id" SERIAL,
 "role_name" VARCHAR(255) NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id")
);

CREATE TABLE "users" (
 "id" SERIAL,
 "name" VARCHAR(255) NOT NULL,
 "birthday" TIMESTAMP NOT NULL,
 "mail" VARCHAR(255) UNIQUE NOT NULL,
 "password" VARCHAR(255) NOT NULL,
 "phone" VARCHAR(255) UNIQUE NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
 "role_id" SERIAL NOT NULL,
   CONSTRAINT "FK_users.role_id"
     FOREIGN KEY ("role_id")
       REFERENCES "users_role"("id")
);

CREATE TABLE "countries" (
 "id" SERIAL,
 "name" VARCHAR(255) NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id")
);

CREATE TABLE "rent_locations" (
 "id" SERIAL,
 "full_address" VARCHAR(255) NOT NULL,
 "id_of_head" INT NOT NULL,
 "country_id" INT NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_rent_locations.country_id"
     FOREIGN KEY ("country_id")
       REFERENCES "countries"("id")
);

CREATE TABLE "car_types" (
 "id" SERIAL,
 "name" VARCHAR(255) NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id")
);

CREATE TABLE "car_brands" (
 "id" SERIAL,
 "name" VARCHAR(255) NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id")
);

CREATE TABLE "car_models" (
 "id" SERIAL,
 "model_name" VARCHAR(255) NOT NULL,
 "brand_id" INT NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_car_models.brand_id"
     FOREIGN KEY ("brand_id")
       REFERENCES "car_brands"("id")
);

CREATE TABLE "car_purposes" (
 "id" SERIAL,
 "name" VARCHAR(255) NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id")
);

CREATE TABLE "cars" (
 "id" SERIAL,
 "type_id" INT NOT NULL,
 "model_id" INT NOT NULL,
 "purpose_id" INT NOT NULL,
 "year" INT NOT NULL,
 "plate" VARCHAR(255) NOT NULL UNIQUE,
 "rental_price" DECIMAL NOT NULL,
 "location_id" INT NOT NULL,
 "available" BOOL NOT NULL DEFAULT TRUE,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_cars.type_id"
     FOREIGN KEY ("type_id")
       REFERENCES "car_types"("id"),
   CONSTRAINT "FK_cars.model_id"
     FOREIGN KEY ("model_id")
       REFERENCES "car_models"("id"),
   CONSTRAINT "FK_cars.location_id"
     FOREIGN KEY ("location_id")
        REFERENCES "rent_locations"("id")
);

CREATE TABLE "car_assignments" (
 "car_id" INT NOT NULL,
 "purpose_id" INT NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
   CONSTRAINT "car_assignments.car_id"
     FOREIGN KEY ("car_id")
       REFERENCES "cars"("id"),
   CONSTRAINT "car_assignments.purpose_id"
     FOREIGN KEY ("purpose_id")
       REFERENCES "car_purposes"("id")
);


CREATE TABLE "rent_info" (
 "id" SERIAL,
 "rental_start_date" TIMESTAMP NOT NULL,
 "rental_end_date" TIMESTAMP NOT NULL,
 "rental_price" DECIMAL NOT NULL,
 "from_location_id" INT NOT NULL,
 "return_location_id" INT ,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_rent_info.from_location_id"
     FOREIGN KEY ("from_location_id")
       REFERENCES "rent_locations"("id"),
   CONSTRAINT "FK_rent_info.return_location_id"
     FOREIGN KEY ("return_location_id")
       REFERENCES "rent_locations"("id")
);

CREATE TABLE "user_history" (
 "id" SERIAL,
 "user_id" INT NOT NULL,
 "rent_info_id" INT NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_user_history.user_id"
     FOREIGN KEY ("user_id")
       REFERENCES "users"("id"),
   CONSTRAINT "FK_user_history.rent_info_id"
     FOREIGN KEY ("rent_info_id")
       REFERENCES "rent_info"("id")
);

CREATE TABLE "car_history" (
 "id" SERIAL,
 "car_id" INT NOT NULL,
 "rent_info_id" INT NOT NULL,
 "created_at" TIMESTAMP,
 "updated_at" TIMESTAMP,
 PRIMARY KEY ("id"),
   CONSTRAINT "FK_car_history.car_id"
     FOREIGN KEY ("car_id")
       REFERENCES "cars"("id"),
   CONSTRAINT "FK_car_history.rent_info_id"
     FOREIGN KEY ("rent_info_id")
       REFERENCES "rent_info"("id")
);


