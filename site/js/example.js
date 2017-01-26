// Example JavaScript code for sfsbook
function addSomething() {
	var el = document.getElementById("insertPoint");
	var p = document.createElement('p');
	p.innerText = "Something here. But not very exciting";
	el.appendChild(p)
}

// Delete the session cookie and reload the page in the unauthed state.
function clearSessionCookie(e) {
	if((e.which === 13) || (e.which === 1)) {
		console.log('event', e.which === 1);
	  document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/";
	  window.location.reload();
	}
}

// Enable visible keyboard navigation of user menu
function dropdownNav(e, node) {
	var theTarget = event.target;

if(e.which === 40){
  	if(theTarget.classList[0].indexOf('user-menu-item') > -1) {
  		var siblings = Array.prototype.slice.call(document.querySelectorAll(".user-menu-item"));
  		var current = siblings.indexOf(theTarget);
  		var next = (current === siblings.length - 1) ? 0 : current + 1;
  		var nextItem = document.querySelectorAll(".user-menu-item")[next];
  		nextItem.focus();
  	} else if (theTarget.classList && theTarget.classList[0] === 'user-menu') {
  	  var child = document.querySelector(".user-menu-dropdown");
  		  console.log('child classlist', child.classList.toString());
  		  child.classList.toggle("visible");
  		  console.log('child classlist toggled', child.classList.toString());
  		  document.querySelector(".user-menu-item").focus();
  	}
  } else if (e.which === 13 && theTarget.classList && theTarget.classList[0].indexOf('user-menu-item') > -1){
    location.href=theTarget.parentElement.href;
  } else if (e.which === 27 && theTarget.classList && theTarget.classList[0] === 'user-menu-item'){
    var parent = document.querySelector('.user-menu');
    var child = document.querySelector(".user-menu-dropdown");
      		  console.log('child classlist', child.classList.toString());
              if(child.classList[1].indexOf('visible') > -1) {
              	child.classList.toggle("visible");
      		    console.log('child classlist toggled', child.classList.toString());
                parent.focus();
              }
  } 


}

// Enable visible keyboard navigation of search filters
function searchNav(e, node) {
	var theTarget = event.target;

if(e.which === 40){
  	if(theTarget.classList[0].indexOf('category') > -1) {
  		var siblings = Array.prototype.slice.call(document.querySelectorAll(".category"));
  		var current = siblings.indexOf(theTarget);
  		var next = (current === siblings.length - 1) ? 0 : current + 1;
  		var nextItem = document.querySelectorAll(".category")[next];
  		nextItem.focus();
  	} else if (theTarget.classList && theTarget.classList[0] === 'search-filter-container') {
  	  var child = document.querySelector('.search-filters');
  		  console.log('child classlist', child.classList.toString());
  		  child.classList.toggle("visible");
  		  console.log('child classlist toggled', child.classList.toString());
  		  console.log(document.querySelector(".category"));
  		  document.querySelector(".category").focus();
  	}
  }/* else if (e.which === 13 && theTarget.classList && theTarget.classList[0].indexOf('user-menu-item') > -1){
    location.href=theTarget.parentElement.href;
  } else if (e.which === 27 && theTarget.classList && theTarget.classList[0] === 'user-menu-item'){
    var parent = document.querySelector('.user-menu');
    var child = document.querySelector(".user-menu-dropdown");
      		  console.log('child classlist', child.classList.toString());
              if(child.classList[1].indexOf('visible') > -1) {
              	child.classList.toggle("visible");
      		    console.log('child classlist toggled', child.classList.toString());
                parent.focus();
              }
  } */


}
