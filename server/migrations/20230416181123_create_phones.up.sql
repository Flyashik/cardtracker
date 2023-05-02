CREATE TABLE IF NOT EXISTS phones
(
    sim_coutn integer NOT NULL,
    producer "char"[] NOT NULL,
    model "char"[] NOT NULL,
    sd boolean NOT NULL,
    mainter "char"[],
    CONSTRAINT "phones_pkey" PRIMARY KEY (producer, model)
);