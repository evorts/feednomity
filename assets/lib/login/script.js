(function (fc, apiUrl, RedirectUrl) {
    fc.onDocumentReady(function () {
        const form = document.getElementById('loginForm');
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
                        fc.setCookie(fc.sessionKey, res.content['token'], 24);
                        fc.toast('Login success! Redirecting...', 'is-success');
                        setTimeout(function () {
                            window.location.replace(RedirectUrl);
                        }, 700);
                        return;
                    }
                    fc.warn(res);
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
})(fc, ApiBaseUrl, RedirectUrl)