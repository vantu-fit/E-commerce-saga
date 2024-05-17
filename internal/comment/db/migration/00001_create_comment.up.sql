CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "comments" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "product_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "content" varchar(255) NOT NULL,
  "left_index" BIGINT NOT NULL,
  "right_index" BIGINT NOT NULL,
  "parent_id" uuid,
  "upadated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

insert into comments (
	id,
	product_id,
	user_id,
	content,
	left_index,
	right_index,
	parent_id
) values 
	('9b3c3532-68ca-4859-bb4f-9e7fda5f886c','747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','ROOT',1,20,NULL),
	('244dec51-570d-4545-9d02-c70ab1aa05e8','747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','TELEVISIONS',2,9,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','TUBES',3,4,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','LCD',5,6,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','PLASMA',7,8,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	('4c2fcee2-f5a7-4f82-8334-5177c3412206','747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','PORTABLE ELECTRONICS',10,19,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	('6032c38a-e677-4df0-b80b-cb086eb5dbe8','747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','MP3 PLAYER',11,14,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','CD PLAYER',15,16,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','2 WAY RADIOS',17,18,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c'),
	( uuid_generate_v4(),'747adae5-3849-4093-9a44-f5a1bc9a8ef9','0903ad25-555d-404d-bf8b-7e3f12b820ad','PLASH',12,13,'9b3c3532-68ca-4859-bb4f-9e7fda5f886c')
;
