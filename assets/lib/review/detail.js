(function (fc, ratings) {
    const attachTabsFunctions = () => {
        let tabs = document.querySelectorAll('.tabs li');
        let tabsContent = document.querySelectorAll('.tab-content');

        let deactivateAllTabs = function () {
            tabs.forEach(function (tab) {
                tab.classList.remove('is-active');
            });
        };

        let hideTabsContent = function () {
            tabsContent.forEach(function (tabContent) {
                tabContent.classList.remove('is-active');
            });
        };

        let activateTabsContent = function (tab) {
            tabsContent[getIndex(tab)].classList.add('is-active');
        };

        let getIndex = function (el) {
            return [...el.parentElement.children].indexOf(el);
        };

        tabs.forEach(function (tab) {
            tab.addEventListener('click', function () {
                deactivateAllTabs();
                hideTabsContent();
                tab.classList.add('is-active');
                activateTabsContent(tab);
            });
        })

        tabs[0].click();
    }

    const attachUpDownScroll = () => {
        const upDownButton = document.getElementsByClassName('btn-up-down')[0];
        const icon = upDownButton.getElementsByClassName('fas')[0];
        upDownButton.addEventListener('click', function () {
            if (icon.classList.contains('fa-caret-down')) {
                fc.scrollToBottom();
            } else {
                fc.scrollToTop();
            }
        });
        window.onscroll = function () {
            if (document.body.scrollTop > 1000 || document.documentElement.scrollTop > 1000) {
                icon.classList.remove('fa-caret-down');
                icon.classList.add('fa-caret-up');
            } else {
                icon.classList.remove('fa-caret-up');
                icon.classList.add('fa-caret-down');
            }
        }
    }

    const hasDataParent = (element) => {
        return element.hasAttribute('data-parent');
    }

    /** score calculator **/
    const calculateScoreSubtotal = ($rowElement) => {
        if (!hasDataParent($rowElement)) {
            return
        }
        const $tbody = $rowElement.closest('tbody');
        const dParent = $rowElement.getAttribute('data-parent');
        let scoreSubtotalElement = $tbody.querySelector('.score-' + dParent);
        if (!fc.elementExist(scoreSubtotalElement)) {
            return;
        }
        scoreSubtotalElement = scoreSubtotalElement.querySelector('.score-subtotal');
        if (!fc.elementExist(scoreSubtotalElement)) {
            return;
        }
        const scoreElements = $tbody.querySelectorAll('.ch-' + dParent);
        let scoreSubtotal = 0;
        for (let i = 0; i < scoreElements.length; i++) {
            let score = scoreElements[i].querySelector('.score').innerText;
            if (score.length) {
                scoreSubtotal += parseFloat(score);
            }
        }
        scoreSubtotalElement.innerText = scoreSubtotal.toFixed(2);
    };

    const calculateRating = (score) => {
        switch (true) {
            case score < 1.5:
                return ratings[0];
            case score >= 1.5 && score < 2.5:
                return ratings[1];
            case score >= 2.5 && score < 3.5:
                return ratings[2];
            case score >= 3.5 && score < 4.5:
                return ratings[3];
            case score >= 4.5:
                return ratings[4];
            default:
                return "";
        }
    }

    const calculateScoreTotal = ($rowElement) => {
        const $tbody = $rowElement.closest('tbody');
        const scoreElements = $tbody.querySelectorAll('.score');
        let scoreTotal = 0;
        for (let i = 0; i < scoreElements.length; i++) {
            let score = scoreElements[i].innerText;
            if (score.length) {
                scoreTotal += parseFloat(score);
            }
        }
        $tbody.querySelector('.score-total').innerText = scoreTotal.toFixed(2);
        $tbody.querySelector('.score-rating').innerText = calculateRating(scoreTotal);
    };

    const checkmarks = document.getElementsByClassName('rating-input');

    const calculate = ($c) => {
        const calc = ($chk) => {
            const $tr = $chk.closest('tr');
            const weightElement = $tr.querySelector('.weight');
            const scoreElement = $tr.querySelector('.score');
            const value = $chk.value;
            const weight = weightElement.innerText.replace(/[^0-9]/, '');
            scoreElement.innerText = ((weight / 100) * value).toFixed(2);
            calculateScoreSubtotal($tr);
            calculateScoreTotal($tr);
        }
        if (fc.elementExist($c)) {
            calc($c);
            return
        }
        for (let i = 0; i < checkmarks.length; i++) {
            if (!checkmarks[i].checked) {
                continue
            }
            calc(checkmarks[i]);
        }
    }

    const attachScoreCalculator = () => {
        for (let i = 0; i < checkmarks.length; i++) {
            checkmarks[i].addEventListener('change', function () {
                calculate(this);
            }, false);
        }
    }

    const init = () => {
        //populate rating guide
        const ratingGuideElements = document.querySelectorAll('.rating-guide .tag');
        if (!fc.elementExist(ratingGuideElements)) {
            return;
        }
        for (let i = 0; i < ratingGuideElements.length; i++) {
            ratingGuideElements[i].innerText = ratings[i];
        }
        calculate();
    }

    const attachOnFormSubmit = () => {
        const form360Element = document.getElementById('form360');
        const submitReview = document.getElementById('submit-review');
        const submitDraft = document.getElementById('save-draft');
        if (!fc.elementExist(submitReview) || !fc.elementExist(submitDraft)) {
            return;
        }
        const enableButton = (v) => {
            if (v) {
                submitReview.disabled = false;
                submitDraft.disabled = false;
                return;
            }
            submitReview.disabled = true;
            submitDraft.disabled = true;
        }
        const submitData = (type, data) => {
            enableButton(false);
            data['submission_type'] = type;
            fc.call(
                'r360',
                form360Element.getAttribute('method'),
                form360Element.getAttribute('action'),
                JSON.stringify(data),
                function (res) {
                    if (res.status === 200) {
                        fc.toast('Your review has been submitted successfully!', 'is-success');
                        setTimeout(function () {
                            location.replace('/mbr/review/list');
                        }, 700);
                    } else {
                        fc.toast(res.error.message, 'is-danger');
                        enableButton(true);
                    }
                }, function (fail) {
                    fc.toast('Your review submission did not complete successfully!', 'is-warning');
                    enableButton(true);
                }, function (aborted) {
                    enableButton(true);
                }
            );
        }
        submitReview.onclick = function (e) {
            e.preventDefault();
            e.stopPropagation();
            fc.overrideEvents('onClickOk', function () {
                submitData('final', fc.getFormData(form360Element));
            });
            fc.dialog(
                true, 'Submit Review',
                "Are you sure you want to submit your FINAL REVIEW now? Once you've submit, you no longer able to modify it.",
                [['ok', 'Yes'], ['cancel', 'No']]
            );
            return false;
        }
        submitDraft.onclick = function (e) {
            e.preventDefault();
            e.stopPropagation();
            fc.overrideEvents('onClickOk', function () {
                submitData('draft', fc.getFormData(form360Element));
            });
            fc.dialog(
                true, 'Saving as Draft', 
                'Are you sure you want to submit your DRAFT REVIEW now?',
                [['ok', 'Yes'], ['cancel', 'No']]
            );
            return false;
        }
    }

    init();
    attachOnFormSubmit();
    attachTabsFunctions();
    attachUpDownScroll();
    attachScoreCalculator();
})(fc, ratings);