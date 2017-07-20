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
 * Shortcut for document.getElementById.
 * @param {string} identifier The selector id.
 * @return {object} The selected node.
 */
function gebi(identifier) {
  return document.getElementById(identifier);
}

/**
 * Shortcut for document.querySelector.
 * @param {string} identifier The selector class.
 * @param {object}  nodeList The container list.
 * @return {object} The selected node.
 */
function qs(identifier, nodeList) {
  nodeList = (nodeList || document);
  return nodeList.querySelector(identifier);
}

/**
 * Shortcut for document.querySelectorAll.
 * @param {string} identifier The selector class.
 * @param {object} nodeList The container list.
 * @return {object} The selected node.
 */
function qsa(identifier, nodeList) {
  nodeList = (nodeList || document);
  return nodeList.querySelectorAll(identifier);
}

/**
 * Adds a keyword to the search input.
 * @param {object} node The selected keyword node.
 */
function toggleCategoryOption(node) {
  var el = gebi('query_field');
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
        var siblings = Array.prototype.slice.call(qsa('.user-menu-item'));
        var current = siblings.indexOf(theTarget);
        var diff = (event.which === 40) ? 1 : -1;
        var next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        var nextItem = qsa('.user-menu-item')[next];
        nextItem.focus();
      } else if (theTarget.classList && theTarget.classList[0] === 'user-menu') {
        var child = qs('.user-menu-dropdown');
        child.classList.toggle('visible');
        qs('.user-menu-item').focus();
      }
      break;

  case keys.enter:
    if (theTarget.classList && theTarget.classList[0] === 'user-menu-item') {
      location.href=theTarget.parentElement.href;
    }
    break;

  case keys.tab:
    if (theTarget.classList && theTarget.classList[0] === 'user-menu-item') {
      var parent = qs('.user-menu');
      var child = qs('.user-menu-dropdown');
      if (child.classList[1] && (child.classList[1] === 'visible')) {
        child.classList.toggle('visible');
      }
    }
    break;

  case keys.esc:
    if (theTarget.classList && theTarget.classList[0] === 'user-menu-item') {
      var parent = qs('.user-menu');
      var child = qs('.user-menu-dropdown');
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
        qs('.search-filters').classList.remove('visible');
        var secondlevel = qsa('.category-dropdown');
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
        siblings = Array.prototype.slice.call(qsa('.category-option', theTarget.parentElement));
        current = siblings.indexOf(theTarget);
        diff = (event.which === 40) ? 1 : -1;
        next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        nextItem = qsa('.category-option', theTarget.parentElement)[next];
        nextItem.focus();
      } else if (theTarget.classList[0] === 'category') {
        siblings = Array.prototype.slice.call(qsa('.category'));
        current = siblings.indexOf(theTarget);
        child = qs('.category-dropdown', theTarget);
        child.classList.remove('visible');
        diff = (event.which === 40) ? 1 : -1;
        next = current + diff;
        next = (next < 0) ? siblings.length - 1 : (next === siblings.length) ? 0 : next;
        nextItem = qsa('.category')[next];
        nextChild = qs('.category-dropdown', nextItem);
        nextChild.classList.add('visible');
        nextItem.focus();
      } else if (theTarget.classList && theTarget.classList[0] === 'search-filter-container') {
        child = qs('.search-filters');
        grandchild = qs('.category-dropdown');
        child.classList.toggle('visible');
        grandchild.classList.toggle('visible');
        qs('.category').focus();
      }
      break;

    case keys.right:
      if (theTarget.classList && theTarget.classList[0] === 'category') {
        child = qs('.category-option', theTarget);
        qs('.category-option', theTarget).focus();
      }
      break;

  }
}
