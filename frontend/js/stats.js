// js/stats.js

// Глобальная переменная для режима просмотра (true - Реальные, false - Усредненные)
window.isRealExpenses = true; 

// Переключатель режима
window.toggleExpenseMode = function() {
    window.isRealExpenses = !window.isRealExpenses;
    const btn = document.getElementById('mode-toggle-btn');
    if (window.isRealExpenses) {
        btn.textContent = ' Показать усреднённые (Бюджет)';
        btn.style.background = '#1f6feb'; // Синий
    } else {
        btn.textContent = 'Показать реальные траты';
        btn.style.background = '#238636'; // Зеленый
    }
    window.renderStats();
};

window.renderStats = function() {
    const subs = window.subscriptions || [];
    
    // Названия месяцев
    const monthNames = ['Янв', 'Фев', 'Мар', 'Апр', 'Май', 'Июн', 'Июл', 'Авг', 'Сен', 'Окт', 'Ноя', 'Дек'];
    
    // Инициализируем массив значений для 12 месяцев нулями
    const monthValues = Array(12).fill(0);
    const now = new Date();
    const currentYear = now.getFullYear();

    // --- РАСЧЁТ РАСХОДОВ ---
    subs.forEach(sub => {
        const cost = sub.cost || 0;
        const payDate = new Date(sub.next_payment_date);
        const monthIndex = payDate.getMonth(); // 0 - Янв, 11 - Дек
        
        // Если дата платежа в другом году, пропускаем (для простоты считаем текущий год)
        if (payDate.getFullYear() !== currentYear) return;

        if (sub.billing_cycle === 'monthly') {
            // Ежемес. добавляем ко всем месяцам
            for (let i = 0; i < 12; i++) {
                monthValues[i] += cost;
            }
        } else {
            // Годовые и Разовые
            if (window.isRealExpenses) {
                // РЕЖИМ "РЕАЛЬНЫЕ": Добавляем полную сумму ТОЛЬКО в месяц оплаты
                monthValues[monthIndex] += cost;
            } else {
                // РЕЖИМ "УСРЕДНЁННЫЙ": Делим на 12 и добавляем ко всем месяцам
                const avgCost = cost / 12;
                for (let i = 0; i < 12; i++) {
                    monthValues[i] += avgCost;
                }
            }
        }
    });

    // --- ОТРИСОВКА ГРАФИКА ---
    const chart = document.getElementById('bar-chart');
    chart.innerHTML = '';
    
    // Находим максимум для масштабирования графика
    const maxVal = Math.max(...monthValues, 10); // Минимум 10, чтобы график не схлопнулся

    monthValues.forEach((val, index) => {
        const heightPercent = (val / maxVal) * 150; // Высота в пикселях (макс 150px)
        
        // Подсветка текущего месяца
        const isCurrent = index === now.getMonth() ? 'border-bottom: 2px solid #58a6ff;' : '';
        const barColor = index === now.getMonth() ? 'linear-gradient(to top, #238636, #3fb950)' : 'linear-gradient(to top, #1f6feb, #58a6ff)';

        chart.innerHTML += `
            <div class="bar" style="height:${Math.max(heightPercent, 4)}px; background: ${barColor}; ${isCurrent}">
                <span class="bar-value">${val > 0 ? '₽' + Math.round(val) : ''}</span>
                <span class="bar-label">${monthNames[index]}</span>
            </div>`;
    });

    // --- ТОП ПОДПИСОК (Всегда реальные цены, без деления) ---
    const top = [...subs].sort((a, b) => (b.cost || 0) - (a.cost || 0)).slice(0, 5);
    document.getElementById('top-expensive').innerHTML = top.map(s => 
        `<li><span>${escapeHtml(s.name)}</span><span><strong>₽${s.cost?.toFixed(2)}</strong></span></li>`
    ).join('') || '<li style="color:#8b949e">Нет данных</li>';

    // --- РАСПРЕДЕЛЕНИЕ ПО ПЕРИОДАМ ---
    const periods = { monthly: 0, yearly: 0, 'one-time': 0 };
    subs.forEach(s => { if (periods[s.billing_cycle] !== undefined) periods[s.billing_cycle]++; });
    
    const periodLabels = { 
        monthly: '📅 Ежемесячно', 
        yearly: '📆 Ежегодно', 
        'one-time': '📊 Разово' 
    };
    
    document.getElementById('period-dist').innerHTML = Object.entries(periods).map(([key, val]) =>
        `<li><span>${periodLabels[key]}</span><span>${val} подписок</span></li>`
    ).join('');
};