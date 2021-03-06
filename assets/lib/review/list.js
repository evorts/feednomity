(function (fc, apiUrl) {
    const createRowElement = (id, no, title, recipient, status) => {
        const $row = document.createElement('tr');
        const $no = document.createElement('td');
        const $title = document.createElement('td');
        const $recipient = document.createElement('td');
        const $status = document.createElement('td');
        const $statusSpan = document.createElement('span');
        const $action = document.createElement('td');
        const $btnReview = document.createElement('button');

        fc.addClass($row, 'item');
        fc.addClass($btnReview, ['button','is-link','is-small']);
        fc.addClass($statusSpan, 'tag');

        if (status !== 'final') {
            $btnReview.setAttribute('data-id', id);
            $btnReview.append('Give Review');
            $btnReview.onclick = function (e) {
                e.stopPropagation();
                window.open(`/mbr/reviews/${this.getAttribute('data-id')}`, '_blank').focus();
            }
            $action.appendChild($btnReview);
            fc.addClass($statusSpan, status === 'draft' ? 'is-warning' : 'is-danger');
        } else {
            fc.addClass($statusSpan, 'is-success');
        }

        $no.append(`${no}`);
        $title.append(`${title}`);
        $recipient.append(`${recipient}`);
        $statusSpan.append(`${status}`)
        $status.append($statusSpan);

        $row.appendChild($no);
        $row.appendChild($title);
        $row.appendChild($recipient);
        $row.appendChild($status);
        $row.appendChild($action);

        return $row;
    }

    const $containerLoadMore = document.querySelectorAll('.table.items tfoot.load-more');
    const $loadMoreButton = $containerLoadMore[0].querySelector('.button');
    const $itemsContainerElement = document.querySelectorAll('.table.items tbody');

    const populateItems = () => {
        const limit = 10;
        let page = 1;

        const $displayedItemsElement = $itemsContainerElement[0].querySelectorAll('.item');
        if (fc.elementExist($displayedItemsElement)) {
            page = ($displayedItemsElement.length / limit) + 1;
        }
        fc.call(
            'review-list',
            'POST',
            `${apiUrl}/v1/reviews/list`,
            JSON.stringify({page: page, limit: limit}),
            function (res) {
                if (res.status !== 200) {
                    return;
                }
                if (res.content.items.length < limit) {
                     $containerLoadMore[0].remove();
                }
                if (res.content.items.length < 1) {
                    return;
                }
                res.content.items.forEach(function (elem, idx) {
                    const id = elem['id'];
                    const no = ((page - 1) * limit) + idx + 1;
                    const title = `${elem['distribution_topic']}`;
                    const recipient = `${elem['recipient_name']}`;
                    const status = `${elem['status']}`;
                    $itemsContainerElement[0].appendChild(createRowElement(id, no, title, recipient, status));
                });
                $loadMoreButton.removeAttribute('disabled');
            }
        )
    }

    fc.onDocumentReady(function () {
        populateItems();
        $loadMoreButton.addEventListener('click', function (e) {
            e.stopPropagation();
            this.setAttribute('disabled', "disabled");
            populateItems();
        });
    });
})(fc, ApiBaseUrl);