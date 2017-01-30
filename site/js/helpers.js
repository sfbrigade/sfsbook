// toggleCategoryOption adds or removes a keyword from the search
function toggleCategoryOption(node) {
  var el = document.getElementById('query_field');
  var searchText = el.value;
  var input = (node.type === 'click') ? node.target : node;
  var optionValue = input.textContent;
  var re = new RegExp('\\b(' + optionValue + ')\\b', 'gi');
  var stringToReplace = optionValue.concat(', ');
  if(!searchText.match(stringToReplace) && !searchText.match(re)) {
    if (searchText.length > 0) {
      el.value = searchText.concat(', ', optionValue);
    } else {
      el.value = optionValue;
    }
  }
}

// Delete the session cookie and reload the page in the unauthed state.
function clearSessionCookie() {
  if((event.which === 13) || (event.which === 1)) {
    document.cookie = "session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/";
    window.location.reload();
  }
}

// Enable visible keyboard navigation of user menu
function userNav(node) {
  var theTarget = event.target;
  if(event.which === 40){
    if(theTarget.classList[0].indexOf('user-menu-item') > -1) {
      var siblings = Array.prototype.slice.call(document.querySelectorAll(".user-menu-item"));
      var current = siblings.indexOf(theTarget);
      var next = (current === siblings.length - 1) ? 0 : current + 1;
      var nextItem = document.querySelectorAll(".user-menu-item")[next];
      nextItem.focus();
    } else if (theTarget.classList && theTarget.classList[0] === 'user-menu') {
      var child = document.querySelector(".user-menu-dropdown");
      child.classList.toggle("visible");
      document.querySelector(".user-menu-item").focus();
    }
  } else if (event.which === 13 && theTarget.classList && theTarget.classList[0].indexOf('user-menu-item') > -1){
    location.href=theTarget.parentElement.href;
  } else if (event.which === 27 && theTarget.classList && theTarget.classList[0] === 'user-menu-item'){
    var parent = document.querySelector('.user-menu');
    var child = document.querySelector(".user-menu-dropdown");
    if(child.classList[1] && (child.classList[1].indexOf('visible') > -1)) {
      child.classList.toggle("visible");
      parent.focus();
    }
  } 
}

// Enable visible keyboard navigation of search filters
function searchNav(node) {
  var theTarget = event.target;
  if(event.which === 9){
    if((theTarget.classList[0].indexOf('category') > -1) || (theTarget.classList[0].indexOf('search-filter-container') > -1)) {
      document.querySelector('.search-filters').classList.remove('visible');
      var secondlevel = document.querySelectorAll('.category-dropdown');
      for(var i = 0; i < secondlevel.length; i++){
        secondlevel[i].classList.remove('visible');
      }
    }
  }
  if(event.which === 27 || event.which === 37){
    if(theTarget.classList[0].indexOf('category') > -1) {
      theTarget.parentElement.classList.toggle('visible');
      theTarget.parentElement.parentElement.focus();
	}
  }
  if(event.which === 13){
    if(theTarget.classList[0].indexOf('category-option') > -1) {
      toggleCategoryOption(theTarget);
    }
  }
  if(event.which === 40 || event.which === 38){
	var diff = (event.which === 40) ? 1 : -1;  		
  	if(theTarget.classList[0].indexOf('category-option') > -1) {
  	  var siblings = Array.prototype.slice.call(theTarget.parentElement.querySelectorAll(".category-option"));
  	  var current = siblings.indexOf(theTarget);
  	  var diff = (event.which === 40) ? 1 : -1;
  	  var next = current + diff;
  	  next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
  	  var nextItem = theTarget.parentElement.querySelectorAll(".category-option")[next];
  	  nextItem.focus();
  	} else if(theTarget.classList[0].indexOf('category') > -1) {
  	  var siblings = Array.prototype.slice.call(document.querySelectorAll(".category"));
  	  var current = siblings.indexOf(theTarget);
  	  var child = theTarget.querySelector('.category-dropdown');
  	  child.classList.remove('visible');
  	  var diff = (event.which === 40) ? 1 : -1;
  	  var next = current + diff;
  	  next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
  	  var nextItem = document.querySelectorAll(".category")[next];
  	  var nextChild = nextItem.querySelector('.category-dropdown');
  	  nextChild.classList.add('visible');
  	  nextItem.focus();
  	} else if (theTarget.classList && theTarget.classList[0] === 'search-filter-container') {
  	  var child = document.querySelector('.search-filters');
  	  var grandchild = document.querySelector('.category-dropdown');
      child.classList.toggle("visible");
      grandchild.classList.toggle("visible");
	  document.querySelector(".category").focus();
  	}
  }
  if(event.which === 39){
  	 if (theTarget.classList && theTarget.classList[0] === 'category') {
  	   var child = theTarget.querySelector('.category-option');
  	   theTarget.querySelector('.category-option').focus();
  	 }
  }
}
