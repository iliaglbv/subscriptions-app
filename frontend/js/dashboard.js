window.renderDashboard = function() {
    const subs = window.subscriptions || [];
    
    // Всего подписок
    document.getElementById('total-count').textContent = subs.length;

    // Расчёт месячных расходов
    let monthly = 0;
    subs.forEach(s => {
        if (s.billing_cycle === 'monthly') monthly += s.cost || 0;
        else if (s.billing_cycle === 'yearly') monthly += (s.cost || 0) / 12;
        else monthly += s.cost || 0; // one-time считаем как единовременный
    });

    document.getElementById('total-monthly').textContent = `₽${monthly.toFixed(2)}`;
    document.getElementById('total-yearly').textContent = `₽${(monthly * 12).toFixed(2)}`;

    // Ближайший платёж
    const sorted = [...subs].sort((a, b) => 
        new Date(a.next_payment_date) - new Date(b.next_payment_date)
    );
    
    const nextPayEl = document.getElementById('next-payment');
    if (sorted.length > 0) {
        const days = window.daysUntil(sorted[0].next_payment_date);
        nextPayEl.textContent = days <= 0 ? 'Сегодня' : `Через ${days} дн.`;
    } else {
        nextPayEl.textContent = '—';
    }

    // Список ближайших платежей
    const upcoming = document.getElementById('upcoming-list');
    upcoming.innerHTML = '';
    sorted.slice(0, 5).forEach(s => {
        const days = window.daysUntil(s.next_payment_date);
        const label = days <= 0 ? 'Сегодня' : days === 1 ? 'Завтра' : `Через ${days} дн.`;
        upcoming.innerHTML += `
        <li>
            <span>${escapeHtml(s.name)} — <strong>₽${s.cost?.toFixed(2)}</strong></span>
            <span class="upcoming-date">${label}</span>
        </li>`;
    });

    // Расходы по категориям
    const cats = {};
    subs.forEach(s => {
        let m = s.billing_cycle === 'monthly' ? (s.cost || 0) : 
                s.billing_cycle === 'yearly' ? (s.cost || 0) / 12 : (s.cost || 0);
        const cat = s.category || 'Другое';
        cats[cat] = (cats[cat] || 0) + m;
    });

    const colors = {
        'Развлечения': '#f85149', 'Музыка': '#3fb950', 'Софт': '#58a6ff',
        'Облако': '#d29922', 'Образование': '#bc8cff', 'Другое': '#8b949e'
    };

    const catList = document.getElementById('cat-list');
    catList.innerHTML = '';
    Object.entries(cats)
        .sort((a, b) => b[1] - a[1])
        .forEach(([name, val]) => {
            catList.innerHTML += `
            <li>
                <span class="cat-name">
                    <span class="cat-dot" style="background:${colors[name] || '#8b949e'}"></span>
                    ${escapeHtml(name)}
                </span>
                <span>₽${val.toFixed(2)}/мес</span>
            </li>`;
        });
};