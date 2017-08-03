/**
 * Attaches event handlers for IE8, stops event propagation and prevents default event behavior.
 * @param {object} el The element to which the event will be attached.
 * @param {object} eventName The event to be attached.
 * @param {function} handler The handler to be invoked when the event occurs.
 */
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

/** Attaches event listeners to toggle dropdown menu elements */
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

function attachSearchListeners() {
  var searchClickTrigger = document.getElementById('search-submit-button');
  var searchEnterTrigger = document.getElementById('query_field');
  addEventListener(searchClickTrigger, 'click', fetchSearchResults);
  addEventListener(searchEnterTrigger, 'keypress', fetchSearchResults);
}

/**
 * Listens for document ready (i.e. readyState not loading), and invokes callback.
 * @param {function} fn The function to be invoked once ready detected.
 */
function ready(fn) {
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

/** Calls attachToggles once ready is detected */
ready(attachToggles);
// will uncomment this out when we refactor to single page search
// ready(attachSearchListeners);
