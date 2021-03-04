(function () {
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

    const scrollToBottom = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = elem.scrollHeight;
    }

    const scrollToTop = () => {
        const elem = document.scrollingElement || document.body;
        elem.scrollTop = 0;
    }

    const attachUpDownScroll = () => {
        const upDownButton = document.getElementsByClassName('btn-up-down')[0];
        const icon = upDownButton.getElementsByClassName('fas')[0];
        upDownButton.addEventListener('click', function () {
            if (icon.classList.contains('fa-caret-down')) {
                scrollToBottom();
            } else {
                scrollToTop();
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

    const hasClass = (element, selector) => {
        return element.classList.contains(selector);
    }

    const hasDataParent = (element) => {
        return element.hasAttribute('data-parent');
    }

    const elementExist = (element) => {
        return element != null && typeof element != "undefined";
    }
    const getSiblings = function (e, c) {
        // for collecting siblings
        let siblings = [];
        // if no parent, return no sibling
        if(!e.parentNode) {
            return siblings;
        }
        // first child of the parent node
        let sibling  = e.parentNode.firstChild;

        // collecting siblings
        while (sibling) {
            if (sibling.nodeType === 1 && sibling !== e) {
                if (c == null) {
                    siblings.push(sibling);
                } else {
                    if (hasClass(sibling, c)) {
                        siblings.push(sibling);
                    }
                }
            }
            sibling = sibling.nextSibling;
        }
        return siblings;
    };

    const attachScoreCalculator = () => {
        const checkmarks = document.getElementsByClassName('rating-input');
        const calculateScoreSubtotal = ($rowElement) => {
            if (!hasDataParent($rowElement)) {
                return
            }
            const $tbody = $rowElement.closest('tbody');
            const dParent = $rowElement.getAttribute('data-parent');
            let scoreSubtotalElement = $tbody.querySelector('.score-'+dParent);
            if (!elementExist(scoreSubtotalElement)) {
                return;
            }
            scoreSubtotalElement = scoreSubtotalElement.querySelector('.score-subtotal');
            if (!elementExist(scoreSubtotalElement)) {
                return;
            }
            const scoreElements = $tbody.querySelectorAll( '.ch-'+dParent);
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
                    return "Pretty Bad";
                case score >= 1.5 && score < 2.5:
                    return "Need Improvement";
                case score >= 2.5 && score < 3.5:
                    return "Meet Expectation";
                case score >= 3.5 && score < 4.5:
                    return "Outstanding";
                case score >= 4.5:
                    return "Excellent";
                default:
                    return "";
            }
        }
        const calculateScoreTotal = ($rowElement) => {
            const $tbody = $rowElement.closest('tbody');
            const scoreElements = $tbody.querySelectorAll( '.score');
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

        for (let i = 0; i < checkmarks.length; i++) {
            checkmarks[i].addEventListener('change', function () {
                const $tr = this.closest('tr');
                const weightElement = $tr.querySelector('.weight');
                const scoreElement = $tr.querySelector('.score');
                const value = this.value;
                const weight = weightElement.innerText.replace(/[^0-9]/,'');
                scoreElement.innerText = ((weight / 100) * value).toFixed(2);
                calculateScoreSubtotal($tr);
                calculateScoreTotal($tr);
            }, false);
        }
    }

    attachTabsFunctions();
    attachUpDownScroll();
    attachScoreCalculator();
})();