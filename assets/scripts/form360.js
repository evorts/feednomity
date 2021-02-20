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

    attachTabsFunctions();
    attachUpDownScroll();
})();