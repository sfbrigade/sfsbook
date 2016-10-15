// Example JavaScript code for sfsbook
function addSomething() {
	var el = document.getElementById("insertPoint");
	var p = document.createElement('p');
	p.innerText = "Something here. But not very exciting";
	el.appendChild(p)
}

// Delete the session cookie and reload the page in the unauthed state.
function clearSessionCookie() {
	document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/";
	window.location.reload();
}
