// Отрисовка таблицы подписок
window.renderTable = function() {
    const tbody = document.getElementById('subs-tbody');
    tbody.innerHTML = '';

    if (!window.subscriptions || window.subscriptions.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" style="text-align:center;color:#8b949e">Нет подписок</td></tr>';
        return;
    }

    window.subscriptions.forEach(sub => {
        const badgeClass = sub.billing_cycle === 'monthly' ? 'monthly' : 
                          sub.billing_cycle === 'yearly' ? 'yearly' : 'one-time';
        const badgeText = sub.billing_cycle === 'monthly' ? 'Ежемес.' : 
                         sub.billing_cycle === 'yearly' ? 'Ежегодн.' : 'Разово';

        tbody.innerHTML += `
        <tr>
            <td><strong>${escapeHtml(sub.name)}</strong></td>
            <td>${escapeHtml(sub.category || '—')}</td>
            <td>₽${sub.cost?.toFixed(2) || '0.00'}</td>
            <td><span class="badge ${badgeClass}">${badgeText}</span></td>
            <td>${window.formatDate(sub.next_payment_date)}</td>
            <td><button class="btn-delete" onclick="window.deleteSub(${sub.id})">✕</button></td>
        </tr>`;
    });
};

// Удаление подписки
window.deleteSub = async function(id) {
    if (!confirm('Удалить эту подписку?')) return;

    try {
        await apiFetch(`/subscriptions/${id}`, { method: 'DELETE' });
        await window.refreshSubscriptions();
        window.renderAll();
    } catch (e) {
        alert('Ошибка: ' + e.message);
    }
};

// Добавление подписки
window.addSubscription = async function(e) {
    e.preventDefault();
    
    const btn = document.getElementById('add-btn');
    btn.disabled = true;
    btn.textContent = 'Добавление...';

    const payload = {
        name: document.getElementById('f-name').value.trim(),
        category: document.getElementById('f-category').value,
        cost: parseFloat(document.getElementById('f-price').value),
        currency: 'RUB',
        billing_cycle: document.getElementById('f-period').value,
        next_payment_date: document.getElementById('f-date').value
    };

    try {
        await apiFetch('/subscriptions', {
            method: 'POST',
            body: JSON.stringify(payload)
        });
        
        await window.refreshSubscriptions();
        window.renderAll();
        e.target.reset();
        window.showPage('subscriptions', document.querySelectorAll('.sidebar nav button')[1]);
    } catch (err) {
        alert('Ошибка: ' + err.message);
    } finally {
        btn.disabled = false;
        btn.textContent = 'Добавить подписку';
    }
};

// Экранирование HTML для защиты от XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}