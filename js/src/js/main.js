import saveScreenshot from './screenshot.js'
import img from '../img/prairie_dog.gif';

document.getElementById('img1').src = img;

function screenshot(elem, filename) {
    const el = document.querySelector(elem);
    saveScreenshot(el, filename);
}

document.querySelector('#btn1').onclick = function() {
    screenshot('#div1', 'shot1.png');
}

document.querySelector('#btn2').onclick = function() {
    screenshot('#div2', 'shot2.png');
}
