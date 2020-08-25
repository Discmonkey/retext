-- main schema
CREATE SCHEMA IF NOT EXISTS qode;

-- so we don't need to reference everything with qode
SET search_path TO qode,public;

-- time_tag is the time entered by the user on project creation
CREATE TABLE IF NOT EXISTS project(
    id SERIAL PRIMARY KEY,
    name text,
    description text,
    created time,
    time_tag time
);

CREATE TABLE IF NOT EXISTS file(
    id SERIAL UNIQUE PRIMARY KEY,
    project_id int,
    name text,
    uploaded timestamp with time zone,

    CONSTRAINT fk_project
        FOREIGN KEY (project_id)
            REFERENCES project(id)
);

-- source file is used for primary documents such as interviews
CREATE TABLE IF NOT EXISTS source_file() inherits (file);

-- demo files are usually tabular data related to linking multiple source files to a set
-- of attributes.
CREATE TABLE IF NOT EXISTS demo_file() inherits (file);

-- store the columns with associated column name, ie age, gender, conditions
CREATE TABLE IF NOT EXISTS attribute_column(
    id SERIAL PRIMARY KEY,
    name text,
    demo_file_id int,
    CONSTRAINT fk_demo_file
        FOREIGN KEY (demo_file_id)
            REFERENCES file(id)
);

-- store the specific row where an attribute was found, useful for linking to source files
CREATE TABLE IF NOT EXISTS attribute_row(
    id SERIAL PRIMARY KEY,
    demo_file_id int,
    
    CONSTRAINT fk_demo_file
        FOREIGN KEY (demo_file_id)
            REFERENCES file(id),
            
    source_file_id int default NULL,
    
    CONSTRAINT fk_source_file
        FOREIGN KEY (source_file_id)
            REFERENCES file(id)
);

-- stores the actual value of the attributes
CREATE TABLE IF NOT EXISTS attribute(
    id SERIAL PRIMARY KEY,
    value text,
    attribute_column_id int,
    CONSTRAINT fk_attribute_column
        FOREIGN KEY (attribute_column_id)
            REFERENCES attribute_column(id),
    attribute_row_id int,
    CONSTRAINT fk_attribute_row
        FOREIGN KEY (attribute_row_id)
            REFERENCES attribute_row(id)
);

-- stores the top level code containers, their names default to the code added to the container
CREATE TABLE IF NOT EXISTS code_container(
    id SERIAL PRIMARY KEY,
    display_order int,
    project_id int,

    CONSTRAINT fk_project
        FOREIGN KEY (project_id)
            REFERENCES project(id)
);

-- about code within the container, text is associated to this item
CREATE TABLE IF NOT EXISTS code(
    id SERIAL PRIMARY KEY,
    display_order int,
    name text,
    code_container_id int,
    CONSTRAINT fk_code_container
        FOREIGN KEY(code_container_id)
            REFERENCES code_container(id)
);

-- create a type for referencing our place within a specific document
DO $$
BEGIN
IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'anchor') THEN
CREATE TYPE anchor AS
(
    paragraph int,
    sentence int,
    word int
);
END IF;
END $$;

-- stores the current commit hash for the parser being used, allows any anchor to be reproduced
CREATE TABLE IF NOT EXISTS parser(
    id SERIAL PRIMARY KEY,
    commit char(40)
);

CREATE TABLE IF NOT EXISTS text(
    id SERIAL PRIMARY KEY,
    start anchor,
    stop anchor,
    value text, 
    parser_id int,

    CONSTRAINT fk_parser
        FOREIGN KEY(parser_id)
            REFERENCES parser(id),

    source_file_id int,
    CONSTRAINT fk_source_file
        FOREIGN KEY(source_file_id)
            REFERENCES file(id),

    code_id int,
    CONSTRAINT fk_code
        FOREIGN KEY(code_id)
            REFERENCES code(id)
);






