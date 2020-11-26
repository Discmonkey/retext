create table qode.insight
(
    id         SERIAL PRIMARY KEY,
    value      text,
    project_id int,
    CONSTRAINT fk_insight_project_id
        FOREIGN KEY (project_id)
            REFERENCES project (id)
);

create table qode.insight_text
(
    rowId      SERIAL PRIMARY KEY,
    insight_id int,
    text_id    int,
    CONSTRAINT fk_insight_id
        FOREIGN KEY (insight_id)
            REFERENCES insight (id),
    CONSTRAINT fk_insight_text
        FOREIGN KEY (text_id)
            REFERENCES text (id)
)