CREATE TABLE public.users (
    id serial NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE TABLE public.posts (
    id serial NOT NULL,
    user_id integer NOT NULL,
    post character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.posts ADD CONSTRAINT posts_pkey PRIMARY KEY (id);
ALTER TABLE public.posts
    ADD CONSTRAINT posts_relation_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;

CREATE TABLE public.post_likes (
    id serial NOT NULL,
    post_id integer NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE public.post_likes ADD CONSTRAINT post_likes_pkey PRIMARY KEY (id);
CREATE UNIQUE INDEX post_likes_unique ON public.post_likes USING btree (post_id, user_id);
ALTER TABLE public.post_likes
    ADD CONSTRAINT post_likes_post_id FOREIGN KEY (post_id) REFERENCES public.posts(id) ON UPDATE CASCADE ON DELETE CASCADE;
ALTER TABLE public.post_likes
    ADD CONSTRAINT post_likes_user_id FOREIGN KEY (user_id) REFERENCES public.users(id) ON UPDATE CASCADE ON DELETE CASCADE;
