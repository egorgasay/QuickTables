let circularProgress = document.querySelector(".circular-progress"),
    progressValue = document.querySelector(".progress-value");

let progressStartValue = 0,
    progressEndValue = 100,
    speed = 120;


const section = document.querySelector('section'),
    overlay = document.querySelector('.overlay'),
    showBtn = document.querySelector('.show-modal'),
    closeBtn = document.querySelector('.close-btn');

showBtn.addEventListener('click', () => {
    section.classList.add('active');
    let progress = setInterval(() => {
        progressStartValue++;

        progressValue.textContent = `${progressStartValue}%`;
        circularProgress.style.background = `conic-gradient(#456990 ${
            progressStartValue * 3.6
        }deg, #ededed 0deg)`;

        if (progressStartValue == progressEndValue) {
            clearInterval(progress);
            window.location = '/';
        }
    }, speed);
});

closeBtn.addEventListener('click', () => {
    section.classList.remove('active');
});
overlay.addEventListener('click', () => {
    section.classList.remove('active');
});
