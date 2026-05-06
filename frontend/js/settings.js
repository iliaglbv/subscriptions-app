// Экспорт данных в JSON
window.exportData = function() {
    const data = {
        exported_at: new Date().toISOString(),
        subscriptions: window.subscriptions || []
    };
    
    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = `subscriptions_${new Date().toISOString().slice(0,10)}.json`;
    a.click();
    URL.revokeObjectURL(a.href);
};

// Очистка всех подписок (через API)
window.clearAll = async function() {
    if (!confirm('Удалить ВСЕ подписки? Это действие нельзя отменить.')) return;

    try {
        // Удаляем по одной (bulk delete не реализован в бэкенде)
        for (const sub of (window.subscriptions || [])) {
            if (sub.id) {
                await apiFetch(`/subscriptions/${sub.id}`, { method: 'DELETE' });
            }
        }
        await window.refreshSubscriptions();
        window.renderAll();
        alert('✅ Все подписки удалены');
    } catch (e) {
        alert('Ошибка: ' + e.message);
    }
};