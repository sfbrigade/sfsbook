// Searchbar JavaScript code for sfsbook

// toggleCategoryOption adds or removes a keyword from the search
function toggleCategoryOption(node) {
  console.log(arguments);
  var el = document.getElementById('query_field');
  var searchText = el.value;
  var input = (node.type === 'click') ? node.target : node;
  console.log('input', input);
  console.log(input.textContent, 'intxtcn');
  var optionValue = input.textContent;
  console.log(optionValue, 'optionValue');
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

function attachToggles() {
  var clickPairs = [['.category-option', toggleCategoryOption]];
  var keydownPairs = [['.search-filter-container', searchNav],
                      ['.search-filters', searchNav],
                      ['.category-option', searchNav],
                      ['.user-menu', userNav],
                      ['.user-menu-item', userNav]];
  for (var i = 0; i < clickPairs.length; i++) {
    var trigger = document.querySelectorAll(clickPairs[i][0]);
    for (var j = 0; j < trigger.length; j++) {
      addEventListener(trigger[j], 'click', clickPairs[i][1]);
    }
  }
  for (var k = 0; k < keydownPairs.length; k++) {
    var trigger = document.querySelectorAll(keydownPairs[k][0]);
    for (var l = 0; l < trigger.length; l++) {
      addEventListener(trigger[l], 'keydown', keydownPairs[k][1]);
    }
  }
}

function ready(fn) {
  if(document.referrer.indexOf('localhost') === -1) {
  var response = confirm("Safety Alert: Computer use can be monitored and is impossible to completely clear. If you are afraid your internet usage might be monitored, exit this site and call the SFWAR Hotline at 415-647-7273. Would you like to exit this site?");
  if (response === true) {
    window.location.replace('https://weather.com/');
  } else {
    txt = "You pressed Cancel!";
    }
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
