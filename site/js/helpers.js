var keys = {
  tab: 9,
  enter: 13,
  esc: 27,
  space: 32,
  left: 37,
  up: 38,
  right: 39,
  down: 40,
};

/**
 * Shortcut for document.querySelector.
 * @param {string} identifier The selector class.
 * @return {object} The selected node.
 */
function query(identifier) {
  return document.querySelector(identifier);
}

/**
 * Adds a keyword to the search input.
 * @param {object} node The selected keyword node.
 */
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

/** Delete the session cookie and reload the page in the unauthed state. */
function clearSessionCookie() {
  if ((event.which === 13) || (event.which === 1)) {
    document.cookie = 'session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/';
    window.location.reload();
  }
}

/**
 * Enables visible keyboard navigation of user menu.
 * @param {object} node The selected navigation node.
 */
function userNav(node) {
  var theTarget = event.target;
  switch (event.which) {
    case keys.up:
    case keys.down:
      if (theTarget.classList[0] === 'user-menu-item') {
        var siblings = Array.prototype.slice.call(document.querySelectorAll('.user-menu-item'));
        var current = siblings.indexOf(theTarget);
        var diff = (event.which === 40) ? 1 : -1;
        var next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        var nextItem = document.querySelectorAll('.user-menu-item')[next];
        nextItem.focus();
      } else if (theTarget.classList && theTarget.classList[0] === 'user-menu') {
        var child = query('.user-menu-dropdown');
        child.classList.toggle('visible');
        query('.user-menu-item').focus();
      }
      break;

  case keys.enter:
    if (theTarget.classList && theTarget.classList[0] === 'user-menu-item') {
      location.href=theTarget.parentElement.href;
    }
    break;

  case keys.esc:
    if (theTarget.classList && theTarget.classList[0] === 'user-menu-item') {
      var parent = query('.user-menu');
      var child = query('.user-menu-dropdown');
      if (child.classList[1] && (child.classList[1] === 'visible')) {
        child.classList.toggle('visible');
        parent.focus();
      }
    }
    break;

  }

}

/**
 * Enable visible keyboard navigation of search filters
 * @param {object} node The selected navigation node.
 */
function searchNav(node) {
  var theTarget = event.target,
    diff,
    siblings,
    current,
    next,
    nextItem,
    child,
    nextChild,
    grandchild;

  switch (event.which) {
    case keys.tab:
      if ((theTarget.classList[0] === 'category') || (theTarget.classList[0] === 'search-filter-container')) {
        query('.search-filters').classList.remove('visible');
        var secondlevel = document.querySelectorAll('.category-dropdown');
        for(var i = 0; i < secondlevel.length; i++) {
          secondlevel[i].classList.remove('visible');
        }
      }
      break;

    case keys.esc:
    case keys.left:
      if ((theTarget.classList[0] === 'category')||(theTarget.classList[0] === 'category-option')) {
        theTarget.parentElement.classList.toggle('visible');
          theTarget.parentElement.parentElement.focus();
      }
      break;

    case keys.enter:
      if (theTarget.classList[0] === 'category-option') {
        toggleCategoryOption(theTarget);
      }
      break;

    case keys.down:
    case keys.up:
      diff = (event.which === 40) ? 1 : -1;
      if (theTarget.classList[0] === 'category-option') {
        siblings = Array.prototype.slice.call(theTarget.parentElement.querySelectorAll('.category-option'));
        current = siblings.indexOf(theTarget);
        diff = (event.which === 40) ? 1 : -1;
        next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        nextItem = theTarget.parentElement.querySelectorAll('.category-option')[next];
        nextItem.focus();
      } else if (theTarget.classList[0] === 'category') {
        siblings = Array.prototype.slice.call(document.querySelectorAll('.category'));
        current = siblings.indexOf(theTarget);
        child = theTarget.querySelector('.category-dropdown');
        child.classList.remove('visible');
        diff = (event.which === 40) ? 1 : -1;
        next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        nextItem = document.querySelectorAll('.category')[next];
        nextChild = nextItem.querySelector('.category-dropdown');
        nextChild.classList.add('visible');
        nextItem.focus();
      } else if (theTarget.classList && theTarget.classList[0] === 'search-filter-container') {
        child = query('.search-filters');
        grandchild = query('.category-dropdown');
        child.classList.toggle('visible');
        grandchild.classList.toggle('visible');
        query('.category').focus();
      }
      break;

    case keys.right:
      if (theTarget.classList && theTarget.classList[0] === 'category') {
        child = theTarget.querySelector('.category-option');
        theTarget.querySelector('.category-option').focus();
      }
      break;

  }
}
