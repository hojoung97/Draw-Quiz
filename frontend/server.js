const http = require('http');
const fs = require('fs');

let html;
let css;
let js;

let drawAppHTML;
let drawAppCSS;
let drawAppJS;

fs.readFile('draw_app.html', function (err, data) {
  if (err) {
    throw err;
  }
  drawAppHTML = data;
});
fs.readFile('css/draw_app.css', function (err, data) {
  if (err) {
    throw err;
  }
  drawAppCSS = data;
});
fs.readFile('js/draw_app.js', function (err, data) {
  if (err) {
    throw err;
  }
  drawAppJS = data;
});

fs.readFile('index.html', function (err, data) {
  if (err) {
    throw err;
  }
  html = data;
});
fs.readFile('css/index.css', function (err, data) {
  if (err) {
    throw err;
  }
  css = data;
});
fs.readFile('js/index.js', function (err, data) {
  if (err) {
    throw err;
  }
  js = data;
});

http.createServer((req, res) => {
  res.statusCode = 200;
  if(req.url.indexOf('css/draw_app.css') != -1){
    res.writeHead(200, {'Content-Type': 'text/css'});
    res.write(drawAppCSS);
    res.end();
    return;
  }
  if(req.url.indexOf('js/draw_app.js') != -1){
    res.writeHead(200, {'Content-Type': 'text/javascript'});
    res.write(drawAppJS);
    res.end();
    return;
  }
  if(req.url.indexOf('draw_app.html') != -1){
    res.writeHeader(200, {"Content-Type": "text/html"});
    res.write(drawAppHTML);
    res.end();
    return;
  }
  if(req.url.indexOf('css/index.css') != -1){
   res.writeHead(200, {'Content-Type': 'text/css'});
   res.write(css);
   res.end();
   return;
  }
  if(req.url.indexOf('index.js') != -1){
   res.writeHead(200, {'Content-Type': 'text/javascript'});
   res.write(js);
   res.end();
   return;
  }
  res.writeHeader(200, {"Content-Type": "text/html"});
  res.write(html);
  res.end();
}).listen(3000);
