CREATE TABLE IF NOT EXISTS public."user" (
    id serial4 NOT NULL,
    email text NOT NULL,
    "password" text NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (id)
);

INSERT INTO public."user" (email, password)
VALUES ('somemeial', 'hehehahha');

CREATE TABLE IF NOT EXISTS public.folder (
    id serial4 NOT NULL,
    "name" text NOT NULL,
    parent_folder_id int4,
    user_id int4 NOT NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT folder_pkey PRIMARY KEY (id),
    CONSTRAINT folder_parent_folder_id_fkey FOREIGN KEY (parent_folder_id) 
        REFERENCES public.folder(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT folder_user_id_fkey FOREIGN KEY (user_id) 
        REFERENCES public."user"(id) ON DELETE CASCADE ON UPDATE CASCADE
);

ALTER TABLE public.folder ADD CONSTRAINT root_folder_check 
CHECK (
    (id = 0 AND parent_folder_id IS NULL) OR 
    (id != 0 AND parent_folder_id IS NOT NULL)
);

INSERT INTO public.folder (id, "name", parent_folder_id, user_id)
VALUES (0, 'THE_ROOT_OF_ALL', NULL, 1);

CREATE TABLE IF NOT EXISTS public.canvas (
    id serial4 NOT NULL,
    "name" text NOT NULL,
    url text NULL,
    updated_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "text" text NULL,
    folder_id int4 NOT NULL,
    user_id int4 NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT canvas_pkey PRIMARY KEY (id),
    CONSTRAINT canvas_folder_id_fkey FOREIGN KEY (folder_id) 
        REFERENCES public.folder(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT canvas_user_id_fkey FOREIGN KEY (user_id) 
        REFERENCES public."user"(id) ON DELETE CASCADE ON UPDATE CASCADE
);