var closesIcon = document.querySelectorAll('.xd-message .close-icon');

closesIcon.forEach(function(closeIcon) {
    closeIcon.addEventListener('click', function() {
        this.parentNode.parentNode.classList.add('hide');
    });
});
