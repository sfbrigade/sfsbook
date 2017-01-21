// Searchbar JavaScript code for sfsbook

// toggleCategoryOption adds or removes a keyword from the search
function toggleCategoryOption() {
  var el = document.getElementById('query_field');
  var searchText = el.value;
  var optionValue = this.textContent;
  console.log(optionValue, this);
  var re = new RegExp('\\b(' + optionValue + ')\\b', 'gi');
  var stringToReplace = optionValue.concat(', ');

  if(!searchText.match(stringToReplace) && !searchText.match(re)) {
    /*var hbox = document.querySelectorAll('.hbox')[0];
    var btn = document.createElement('BUTTON');
    btn.value = optionValue;
    var t = document.createTextNode(optionValue);
    btn.appendChild(t);
    var updateInsert = function(buttonValue) {
      var searchfield = document.getElementById('query_field');
      if (searchfield.value.match(optionValue.concat(', '))) {
        searchfield.value = searchfield.value.replace(buttonValue.concat(', '), '');
      } else {
        searchfield.value = searchfield.value.replace(buttonValue, '');
      }
    };
    btn.onclick = function(event) {
      event.preventDefault();
      event.stopPropagation();
      updateInsert(this.value);
      this.remove();
    };
    hbox.appendChild(btn);*/
    if (searchText.length > 0) {
      el.value = searchText.concat(', ', optionValue);
    } else {
      el.value = optionValue;
    }
  }
}

// toggleActiveClass adds or removes active from user-menu class element
function toggleActiveClass() {
  if (this.classList.length > 1) {
    this.classList.toggle('user-menu-active');
  } else {
    var classes = this.classList.value.split(' ');
    classes.push('user-menu-active');
    this.classList.value = classes.join(' ');
  }
}

/* toggleHiddenNav adds or removes hidden from nav on mobile
function toggleHiddenNav() {
  var navbar = document.querySelectorAll('.nav')[0];
  if (navbar.classList.length > 1) {
    navbar.classList.toggle('nav-hidden');
  } else {
    var classes = navbar.classList.value.split(' ');
    classes.push('nav-hidden');
    navbar.classList.value = classes.join(' ');
  }
}
*/
// toggleHiddenCategory adds or removes hidden from category on mobile
function toggleHiddenCategory() {
  console.log('in hidden', this);
  var category = this.parentElement.children[1];
  if (category.classList.length > 1) {
    category.classList.toggle('category-hidden');
  } else {
    var classes = category.classList.value.split(' ');
    classes.push('category-hidden');
    category.classList.value = classes.join(' ');
  }
}


// addEventListener attaches toggleCategoryOption to the checkboxes
function addEventListener(el, eventName, handler) {
  eventName.preventDefault ? eventName.preventDefault() :
  (eventName.returnValue = false);
  eventName.stopPropagation ? eventName.stopPropagation() :
  (eventName.cancelBubble = true);
  if (el.addEventListener) {
    el.addEventListener(eventName, handler);
  } else {
    el.attachEvent('on' + eventName, function() {
      handler.call(el);
    });
  }
}

// attachToggles attaches event listener to each category option
function attachToggles() {
/*addEventListener(document.querySelectorAll('.logo')[0], 'click',
  toggleHiddenNav);*/
  var clickPairs = [['.expandtab', toggleHiddenCategory],
                     ['.user-menu', toggleActiveClass],
                     ['.category-option', toggleCategoryOption]];
  for (var i = 0; i < clickPairs.length; i++) {
    var trigger = document.querySelectorAll(clickPairs[i][0]);
    for (var j = 0; j < trigger.length; j++) {
      addEventListener(trigger[j], 'click', clickPairs[i][1]);
    }
  }
}

// ready calls addEventListener once the page has loaded
function ready(fn) {
  var txt;
  var r = confirm("Safety Alert: Computer use can be monitored and is impossible to completely clear. If you are afraid your internet usage might be monitored, call the SFWAR Hotline at 415-647-7273. Would you like to exit this site?");
  if (r == true) {
    window.location.replace('https://weather.com/');
  } else {
    txt = "You pressed Cancel!";
    }

  if (document.readyState != 'loading') {
    fn();
  } else if (document.addEventListener) {
    document.addEventListener('DOMContentLoaded', fn);
  } else {
    document.attachEvent('onreadystatechange', function() {
      if (document.readyState != 'loading') {
        fn();
      }
    });
  }
}

// calls ready function with attachToggles
ready(attachToggles);
