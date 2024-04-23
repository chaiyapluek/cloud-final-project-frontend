function loadding(event, text="Loading...") {
    var buttonText = document.querySelector(".load-button span");
    if (buttonText) {
        buttonText.setAttribute("disabled", "disabled");
        buttonText.innerText = text;
        var buttonSpinner = document.querySelector("svg");
        if (buttonSpinner) {
            buttonSpinner.classList.remove("hidden");
            buttonSpinner.classList.add("inline");
        }
    }
}

function stopLoadding(event, text) {
    var buttonText = document.querySelector(".load-button span");
    if (buttonText) {
        buttonText.removeAttribute("disabled");
        buttonText.innerText = text;
        var buttonSpinner = document.querySelector("svg");
        if (buttonSpinner) {
            buttonSpinner.classList.remove("inline");
            buttonSpinner.classList.add("hidden");
        }
    }
}

function showLoaddingModal(event, selector) {
    var e = document.querySelector(selector);
    if (e) {
        e.classList.remove("hidden");
        e.classList.add("inline");
    }
}

function hideLoaddingModal(event, selector) {
    var e = document.querySelector(selector);
    if (e) {
        e.classList.remove("inline");
        e.classList.add("hidden");
    }
}

function timer(seconds) {
    return {
        current: seconds,
        init() {
            setInterval(() => {
                this.setCurrent()
            }, 1000);
        },
        setCurrent() {
            if (this.current <= 0) {
                return
            }
            this.current -= 1
        }
    }
}