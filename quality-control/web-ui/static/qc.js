console.log("ok");

let scene = document.getElementById("scene");
let flash = document.getElementById("flash");
let foot = document.getElementById("foot");
foot.firing = false;
foot.rot = 0;
let sprites = [];
let t=0;
function next() {
    t++;
    if(t==1000)
        return; // END
    if(Math.floor(t/20)%2 == 0) {
        for(let i=0;i<sprites.length;i++){
            let a = sprites[i];
            a.age ++;
            switch(a.state) {
            case ROLLING:
                a.posx += 8.6;
                a.posy -= 5.8;
            }
            draw(a);
        }
    }else{
        if(t%20==9)
            flash.classList.remove("hidden");
        if(t%20==12)
            flash.classList.add("hidden");

        for(let i=0;i<sprites.length;i++){
            let a = sprites[i];
            a.age ++;
            if( !a.ok && a.posx>300)
                foot.style.transform = 'rotate(-' + (30*t) + 'deg)';
            if( !a.ok && a.posx>350)
                a.state = DISCARDING;
            switch(a.state) {
            case DISCARDING:
                a.posx += 7;
                a.posy += 5;
                a.rot += 30;
            }
            draw(a);
        }
    }
    window.setTimeout(next, 100);
}

function scenario(itemfile, ok) {
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

// State machine
const ROLLING = 1;
const DISCARDING = 4;

function draw(pic){
    pic.style.marginLeft = pic.posx + "px";
    pic.style.marginTop = pic.posy + "px";
    if(pic.rot != 0)
        pic.style.transform = 'rotate(' + pic.rot + 'deg)'; 
}

scenario("item-ok.png", true);
window.setTimeout(function(){
    scenario("item-defect.png", false);
}, 3000);
next();