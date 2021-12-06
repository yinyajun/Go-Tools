// https://codepen.io/KaioRocha/pen/MdvWmg?ref=devawesome.io
(function () {
    let checkbox = document.querySelector('input[name=mode]');
    let themeSwitch = (isDark) => {
        localStorage.setItem("theme", isDark ? "dark" : "default");
        if (isDark) {
            document.documentElement.setAttribute("data-theme", "dark")
        } else {
            document.documentElement.setAttribute("data-theme", "default")
        }
    };

    let trans = () => {
        document.documentElement.classList.add('transition');
        window.setTimeout(() => {
            document.documentElement.classList.remove('transition');
        }, 500);
    };

    document.addEventListener('DOMContentLoaded', function () {
        let c = localStorage.getItem("theme") || "default";
        checkbox.checked = (c === 'dark');
        themeSwitch(c === 'dark')
    });

    checkbox.addEventListener('change', function () {
        trans();
        themeSwitch(this.checked);
    });
})();