// Machine generated. Do not edit. Go read ../generator/README.md
package server

var Resources map[string]string = map[string]string{
	"/index.html":    "<html>\n<head>\n<script type=\"text/javascript\" src=\"js/example.js\"></script>\n</head>\n<body>\n<p>Hello from sfsbook!</p>\n<button onclick=\"addSomething();\">Click Me Now!</button>\n<div id=\"insertPoint\"></div>\n</body>\n</html>\n",
	"/js/example.js": "// Example JavaScript code for sfsbook\nfunction addSomething() {\n\tvar el = document.getElementById(\"insertPoint\");\n\tvar p = document.createElement('p');\n\tp.innerText = \"Something here. But not very exciting\";\n\tel.appendChild(p)\n}\n",
}
