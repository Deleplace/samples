let scene = document.getElementById("scene");

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
var swimmers = [];

function newSwimmer() {
    let s = document.createElement('img');
    s.classList.add('swimmer');
    s.classList.add('back');
    s.setAttribute('src', 'resources/swimmer.png');
    scene.appendChild(s);

    s.posX = xEntry;
    s.posY = 70 + Math.random() * heightPool;

    s.animate([
        { transform: 'translateX(' + 0 + 'px) translateY(' + s.posY + 'px)' },
        { transform: 'translateX(' + xEntry + 'px) translateY(' + s.posY + 'px)' }
    ],{ 
        duration: 1000 / speed,
        easing: "ease-out",
        fill: "forwards"
    });

    return s;
}

function goSwim(s) {
    takeCap();
    s.setAttribute('src', 'resources/swimmer-red-cap.png');
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
        { transform: 'translateX(' + xEntry + 'px) translateY(' + s.posY + 'px)' },
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
        putCap();
        s.setAttribute('src', 'resources/swimmer.png');
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
    const baseY = 550;
    for(let i=0;i<C;i++){
        let cap = document.createElement("img");
        cap.classList.add("cap");
        cap.setAttribute("src", "resources/cap.png");
        let x = 18 * Math.floor(i/5);
        let y = 15 * (i % 5);
        cap.style.left = (baseX + x) + "px";
        cap.style.top = (baseY - y) + "px";
        scene.appendChild(cap);
        caps.push(cap);
    }
    console.log("Basket filled up with " + C + " caps");
}

function takeCap() {
    utilization++;
    updateCapDisplay();
}

function putCap() {
    utilization--;
    updateCapDisplay();
}

function updateCapDisplay() {
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