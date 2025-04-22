CREATE TABLE public.canvas (
	id serial4 NOT NULL,
	canvas_name text NOT NULL,
	url text NULL,
	update_time timestamp DEFAULT CURRENT_DATE NOT NULL,
	"text" text NULL,
	CONSTRAINT canvas_pkey PRIMARY KEY (id)
);