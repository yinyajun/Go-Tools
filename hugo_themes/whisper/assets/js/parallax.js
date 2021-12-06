let parallax = (el) => {
    let height = el.offsetHeight;
    let offset = window.scrollY;
    if (offset > height) {
        return;
    }
    el.style.transform = 'translate3d(0,' + offset / 4 + 'px,0)'
};

let slideParallax = () => {
    let el = document.getElementById("slide-img");
    parallax(el)
};

(function () {
    window.addEventListener("scroll", throttle(slideParallax, 15))
})();