CREATE EXTENSION IF NOT EXISTS pgcrypto;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_type
        WHERE typname = 'task_status'
    ) THEN
CREATE TYPE task_status AS ENUM ('todo', 'in_progress', 'done');
END IF;
END
$$;



CREATE TABLE IF NOT EXISTS boards (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_boards PRIMARY KEY (id),
    CONSTRAINT uq_boards_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS board_columns (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    position INT NOT NULL DEFAULT 0,
    board_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_board_columns PRIMARY KEY (id),
    CONSTRAINT fk_board_columns_boards FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE,
    CONSTRAINT uq_board_columns_board_id_name UNIQUE (board_id, name),
    CONSTRAINT chk_board_column_position CHECK (position >= 0)
);

CREATE TABLE IF NOT EXISTS sprints (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    board_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_sprints PRIMARY KEY (id),
    CONSTRAINT fk_sprints_boards FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS groups (
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_groups PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'todo',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    author_id UUID NOT NULL,
    assignee_id UUID,
    board_id UUID,
    board_column_id UUID,
    sprint_id UUID,
    group_id UUID,

    CONSTRAINT pk_tasks PRIMARY KEY (id),
    CONSTRAINT fk_tasks_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE SET NULL,
    CONSTRAINT fk_tasks_board_column FOREIGN KEY (board_column_id) REFERENCES board_columns(id) ON DELETE SET NULL,
    CONSTRAINT fk_tasks_sprint FOREIGN KEY (sprint_id) REFERENCES sprints(id) ON DELETE SET NULL,
    CONSTRAINT fk_tasks_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS task_watchers (
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,

    CONSTRAINT pk_task_watchers PRIMARY KEY (task_id, user_id),
    CONSTRAINT fk_task_watchers_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS task_executors (
    task_id UUID NOT NULL,
    user_id UUID NOT NULL,

    CONSTRAINT pk_task_executors PRIMARY KEY (task_id, user_id),
    CONSTRAINT fk_task_executors_task FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE
    );

CREATE INDEX idx_tasks_created_at ON tasks(created_at);
CREATE INDEX idx_tasks_author_id ON tasks(author_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_tasks_board_column_id ON tasks(board_column_id);
CREATE INDEX idx_tasks_board_id ON tasks(board_id);
CREATE INDEX idx_tasks_sprint_id ON tasks(sprint_id);
CREATE INDEX idx_tasks_group_id ON tasks(group_id);

CREATE INDEX idx_sprints_board_id ON sprints(board_id);
CREATE INDEX idx_boards_created_at ON boards(created_at);
CREATE INDEX idx_sprints_created_at ON sprints(created_at);
CREATE INDEX idx_groups_created_at ON groups(created_at);