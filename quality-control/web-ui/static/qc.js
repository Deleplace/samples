let scene = document.getElementById("scene");
let flash = document.getElementById("flash");
let foot = document.getElementById("foot");
let mugshot = document.getElementById("mugshot");
let mugshotPic = document.getElementById("mugshot-pic");
let uploadingMugshot = document.getElementById("uploading-mugshot");
let gettingPrediction = document.getElementById("getting-prediction");
let threshold = document.getElementById("threshold");
foot.firing = false;
foot.rot = 0;
let sprites = [];

let t=0;
let speed = 40; // in Hz
let defectRate = 0.30;

// next moves forward 1 animation step
function next() {
    t++;
    if(t==1000)
        return; // END
    if(Math.floor(t/20)%2 == 0) {
        // Move things along the assembly line
        resetFoot();
        for(let i=0;i<sprites.length;i++){
            let a = sprites[i];
            a.age ++;
            switch(a.state) {
            case ROLLING:
                a.posx += 8.4;
                a.posy -= 5.55;
                if(a.age>120)
                    destroy(a);
            }
            draw(a);
        }
    }else{
        // Pause the assembly line
        let tt = t%20;
        if(tt==9)
            takePicture();

        for(let i=0;i<sprites.length;i++){
            let a = sprites[i];
            a.age ++;
            if(tt==0 && a.posx>220 && a.posx<300){
                // console.log(a.posx + " -> " + a.src);
                mugshotPic.src = a.src;
            }
            if( foot.firing && a.posx>300 && a.posx<500)
                foot.style.transform = 'rotate(-' + (35*tt) + 'deg)';
            if( foot.firing && a.posx>350 )
                a.state = DISCARDING;
            switch(a.state) {
            case DISCARDING:
                a.posx += 7;
                a.posy += 3;
                a.rot += 30;
                if(a.age>90)
                    destroy(a);
            }
            draw(a);
        }
    }
    purge();
    window.setTimeout(next, 1000/speed);
}

// Puts a single item on the assembly line, and set it into motion
function scenarioSingle(itemfile, ok) {
    let pic = document.createElement("img");
    pic.src = "static/" + itemfile;
    pic.classList.add("sprite");
    pic.style.height = '60';
    pic.age = 0;
    pic.ok = ok;
    pic.posx = 100;
    pic.posy = 490;
    pic.rot = 0;
    draw(pic);
    pic.state = ROLLING;

    scene.appendChild(pic);
    sprites.push(pic);
}

function currentThreshold(){
    // threshold.value is in 0..100
    // currentThreshold is in 0.0..1.0
    return threshold.value/100;
}

threshold.oninput = function(){
    document.getElementById("threshold-display-value").textContent = currentThreshold();
};

// State machine
const ROLLING = 1;
const DISCARDING = 4;

function draw(pic){
    pic.style.marginLeft = pic.posx + "px";
    pic.style.marginTop = pic.posy + "px";
    if(pic.rot != 0)
        pic.style.transform = 'rotate(' + pic.rot + 'deg)';
}

function takePicture(){
    flash.classList.remove("hidden");
    window.setTimeout(function(){
        flash.classList.add("hidden");
    }, 3000/speed);

    mugshot.style.marginLeft = "45px";
    mugshot.style.marginTop = "300px";
    mugshot.classList.remove("hidden");

    window,setTimeout(function(){
        uploadingMugshot.classList.remove("hidden");
    }, 6000/speed);

    let ok = (mugshotPic.src.indexOf("-defect") == -1);
    let score = simulateDefectScore(ok);
    let detectedDefect = score > currentThreshold();
    gettingPrediction.score = score;
    if(detectedDefect)
        gettingPrediction.innerHTML =
            score.toFixed(2) + " > " + currentThreshold() + "<br/>"
            + "score is above threshold <br/>"
            + "DEFECT detected";
    else
        gettingPrediction.innerHTML =
            score.toFixed(2) + " < " + currentThreshold() + "<br/>"
            + "score is below threshold";

    window,setTimeout(function(){
        gettingPrediction.classList.remove("ok");
        gettingPrediction.classList.remove("defect");
        if(detectedDefect)
            gettingPrediction.classList.add("defect");
        else
            gettingPrediction.classList.add("ok");
        gettingPrediction.classList.remove("hidden");
    }, 10000/speed);

    window,setTimeout(function(){
        foot.firing = detectedDefect;
    }, 20000/speed);

    window.setTimeout(function(){
        mugshot.classList.add("hidden");
        mugshotPic.src = "";
        uploadingMugshot.classList.add("hidden");
    }, 20000/speed);

    window.setTimeout(function(){
        gettingPrediction.classList.add("hidden");
        gettingPrediction.innerHTML = "";
    }, 35000/speed);
}

function resetFoot(){
    foot.style.transform = 'rotate(0deg)';
}

function destroy(pic) {
    // Just mark it and hide it, for now
    pic.destroyed = true;
    pic.classList.add("hidden");
}

function purge() {
    // Remove "destroyed" sprites
    let keep = [];
    for(let i=0;i<sprites.length;i++){
        let a = sprites[i];
        if(a.destroyed)
            a.parentNode.removeChild(a);
        else
            keep.push(a);
    }
    sprites = keep;
}

function clear() {
    for(let i=0;i<sprites.length;i++){
        let a = sprites[i];
        a.parentNode.removeChild(a);
    }
    sprites = [];
}

document.getElementById("reset").onclick = function(){
    document.location.reload();
}

function fullScenario(nItems) {
    clear();

    function addNewItem() {
        let defect = (Math.random() < defectRate);
        let ok = !defect;
        scenarioSingle(pickItem(ok), ok);
    }

    for(let i=0;i<nItems;i++)
        window.setTimeout(addNewItem, (i*40000)/speed);

    // Start animation loop
    next();
}

function pickItem(ok) {
    let filename = "item-"
    if(ok)
        filename += "ok-";
    else
        filename += "defect-";
    filename += Math.floor(Math.random() * 4);
    filename += ".png";
    return filename;
}

function simulateDefectScore(ok){
    if(ok){
        // Random low score for label "DEFECT"
        return 0.75 * Math.random();
    }else{
        // Random high score for label "DEFECT"
        return 0.25 + (0.75 * Math.random());
    }
    // This is a fake classifer
    // TODO cleverer algorithm
    // TODO plug on real AutoML Vision Predict API
}

fullScenario(20);
