/*
Designed by: Oleg Frolov
Original image : https://dribbble.com/shots/6629408-Info-Button-Interaction
*/

const btn = document.getElementById("btn");
const second = document.getElementById("second");

let showCard = (event) => {
    btn.classList.toggle("is-rotate");
    second.classList.toggle("is-visible");
}

btn.addEventListener("click", showCard);

function CopyToClipboard(id)
{
    var r = document.createRange();
    r.selectNode(document.getElementById(id));
    window.getSelection().removeAllRanges();
    window.getSelection().addRange(r);
    document.execCommand('copy');
    window.getSelection().removeAllRanges();
}