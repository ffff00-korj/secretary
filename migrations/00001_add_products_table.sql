-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS public.products (
	"name" text NOT NULL,
	id serial4 NOT NULL,
	sum int4 DEFAULT 0 NOT NULL,
	paymentday int4 DEFAULT 1 NOT NULL,
	CONSTRAINT products_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.products;
-- +goose StatementEnd
