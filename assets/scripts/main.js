const feednomity = (function ($) {
    // value should be in serializeArray
    const toJson = function (value) {
        if (typeof value !== 'object') {
            return {};
        }
        console.log(value);
        console.log(typeof value);
        return JSON.stringify(value);
    }

    return {
        toJson: toJson,
        toast: function (msg) {
            alert(msg)
        }
    }
})(jQuery);
// jQuery Plugin section
(function ($) {
    $.fn.serializeObject = function () {
        let obj = {};
        const arr = this.serializeArray();
        $.each(arr, function () {
            if (obj[this.name]) {
                if (!obj[this.name].push) {
                    obj[this.name] = [obj[this.name]];
                }
                obj[this.name].push(this.value || '');
            } else {
                obj[this.name] = this.value || '';
            }
        });
        console.log(obj);
        return obj;
    }
})(jQuery)