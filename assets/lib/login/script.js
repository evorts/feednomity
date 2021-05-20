(function (fc, apiUrl) {
    const getRedirectUrl = (def) => {
        const urlParams = new URLSearchParams(window.location.search);
        const ref = urlParams.get('ref');
        if (!fc.isEmpty(ref)) {
            return fc.toFullUrl(ref);
        }
        if (typeof window['redirectUrl'] !== "undefined" && !fc.isEmpty(window['redirectUrl'])) {
            return fc.toFullUrl(window['redirectUrl']);
        }
        return fc.toFullUrl(def);
    }
    fc.onDocumentReady(function () {
        const form = document.getElementById('loginForm');
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
        const btnLogin = document.getElementById('submitLogin');
        btnLogin.addEventListener('click', function (e) {
            const $this = this;
            $this.setAttribute('disabled', "disabled");
            e.preventDefault();
            const data = fc.getFormData(form);
            fc.call(
                "login", "POST", `${apiUrl}/v1/users/login`,
                JSON.stringify(data),
                function (res) {
                    $this.removeAttribute('disabled');
                    if (res.status === 200) {
                        fc.setCookie(fc.sessionKey, res.content['token'], 24)
                        window.location.replace(getRedirectUrl());
                        return;
                    }
                    warn(res);
                },
                function () {
                    $this.removeAttribute('disabled');
                },
                function () {
                    $this.removeAttribute('disabled');
                },
            )
        })
    });
})(fc, ApiBaseUrl)