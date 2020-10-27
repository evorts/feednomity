(function ($, $f) {
    $(document).ready(function () {
        const warn = function(response) {
            $('*[class^="err-"]').text('').removeClass('hide').addClass('hide');
            const $globalElement = $('.err-global');
            if (response['error']['message'].length > 0) {
                $globalElement.text(response['error']['message']).removeClass('hide');
            }
            if (response['error']['reasons'].length > 0) {
                $.each(response['error']['reasons'], function (key, value) {
                    const $element = $('.err-'+key);
                    if ($element.length) {
                        $element.text(value).removeClass('hide')
                    }
                });
            }
        };
        $('#submitLogin').on('click', function (e) {
            const $this = $(this);
            $this.prop('disabled', 1);
            e.preventDefault();
            $.ajax({
                url: '/api/login',
                method: 'POST',
                dataType: 'json',
                contentType: 'application/json',
                data: JSON.stringify($('form[name=loginForm]').find('input').serializeObject())
            }).done(function () {
                window.location.replace('/dashboard');
            }).fail(function (xhr) {
                warn(JSON.parse(xhr.responseText));
            }).always(function () {
                $this.prop('disabled', 0);
            });
        });
    })
})(jQuery, feednomity)