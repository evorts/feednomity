(function (fc, apiUrl, RedirectUrl) {
    fc.onDocumentReady(function () {
        const form = document.getElementById('fpForm');
        const btnLogin = document.getElementById('submitForgot');
        btnLogin.addEventListener('click', function (e) {
            const $this = this;
            $this.setAttribute('disabled', "disabled");
            e.preventDefault();

            const data = fc.getFormData(form);
            fc.call(
                "fp", "POST", `${apiUrl}/v1/users/forgot-password`,
                JSON.stringify(data),
                function (res) {
                    $this.removeAttribute('disabled');
                    if (res.status === 200) {
                        fc.toast('Request forgot password sent! Check your email. Redirecting to login page...', 'is-success');
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