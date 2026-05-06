// js/api.js
async function apiFetch(endpoint, options = {}) {
    const token = localStorage.getItem('jwt_token');
    
    // Гарантируем заголовок Content-Type
    const headers = {
        'Content-Type': 'application/json',
        ...(options.headers || {})
    };
    
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    // Если есть body и это объект — строкуем в JSON
    let body = options.body;
    if (body && typeof body === 'object' && !(body instanceof FormData)) {
        body = JSON.stringify(body);
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
        body
    });

    // Читаем ответ
    const text = await response.text();
    let data;
    try {
        data = text ? JSON.parse(text) : {};
    } catch {
        data = { error: 'Неверный формат ответа' };
    }

    if (!response.ok) {
        throw new Error(data.error || data.details || `Ошибка ${response.status}`);
    }

    return data;
}