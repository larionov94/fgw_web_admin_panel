document.addEventListener('DOMContentLoaded', function() {
    const searchForm = document.getElementById('searchForm');

    const searchFIO = document.getElementById('searchFIO');
    const clearSearchIconFIO = document.getElementById('clearSearchIconFIO');

    const searchTabNum = document.getElementById('searchTabNum');
    const clearSearchIconTabNum = document.getElementById('clearSearchIconTabNum');

    let debounceTimer;

    if (clearSearchIconFIO) {
        clearSearchIconFIO.addEventListener('click', function() {
            searchFIO.value = '';
            searchForm.submit();
        });
    }

    if (clearSearchIconTabNum) {
        clearSearchIconTabNum.addEventListener('click', function() {
            searchTabNum.value = '';
            searchForm.submit();
        });
    }

    // 1. Авто-поиск при вводе в поле ФИО.
    searchFIO.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });

    // 1. Авто-поиск при вводе в поле табельный номер
    searchTabNum.addEventListener('input', function(e) {
        clearTimeout(debounceTimer);

        // 2. Если поле очищено - сразу отправляем
        if (e.target.value === '') {
            searchForm.submit();
            return;
        }

        // 3. Ждем 800ms после последнего ввода
        debounceTimer = setTimeout(() => {
            searchForm.submit();
        }, 800);
    });
});