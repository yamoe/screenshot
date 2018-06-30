import html2canvas from 'html2canvas';
import FileSaver from 'file-saver';
import 'canvas-toBlob';
import moment from 'moment';

const prefix = {
    svg: 'http://www.w3.org/2000/svg',
    xhtml: 'http://www.w3.org/1999/xhtml',
    xlink: 'http://www.w3.org/1999/xlink',
    xml: 'http://www.w3.org/XML/1998/namespace',
    xmlns: 'http://www.w3.org/2000/xmlns/',
};

class Screenshot {

    constructor() {
        this.el = null;
        this.filename = '';

        this.nodesToRecover = [];
        this.nodesToRemove = [];
        this.svgsLength = 0;
    }

    save(el, filename) {
        this.el = el;
        this.filename = filename;

        if (el.nodeName === 'svg' && el.nodeType === 1) {
            const cb = (function (canvas) {
                this.inDownloadFile(canvas, filename);
            }).bind(this);

            this.inConvertSvgToCanvas(svg, cb);
            return;
        }

        const svgs = el.querySelectorAll('svg');
        this.svgsLength = svgs.length;

        if (this.svgsLength > 0) {
            const cb = this.inOnCanvas.bind(this);
            for (const svg of svgs) {
                this.inConvertSvgToCanvas(svg, cb, svg);
            }
        } else {
            this.inCapture();
        }
    }

    inConvertSvgToCanvas(svg, cb, arg) {
        this.inSvgPreprocess(svg);

        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');
        const image = new Image();
        image.onload = function load() {
            canvas.height = image.height;
            canvas.width = image.width;
            ctx.fillStyle = '#FFFFFF';
            ctx.fillRect(0, 0, image.width, image.height);
            ctx.drawImage(image, 0, 0);
            cb(canvas, arg);
        };
        image.src = 'data:image/svg+xml;charset-utf-8,' + encodeURIComponent(svg.outerHTML);
    }

    inDownloadFile(canvas, filename) {
        const fn = filename
            .replace(/[&/\\#,+()$~%.'":*?<>{}]/g, '')
            .replace(/ /g, '_')
            .toLowerCase()
            + '_' + moment().format('YYYYMMDD_HHmmss') + '.png';

        canvas.toBlob(function (blob) {
            FileSaver.saveAs(blob, fn);
        });
    }

    inOnCanvas(canvas, svg) {
        const parentNode = svg.parentNode;

        this.nodesToRecover.push({ parent: parentNode, child: svg });
        parentNode.removeChild(svg);

        this.nodesToRemove.push({ parent: parentNode, child: canvas });
        parentNode.appendChild(canvas);

        if (--this.svgsLength <= 0) {
            this.inCapture();
        }
    }

    inCapture() {
        const cb = (function (canvas) {
            const ctx = canvas.getContext('2d');
            ctx.webkitImageSmoothingEnabled = false;
            ctx.mozImageSmoothingEnabled = false;
            ctx.imageSmoothingEnabled = false;

            for (const pair of this.nodesToRemove) {
                pair.parent.removeChild(pair.child);
            }
            for (const pair of this.nodesToRecover) {
                pair.parent.appendChild(pair.child);
            }

            this.inDownloadFile(canvas, this.filename);

        }).bind(this);

        html2canvas(this.el, {
            logging: false,
            allowTaint: true,
            useCORS: true,
        }).then(cb);
    }

    inSvgExplicitlySetStyle(element, emptySvgDeclarationComputed) {
        const cSSStyleDeclarationComputed = window.getComputedStyle(element);
        let i;
        let len;
        let key;
        let value;
        let computedStyleStr = '';

        for (i = 0, len = cSSStyleDeclarationComputed.length; i < len; i++) {
            key = cSSStyleDeclarationComputed[i];
            value = cSSStyleDeclarationComputed.getPropertyValue(key);
            if (value !== emptySvgDeclarationComputed.getPropertyValue(key)) {
                if ((key !== 'height') && (key !== 'width')) {
                    computedStyleStr += key + ':' + value + ';';
                }

            }
        }
        element.setAttribute('style', computedStyleStr);
    }

    inVisitAll(node, tree) {
        if (node && node.hasChildNodes()) {
            let child = node.firstChild;
            while (child) {
                if (child.nodeType === 1 && child.nodeName !== 'SCRIPT') {
                    tree.push(child);
                    this.inVisitAll(child, tree);
                }
                child = child.nextSibling;
            }
        }
    }

    inGetAllElements(obj) {
        const tree = [];
        tree.push(obj);
        this.inVisitAll(obj, tree);
        return tree;
    }

    inSvgPreprocess(svg) {
        svg.setAttribute('version', '1.1');

        // removing attributes so they aren't doubled up
        svg.removeAttribute('xmlns');
        svg.removeAttribute('xlink');

        // These are needed for the svg
        if (!svg.hasAttributeNS(prefix.xmlns, 'xmlns')) {
            svg.setAttributeNS(prefix.xmlns, 'xmlns', prefix.svg);
        }

        if (!svg.hasAttributeNS(prefix.xmlns, 'xmlns:xlink')) {
            svg.setAttributeNS(prefix.xmlns, 'xmlns:xlink', prefix.xlink);
        }

        // svg set inline style
        // add empty svg element
        const emptySvg = window.document.createElementNS(prefix.svg, 'svg');
        window.document.body.appendChild(emptySvg);
        const emptySvgDeclarationComputed = window.getComputedStyle(emptySvg);

        // hardcode computed css styles inside svg
        const allElements = this.inGetAllElements(svg);
        let i = allElements.length;
        while (i--) {
            this.inSvgExplicitlySetStyle(allElements[i], emptySvgDeclarationComputed);
        }
        emptySvg.parentNode.removeChild(emptySvg);
    }

}

export default function saveScreenshot(el, filename) {
    if (!el) return;
    const shot = new Screenshot();
    shot.save(el, filename);
}
