// Форматирование даты для отображения
window.formatDate = function(dateStr) {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('ru-RU', {
        day: 'numeric',
        month: 'short',
        year: 'numeric'
    });
};

// Получение сегодняшней даты без времени
window.getToday = function() {
    const d = new Date();
    d.setHours(0, 0, 0, 0);
    return d;
};

// Расчёт дней до даты
window.daysUntil = function(dateStr) {
    const target = new Date(dateStr);
    const diff = target - window.getToday();
    return Math.ceil(diff / 86400000);
};