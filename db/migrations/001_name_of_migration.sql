CREATE TABLE IF NOT EXISTS public.canvas (
	id serial4 NOT NULL,
	canvas_name text NOT NULL,
	url text NOT NULL,
	update_time date DEFAULT CURRENT_DATE NOT NULL,
	CONSTRAINT canvas_pkey PRIMARY KEY (id)
);