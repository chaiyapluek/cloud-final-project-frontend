function loadding(event) {
    var buttonText = document.querySelector(".load-button span");
    buttonText.setAttribute("disabled", "disabled");
    buttonText.innerText = "Loading...";
    var buttonSpinner = document.querySelector("svg");
    if(buttonSpinner){
        buttonSpinner.classList.remove("hidden");
        buttonSpinner.classList.add("inline");
    }    
}

function stopLoadding(event, text) {
    var buttonText = document.querySelector(".load-button span");
    buttonText.removeAttribute("disabled");
    buttonText.innerText = text;
    var buttonSpinner = document.querySelector("svg");
    if(buttonSpinner){
        buttonSpinner.classList.remove("inline");
        buttonSpinner.classList.add("hidden");
    }    
}

function showLoaddingModal(event, selector){
    var e = document.querySelector(selector);
    console.log("show", e);
    if(e){
        e.classList.remove("hidden");
        e.classList.add("inline");
    }
}

function hideLoaddingModal(event, selector){
    var e = document.querySelector(selector);
    console.log("hide", e);
    if(e){
        e.classList.remove("inline");
        e.classList.add("hidden");
    }
}

function timer(seconds){
    return {
        current: seconds,
        init() {
            setInterval(()=>{
                this.setCurrent()
            }, 1000);
        },
        setCurrent() {
            if(this.current <= 0){
                return
            }
            this.current -= 1
        }
    }
}