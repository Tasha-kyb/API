-- Создание таблицы lists
CREATE TABLE IF NOT EXISTS lists (
    id UUID PRIMARY KEY,
    title VARCHAR(100) NOT NULL CHECK (length(title) >= 1 AND length(title) <= 100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Индекс для сортировки по дате создания
CREATE INDEX idx_lists_created_at ON lists(created_at DESC);

-- Комментарии для документации
COMMENT ON TABLE lists IS 'Списки задач пользователей';
COMMENT ON COLUMN lists.id IS 'Уникальный идентификатор списка (UUID)';
COMMENT ON COLUMN lists.title IS 'Название списка (1-100 символов)';
COMMENT ON COLUMN lists.created_at IS 'Дата и время создания списка';