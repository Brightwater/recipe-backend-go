CREATE TABLE public.recipes (
	id int4 DEFAULT nextval('recipes_column1_seq'::regclass) NOT NULL,
	author varchar NOT NULL,
	"data" jsonb NOT NULL,
	ts timestamp NULL,
	CONSTRAINT recipes_pk PRIMARY KEY (id)
);