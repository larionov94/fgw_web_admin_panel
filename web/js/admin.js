// Дебаунс функция для ограничения частоты вызовов
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// 1. Функция проверки сессии с кешированием
let lastSessionCheck = 0;
const SESSION_CHECK_CACHE_TIME = 30000; // 30 секунд

function checkSessionBeforeNavigation(callback) {
    const now = Date.now();

    // Используем кеширование, если недавно проверяли
    if (now - lastSessionCheck < SESSION_CHECK_CACHE_TIME) {
        if (typeof callback === 'function') {
            callback();
        }
        return Promise.resolve(true);
    }

    return fetch('/api/session-check', {
        method: 'HEAD',
        credentials: 'include',
        cache: 'no-store'
    })
        .then(function (response) {
            lastSessionCheck = Date.now();

            if (!response.ok) {
                // Сессия невалидна - logout
                alert('Ваша сессия истекла. Пожалуйста, войдите снова.');
                window.location.href = '/logout';
                return false;
            }

            // Сессия валидна - выполняем callback
            if (typeof callback === 'function') {
                callback();
            }
            return true;
        })
        .catch(function (error) {
            console.error('Ошибка проверки сессии:', error);
            // При ошибке сети разрешаем навигацию (опционально)
            if (typeof callback === 'function') {
                callback();
            }
            return false;
        });
}

// 2. Оптимизированная защита ссылок
function protectNavigationLinks() {
    // Единый обработчик для всех ссылок (делегирование событий)
    document.addEventListener('click', function (e) {
        const link = e.target.closest('a[href^="/admin"]');
        if (!link) return;

        // Не проверяем для текущей страницы
        if (link.href === window.location.href) {
            return;
        }

        e.preventDefault();
        e.stopPropagation();

        checkSessionBeforeNavigation(function () {
            window.location.href = link.href;
        });
    });
}

// 3. Оптимизированный мониторинг сессии и активности
(function () {
    let sessionCheckInterval;
    let lastActivity = Date.now();
    let inactivityWarningShown = false;

    // Дебаунсированное обновление активности
    const updateActivity = debounce(() => {
        lastActivity = Date.now();
        inactivityWarningShown = false;
    }, 1000); // Обновляем не чаще чем раз в секунду

    // Проверка сессии с защитой от повторных вызовов
    function checkSession() {
        return fetch('/api/session-check', {
            method: 'HEAD',
            credentials: 'include',
            cache: 'no-store'
        })
            .then(function (response) {
                if (!response.ok) {
                    console.log('Сессия истекла, выполняется выход...');
                    window.location.href = '/logout';
                    return false;
                }
                return true;
            })
            .catch(function (error) {
                console.error('Ошибка проверки сессии:', error);
                return false; // Не выходим при ошибке сети
            });
    }

    // Проверка неактивности
    function checkInactivity() {
        const now = Date.now();
        const inactiveTime = now - lastActivity;

        // Предупреждение через 10 минут
        if (!inactivityWarningShown && inactiveTime > 10 * 60 * 1000) {
            inactivityWarningShown = true;
            const userConfirmed = confirm('Вы неактивны более 10 минут. Нажмите OK для продления сессии.');
            if (userConfirmed) {
                updateActivity();
                checkSession(); // Дополнительная проверка сессии
            }
        }

        // Выход через 15 минут
        if (inactiveTime > 15 * 60 * 1000) {
            alert('Сессия завершена из-за неактивности.');
            window.location.href = '/logout';
        }
    }

    // Инициализация
    function initSessionMonitoring() {
        // События активности с дебаунсом
        ['click', 'keypress', 'mousemove', 'scroll', 'touchstart'].forEach(function (eventName) {
            document.addEventListener(eventName, updateActivity, { passive: true });
        });

        // Проверка сессии каждые 5 минут
        sessionCheckInterval = setInterval(checkSession, 300000);

        // Проверка неактивности каждую минуту
        setInterval(checkInactivity, 60000);

        // Первая проверка
        setTimeout(checkSession, 60000);
    }

    // Запускаем при загрузке
    document.addEventListener('DOMContentLoaded', function () {
        protectNavigationLinks();
        initSessionMonitoring();

        // Очистка при закрытии
        window.addEventListener('beforeunload', function () {
            if (sessionCheckInterval) {
                clearInterval(sessionCheckInterval);
            }
        });
    });
})();