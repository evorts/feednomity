(function (fc, apiUrl, RedirectUrl) {
    fc.onDocumentReady(function () {
        const form = document.getElementById('crpForm');
        const btnLogin = document.getElementById('submitPassword');
        btnLogin.addEventListener('click', function (e) {
            const $this = this;
            $this.setAttribute('disabled', "disabled");
            e.preventDefault();
            const data = fc.getFormData(form);
            fc.call(
                "crp", "POST", `${apiUrl}/v1/users/create-password`,
                JSON.stringify(data),
                function (res) {
                    $this.removeAttribute('disabled');
                    if (res.status === 200) {
                        fc.toast('Create password success!', 'is-success');
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