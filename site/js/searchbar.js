// Searchbar JavaScript code for sfsbook

// toggleCategoryOption adds or removes a keyword from the search
function toggleCategoryOption() {
	var el = document.getElementById("query_field");
        var searchText = el.value
	var optionValue = this.children[0].value;
	var re = new RegExp("\\b(" + optionValue + ")\\b", "gi");
        var stringToReplace = optionValue.concat(", ");
	if(searchText.match(stringToReplace)) {
			var stringToReplace = optionValue.concat(", ");
			searchText = searchText.replace(stringToReplace, "");
			el.value = searchText;
		} else if (searchText.match(re)) {
			searchText = searchText.replace(optionValue, "");
			el.value = searchText;
		} else {
		if(searchText.length > 0){
			el.value = searchText.concat(", ",this.children[0].value)
		} else {
			el.value = optionValue;
		}
	}
}

// addEventListener attaches toggleCategoryOption to the checkboxes
function addEventListener(el, eventName, handler) {
	eventName.preventDefault ? eventName.preventDefault() : (eventName.returnValue = false);
  	if (el.addEventListener) {
    		el.addEventListener(eventName, handler);
  	} else {
    		el.attachEvent('on' + eventName, function(){
      			handler.call(el);
    		});
  	}
}

// attachToggles attaches event listener to each category option
function attachToggles() {
	var categories = document.getElementsByClassName("category-option");
	for(var i = 0; i < categories.length; i++){
		el = categories[i];
		addEventListener(el, 'click', toggleCategoryOption);
	}
}

// ready calls addEventListener once the page has loaded
function ready(fn) {
  	if (document.readyState != 'loading'){
    		fn();
  	} else if (document.addEventListener) {
    		document.addEventListener('DOMContentLoaded', fn);
  	} else {
    		document.attachEvent('onreadystatechange', function() {
      		if (document.readyState != 'loading')
        		fn();
    		});
  	}
}

//calls ready with attachToggles
ready(attachToggles);
