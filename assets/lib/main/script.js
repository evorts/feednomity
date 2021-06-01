const fc = (function () {
    const getCookieDomain = () => {
        const hName = location.hostname;
        if (hName.indexOf(".") < 0) {
            return hName;
        }
        const hArr = hName.split('.');
        if (hArr.length < 2) {
            return hName;
        }
        return `${hArr[hArr.length - 2]}.${hArr[hArr.length - 1]}`;
    }

    const setCookie = (key, value, hours) => {
        let expires = '';
        const secure = location.protocol.indexOf('https') > -1 ? 'Secure;' : '';
        value = value || '';
        if (hours) {
            const date = new Date();
            date.setTime(date.getTime() + (hours * 60 * 60 * 1000));
            expires = `; expires=${date.toUTCString()}`;
        }
        console.log(getCookieDomain());
        document.cookie = `${key}=${value}${expires}; SameSite=Strict; path=/;${secure} domain=${getCookieDomain()}`;
    }

    const getCookie = (key) => {
        const eq = `${key}=`;
        const ca = document.cookie.split(';');
        for (let i = 0; i < ca.length; i++) {
            let c = ca[i];
            while (c.charAt(0) === ' ') {
                c = c.substring(1, c.length);
            }
            if (c.indexOf(eq) === 0) {
                return c.substring(eq.length, c.length);
            }
        }
        return null;
    }

    const removeCookie = (key) => {
        document.cookie = `${key}=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;`;
    }
    const sessionKey = "feednomisess";
    const getRequestHeaders = () => {
        return [
            ['X-Authorization', getCookie(sessionKey)]
        ];
    }

    const onDocumentReady = (fn) => {
        if (document.readyState === "complete" || document.readyState === "interactive") {
            setTimeout(fn, 1);
        } else {
            document.addEventListener("DOMContentLoaded", fn);
        }
    }

    const scrollToBottom = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = elem.scrollHeight;
    }

    const scrollToTop = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = 0;
    }

    const toFullUrl = (v, base) => {
        v = v || '';
        if (/http?s:.*/.test(v)) {
            return v;
        }
        if (fc.isEmpty(base)) {
            base = `${location.protocol}//${location.hostname}`;
            if (!fc.isEmpty(location.port)) {
                base = `${base}:${location.port}`
            }
        }
        return `${base}/${v}`;
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

    const isEmpty = (v) => {
        return typeof v === "undefined" || v === null || v === 0 || v === "";
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
        if (Array.isArray(c)) {
            c.forEach(function (ce) {
                e.classList.add(ce);
            })
            return
        }
        e.classList.add(c);
    }

    const removeClass = (e, c) => {
        if (!elementExist(e)) {
            return;
        }
        if (Array.isArray(c)) {
            c.forEach(function (ce) {
                e.classList.remove(ce);
            })
            return
        }
        e.classList.remove(c);
    }

    /** ajax **/
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
            try {
                const res = JSON.parse(ajax[key].response);
                fc.isFunc(success) && success(res);
            } catch (e) {
                console.log(e);
            }
        });
        ajax[key].addEventListener('error', function (e) {
            fc.isFunc(fail) && fail(e);
        });
        ajax[key].addEventListener('abort', function (e) {
            fc.isFunc(aborted) && aborted(e);
        });
        ajax[key].open(method, endpoint);
        ajax[key].setRequestHeader('Content-Type', 'application/json');
        const headers = getRequestHeaders();
        if (Array.isArray(headers) && headers.length > 0) {
            headers.forEach(function (item) {
                ajax[key].setRequestHeader(item[0], item[1]);
            })
        }
        ajax[key].send(data);
    }

    /** form **/
    const getFormData = (form, inputSelector) => {
        let data = {};
        const selector = inputSelector || '.form-input';
        const formElements = form.elements;
        const inputElements = form.querySelectorAll(selector);
        for (let i = 0; i < inputElements.length; i++) {
            const name = inputElements[i].name;
            let value = formElements[name].value;
            const type = inputElements[i].type;
            if (type === 'radio') {
                if (value === "") {
                    value = 0;
                } else {
                    value = parseInt(value);
                }
            }
            if (name.endsWith(']')) {
                const n = name.split('[')[0];
                if (fc.keyExist(n, data)) {
                    let v = data[n];
                    if (!Array.isArray(v)) {
                        continue;
                    }
                    v.push(value);
                    data[n] = v;
                } else {
                    data[n] = [value];
                }
                continue;
            }
            if (keyExist(name, data)) {
                continue;
            }
            if (name.indexOf('.') < 1) {
                data[name] = value;
                continue;
            }
            //parse and rebuild into nested object when naming convention contains dot
            const item = name.split('.').reduceRight(
                (all, item) => ({[item]: all}), value
            );
            data = deepMerge(data, item);
        }
        return data;
    }

    /** toast **/
    const toastTypes = ['is-warning', 'is-danger', 'is-success', 'is-info'];
    let t = document.getElementById('snackbar');
    if (!elementExist(t)) {
        t = document.createElement('div');
        t.id = "snackbar";
        addClass(t, 'notification');
        document.body.appendChild(t);
    }
    const toast = (text, type) => {
        const c = ['show'];
        t.innerText = text;
        removeClass(t, toastTypes);
        if (toastTypes.includes(type)) {
            c.push(type);
        }
        addClass(t, c);
        setTimeout(function () {
            removeClass(t, 'show');
        }, 1500);
    }

    /** warning field **/
    const warn = function (response) {
        const $elements = document.querySelectorAll('*[class^="err-"]');
        if ($elements.length < 1) {
            return;
        }
        for (let i = 0; i < $elements.length; i++) {
            $elements[i].innerText = '';
            fc.removeClass($elements[i],'is-hidden');
            fc.addClass($elements[i], 'is-hidden');
        }
        if (response['error']['message'].length > 0) {
            fc.toast(response['error']['message'], 'is-danger');
        }
        if (fc.isObject(response['error']['reasons'])) {
            for (let [key, value] of Object.entries(response['error']['reasons'])) {
                const $element = document.querySelector('.err-' + key);
                if (fc.elementExist($element)) {
                    $element.innerText = value;
                    fc.removeClass($element, 'is-hidden');
                }
            }
        }
    };

    /** dialog **/
    const dialogTemplate = `
        <div class="modal-background"></div>
        <div class="modal-card" style="min-height: 12rem;">
            <header class="modal-card-head">
                <p class="modal-card-title">Modal title</p>
                <button class="delete" aria-label="close"></button>
            </header>
            <section class="modal-card-body">
                <div class="content" style="min-height: 3rem;"></div>
            </section>
            <footer class="modal-card-foot">
                <button class="button button-ok is-success">Save changes</button>
                <button class="button button-cancel">Cancel</button>
            </footer>
        </div>`
    let dme = document.getElementById('dialog-modal');
    if (!elementExist(dme)) {
        dme = document.createElement('div');
        dme.id = 'dialog-modal';
        addClass(dme, 'modal');
        dme.innerHTML = dialogTemplate;
        document.body.appendChild(dme);
    }
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
            triggerEvent('onClickCancel');
            dialog(false);
        });
    }
    if (elementExist(dialogButtonCancel)) {
        dialogButtonCancel.addEventListener('click', function (e) {
            e.stopPropagation();
            triggerEvent('onClickCancel');
            dialog(false);
        });
    }
    if (elementExist(dialogButtonOk)) {
        dialogButtonOk.addEventListener('click', function (e) {
            e.stopPropagation();
            triggerEvent('onClickOk');
            dialog(false);
        });
    }

    return {
        onDocumentReady,
        getFormData,
        scrollToTop,
        scrollToBottom,
        elementExist,
        keyExist,
        isFunc,
        isObject,
        isEmpty,
        deepMerge,
        call,
        dialog,
        warn,
        toast,
        addClass,
        removeClass,
        overrideEvents,
        toFullUrl,
        setCookie,
        getCookie,
        removeCookie,
        sessionKey,
    };
})();
