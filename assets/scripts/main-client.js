const fc = (function () {
    const scrollToBottom = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = elem.scrollHeight;
    }

    const scrollToTop = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = 0;
    }

    const elementExist = (element) => {
        return element != null && typeof element != "undefined";
    }

    const keyExist = (key, collection) => {
        return key in collection;
    }

    const isFunc = (v) => {
        return typeof v == "function"
    }

    const isObject = (v) => {
        return (v && typeof v === 'object' && !Array.isArray(v));
    }

    const deepMerge = (dst, ...src) => {
        if (!src.length) return dst;
        const source = src.shift();
        if (isObject(dst) && isObject(source)) {
            for (const key in source) {
                if (source.hasOwnProperty(key) && isObject(source[key])) {
                    if (!dst[key]) Object.assign(dst, {[key]: {}});
                    deepMerge(dst[key], source[key]);
                } else {
                    Object.assign(dst, {[key]: source[key]});
                }
            }
        }
        return deepMerge(dst, ...src);
    }

    const addClass = (e, c) => {
        if (!elementExist(e)) {
            return;
        }
        e.classList.add(c);
    }

    const removeClass = (e, c) => {
        if (!elementExist(e)) {
            return;
        }
        e.classList.remove(c);
    }

    const ajax = {};
    const call = (key, method, endpoint, data, success, fail, aborted, progress) => {
        //abort previous request when exist
        if (fc.keyExist(key, ajax) && [1, 2, 3].includes(ajax[key].status)) {
            ajax[key].abort();
        }
        ajax[key] = new XMLHttpRequest();
        ajax[key].withCredentials = true;
        ajax[key].addEventListener('progress', function (e) {
            fc.isFunc(progress) && progress(e);
        });
        ajax[key].addEventListener('load', function (e) {
            fc.isFunc(success) && success(e);
        });
        ajax[key].addEventListener('error', function (e) {
            fc.isFunc(fail) && fail(e);
        });
        ajax[key].addEventListener('abort', function (e) {
            fc.isFunc(aborted) && aborted(e);
        });
        ajax[key].open(method, endpoint);
        ajax[key].setRequestHeader('Content-Type', 'application/json');
        ajax[key].send(data);
    }

    const dme = document.querySelector('.modal');
    const dialogTitle = dme.querySelector('.modal-card-title');
    const dialogContent = dme.querySelector('.content');
    const dialogButtonOk = dme.querySelector('.button-ok');
    const dialogButtonCancel = dme.querySelector('.button-cancel');
    const dialogClose = dme.querySelector('.delete');
    const events = {};

    const overrideEvents = (key, func) => {
        if (keyExist(key, events)) {
            delete events[key];
        }
        events[key] = func;
    };

    const triggerEvent = (key) => {
        if (keyExist(key, events) && isFunc(events[key])) {
            events[key]();
        }
    }

    // buttons format: [ ['ok', 'Yes'], ['cancel', 'No'] ]
    const dialog = (visible, title, content, buttons) => {
        if (!elementExist(dme)) {
            return;
        }
        if (title) {
            dialogTitle.innerText = title;
        }
        if (content) {
            dialogContent.innerHTML = content;
        }
        dialogButtonOk.style.display = 'none';
        dialogButtonCancel.style.display = 'none';
        if (buttons) {
            buttons.forEach(i => {
                if (i.length !== 2) {
                    return;
                }
                switch (i[0]) {
                    case 'ok':
                        dialogButtonOk.innerText = i[1];
                        dialogButtonOk.style.display = 'block';
                        break;
                    case 'cancel':
                        dialogButtonCancel.innerText = i[1];
                        dialogButtonCancel.style.display = 'block';
                }
            })
        }
        if (visible) {
            addClass(dme, 'is-active');
        } else {
            removeClass(dme, 'is-active');
        }
    }
    // dialog global listener
    if (elementExist(dialogClose)) {
        dialogClose.addEventListener('click', function (e) {
            e.stopPropagation();
            dialog(false);
        });
    }
    if (elementExist(dialogButtonCancel)) {
        dialogButtonCancel.addEventListener('click', function (e) {
            e.stopPropagation();
            dialog(false);
        });
    }
    if (elementExist(dialogButtonOk)) {
        dialogButtonOk.addEventListener('click', function (e) {
            e.stopPropagation();
            triggerEvent('onClickOk');
            removeClass(dme, 'is-active');
        });
    }

    // value should be in serializeArray
    const toJson = function (value) {
        if (typeof value !== 'object') {
            return {};
        }
        return JSON.stringify(value);
    }

    return {
        scrollToTop,
        scrollToBottom,
        elementExist,
        keyExist,
        isFunc,
        isObject,
        deepMerge,
        call,
        dialog,
        addClass,
        removeClass,
        toJson,
        overrideEvents
    };
})();
