-- Создание таблицы tasks
CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    list_id UUID NOT NULL,
    text VARCHAR(500) NOT NULL,
    completed boolean DEFAULT False,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Добавляем внешний ключ с каскадным удалением
ALTER TABLE tasks 
ADD CONSTRAINT fk_tasks_list_id 
FOREIGN KEY (list_id) 
REFERENCES lists(id) 
ON DELETE CASCADE;

-- Создаем индексы для производительности
CREATE INDEX idx_tasks_id ON tasks(id);
CREATE INDEX idx_tasks_list_id ON tasks(list_id);
CREATE INDEX idx_tasks_completed ON tasks(completed);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);

-- Комментарии для документации
COMMENT ON TABLE tasks IS 'Задачи пользователей';
COMMENT ON COLUMN tasks.id IS 'Уникальный идентификатор задачи (UUID)';
COMMENT ON COLUMN tasks.text IS 'Описание задачи (1-500 символов)';
COMMENT ON COLUMN tasks.completed IS 'Статус выполнения задачи';
COMMENT ON COLUMN tasks.created_at IS 'Дата и время создания задачи';
COMMENT ON COLUMN tasks.updated_at IS 'Дата и время обновления задачи';