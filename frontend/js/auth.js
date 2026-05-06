// Переключение между формами входа/регистрации
window.toggleAuth = function(type) {
    document.getElementById('login-form').classList.toggle('active', type === 'login');
    document.getElementById('register-form').classList.toggle('active', type === 'register');
    document.getElementById('auth-error').textContent = '';
    document.getElementById('auth-error').style.color = '#f85149';
};

// Обработка входа
window.handleLogin = async function() {
    const username = document.getElementById('login-username').value.trim();
    const password = document.getElementById('login-password').value;
    const btn = document.getElementById('login-btn');
    
    if (!username || !password) {
        document.getElementById('auth-error').textContent = 'Заполните все поля';
        return;
    }

    btn.disabled = true;
    btn.textContent = 'Вход...';

    try {
        const data = await apiFetch('/login', {
            method: 'POST',
            body: JSON.stringify({ username, password })
        });
        
        localStorage.setItem('jwt_token', data.access_token);
        window.initApp();
    } catch (e) {
        document.getElementById('auth-error').textContent = e.message;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Войти';
    }
};

// Обработка регистрации
window.handleRegister = async function() {
    const username = document.getElementById('reg-username').value.trim();
    const email = document.getElementById('reg-email').value.trim();
    const password = document.getElementById('reg-password').value;
    const btn = document.getElementById('register-btn');
    
    if (!username || !email || !password) {
        document.getElementById('auth-error').textContent = 'Заполните все поля';
        return;
    }

    if (password.length < 6) {
        document.getElementById('auth-error').textContent = 'Пароль должен быть не менее 6 символов';
        return;
    }

    btn.disabled = true;
    btn.textContent = 'Регистрация...';

    try {
        await apiFetch('/register', {
            method: 'POST',
            body: JSON.stringify({ username, email, password })
        });
        
        document.getElementById('auth-error').style.color = '#3fb950';
        document.getElementById('auth-error').textContent = '✅ Регистрация успешна! Войдите.';
        window.toggleAuth('login');
    } catch (e) {
        document.getElementById('auth-error').style.color = '#f85149';
        document.getElementById('auth-error').textContent = e.message;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Создать аккаунт';
    }
};

// Выход из аккаунта
window.handleLogout = function() {
    localStorage.removeItem('jwt_token');
    location.reload();
};