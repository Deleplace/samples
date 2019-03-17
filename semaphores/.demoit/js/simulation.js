let scene = document.getElementById("scene");
let safetyIndicator = document.getElementById("safety-indicator");

const xEntry = 120;
const xPool = 350;
const widthPool = 500;
const heightPool = 200;
const yOut = 400;

// Simulation
var N = 200; // Total number of swimmers
var capacity = 20; // Max capacity of the pool
var speed = 1;
var utilization = 0; // Current number of swimmers in the pool
var caps;
var bags;
var swimmers = [];
var metaphor; // "swimcaps" or "gymbags" or "nosync"

function newSwimmer() {
    let s = document.createElement('img');
    s.classList.add('swimmer');
    s.classList.add('back');
    if(metaphor=="gymbags")
        s.setAttribute('src', '/images/swimmer-with-bag.png');
    else
        s.setAttribute('src', '/images/swimmer.png');
    scene.appendChild(s);

    s.posX = xEntry - 50 + 70*Math.random();
    s.posY = 70 + Math.random() * heightPool;

    s.animate([
        { transform: 'translateX(' + 0 + 'px) translateY(' + s.posY + 'px)' },
        { transform: 'translateX(' + s.posX + 'px) translateY(' + s.posY + 'px)' }
    ],{ 
        duration: 1000 / speed,
        easing: "ease-out",
        fill: "forwards"
    });

    return s;
}

function goSwim(s) {
    utilization++;
    updateSafetyIndicator();
    putGymBag();
    takeCap();
    if(metaphor=="swimcaps")
        s.setAttribute('src', '/images/swimmer-red-cap.png');
    else
        s.setAttribute('src', '/images/swimmer.png');
    let oldX = s.posX;
    let newX = xPool + widthPool * Math.random();
    let newY = s.posY;
    if(utilization > capacity) {
        // Ouch.
        console.log("OVER CAPACITY: " + utilization);
        newY -= 1 * (utilization-capacity);
        // newY -= 100;
        // newX = xPool + (widthPool/4) + (widthPool * Math.random() / 2);
        let minX = xPool + 2*(utilization-capacity);
        let maxX = xPool + widthPool - 2*(utilization-capacity);
        newX = minX + (maxX - minX) * Math.random();
    }
    s.animate([
        { transform: 'translateX(' + oldX + 'px) translateY(' + s.posY + 'px)' },
        { transform: 'translateX(' + newX + 'px) translateY(' + newY + 'px)' }
    ],{ 
        duration: 1000 / speed,
        easing: "ease-in-out",
        fill: "forwards"
    });
    s.posX = newX;
    s.posY = newY;
}

function getOut(s) {
    // Dramatic hack: nobody gets out anymore when congestion is high
    if (utilization > 2*capacity)
        return;

    s.classList.add('back');
    let d = 3000 / speed;
    let anim = s.animate([
        { transform: 'translateX(' + s.posX + 'px) translateY(' + s.posY + 'px) scaleX(-1)' },
        { transform: 'translateX(' + s.posX + 'px) translateY(' + yOut + 'px) scaleX(-1)' },
        { transform: 'translateX(' + xEntry + 'px) translateY(' + yOut + 'px) scaleX(-1)' }
    ],{ 
        duration: d,
        easing: "ease-out",
        fill: "forwards"
    });
    s.posX = xEntry;

    window.setTimeout(function() {
        utilization--;
        updateSafetyIndicator();
        putCap();
        takeGymBag();
        if(metaphor=="gymbags")
            s.setAttribute('src', '/images/swimmer-with-bag.png');
        else
            s.setAttribute('src', '/images/swimmer.png');
        s.animate([
            { transform: 'translateX(' + xEntry + 'px) translateY(' + yOut + 'px) scaleX(-1)' },
            { transform: 'translateX(' + -200 + 'px) translateY(' + yOut + 'px) scaleX(-1)' },
        ],{
            duration: 4000 / speed,
            easing: "ease-out",
            fill: "forwards"
        });
    }, d);
}

function makeBasketCaps(C) {
    if(caps){
        for(let i=0;i<caps.length;i++)
            document.removeChild(caps[i]);
    }
    caps = [];
    const baseX = 130;
    const baseY = 350;
    for(let i=0;i<C;i++){
        let cap = document.createElement("img");
        cap.classList.add("cap");
        cap.setAttribute("src", "/images/cap.png");
        let x = 18 * Math.floor(i/5);
        let y = 15 * (i % 5);
        cap.style.left = (baseX + x) + "px";
        cap.style.top = (baseY - y) + "px";
        scene.appendChild(cap);
        caps.push(cap);
    }
    console.log("Basket filled up with " + C + " caps");
}

function makeGymbagsShelf(C) {
    if(bags){
        for(let i=0;i<caps.length;i++)
            document.removeChild(caps[i]);
    }
    bags = [];
    const baseX = 50;
    const baseY = 348;
    for(let i=0;i<C;i++){
        let bag = document.createElement("img");
        bag.classList.add("gymbag");
        bag.setAttribute("src", "/images/gymbag.png");
        bag.style.display = "none";
        let x = 49 * i;
        bag.style.left = (baseX + x) + "px";
        bag.style.top = (baseY) + "px";
        scene.appendChild(bag);
        bags.push(bag);
    }
    console.log("Shelf ready for " + C + " gym bags");
}

function takeCap() {
    updateCapDisplay();
}

function putCap() {
    updateCapDisplay();
}

function takeGymBag() {
    updateGymbagsDisplay();
}

function putGymBag() {
    updateGymbagsDisplay();
}

function updateCapDisplay() {
    if(metaphor != "swimcaps")
        return;
    console.log("utilization="+utilization);
    let stackSize = capacity - utilization;
    if(utilization>capacity){
        console.warn("Yes there is a race condition between GopherJS controller and the JS view. But it's not the topic today.")
        stackSize = 0;
    }
    for(let i=0;i<stackSize;i++)
        caps[i].style.display = "block";
    for(let i=stackSize;i<capacity;i++)
        caps[i].style.display = "none";
}

function updateGymbagsDisplay() {
    if(metaphor != "gymbags")
        return;
    console.log("utilization="+utilization);
    let stackSize = utilization;
    if(utilization<0){
        console.warn("Yes there is a race condition between GopherJS controller and the JS view. But it's not the topic today.")
        stackSize = 0;
    }
    if(utilization>capacity){
        console.warn("Yes there is a race condition between GopherJS controller and the JS view. But it's not the topic today.")
        stackSize = capacity-1;
    }
    for(let i=0;i<stackSize;i++)
        bags[i].style.display = "block";
    for(let i=stackSize;i<capacity;i++)
        bags[i].style.display = "none";
}

function updateSafetyIndicator() {
    if(!safetyIndicator)
        return;

    if(utilization>capacity) {
        safetyIndicator.src = "/images/safety-red.png";
        return;
    }
    /*
    if(utilization>=(9*capacity/10)) {
        safetyIndicator.src = "/images/safety-orange.png";
        return;
    }
    */
    safetyIndicator.src = "/images/safety-green.png";
}

//
// To be called by the Go code:
//
function arrive(i) {
    swimmers[i] = newSwimmer();
}

function swim(i, d) {
    let s = swimmers[i];
    console.log("goSwim(" + i + ")");
    goSwim(s);
    window.setTimeout(function(){
        console.log("getOut(" + i + ")");
        getOut(s);
    }, d/speed);
}