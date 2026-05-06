// Глобальное состояние
window.subscriptions = [];

// Навигация между страницами
window.showPage = function(pageId, btn) {
    // Скрыть все страницы
    document.querySelectorAll('.page').forEach(p => p.classList.remove('active'));
    // Убрать активный класс с кнопок
    document.querySelectorAll('.sidebar nav button').forEach(b => b.classList.remove('active'));
    
    // Показать нужную
    document.getElementById(pageId).classList.add('active');
    if (btn) btn.classList.add('active');
};

// Обновление данных с сервера
window.refreshSubscriptions = async function() {
    window.subscriptions = await apiFetch('/subscriptions');
};

// Полный рендер всех компонентов
window.renderAll = function() {
    window.renderTable();
    window.renderDashboard();
    window.renderStats();
};

// Инициализация приложения
window.initApp = async function() {
    const token = localStorage.getItem('jwt_token');
    
    if (!token) {
        // Не авторизован: показать форму входа
        document.getElementById('auth-container').style.display = 'flex';
        document.getElementById('app-container').style.display = 'none';
        return;
    }

    // Авторизован: показать приложение
    document.getElementById('auth-container').style.display = 'none';
    document.getElementById('app-container').style.display = 'flex';

    try {
        await window.refreshSubscriptions();
        window.renderAll();
    } catch (e) {
        console.error('Ошибка загрузки данных:', e);
        // Если токен протух — разлогинить
        if (e.message.includes('401') || e.message.includes('token')) {
            localStorage.removeItem('jwt_token');
            location.reload();
        }
    }
};

// Обработчик формы добавления
document.getElementById('add-form')?.addEventListener('submit', window.addSubscription);

// Запуск при загрузке страницы
document.addEventListener('DOMContentLoaded', window.initApp);