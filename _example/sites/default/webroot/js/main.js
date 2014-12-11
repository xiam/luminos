(function() {

  var getParent = function(start, fn) {
    while (start) {
      if (fn(start)) {
        return start;
      };
      start = start.parentNode;
    };
    return null;
  };

  var normalizeLink = function(s) {
    if (s != null) {
      s = s.replace(/^\/+/g, '');
      s = s.replace(/\/+$/g, '');
      return s;
    };
    return "";
  };

  // setActiveLinks loops over relevant links and sets the active class if
  // required.
  var setActiveLinks = function() {
    var anchors = document.body.querySelectorAll('a');

    for (var i = 0; i < anchors.length; i++) {
      var anchor = anchors[i];
      if (anchor) {
        if (normalizeLink(anchor.getAttribute('href')) == normalizeLink(document.location.pathname)) {
          // Look for a suitable parent.
          var li = getParent(anchor, function(el) {
            if (el.tagName == 'LI') {
              return true;
            };
            return false;
          });
          if (li) {
            li.className += ' active';
          } else {
            anchor.className += ' active';
          };
        };
      };
    };
  };

  // load sets the task to be run at page loaded event.
  var load = function() {
    setActiveLinks();
  };

  window.onload = load;
})();


/*
$(document.body).ready(
  function() {
    // Code (marking code blocks for prettyPrint)
    var code = $('code');

    for (var i = 0; i < code.length; i++) {
      var el = $(code[i])
      var className = el.attr('class');
      if (className) {
        el.addClass('language-'+className);
      }
    };

    // An exception, LaTeX blocks.
    var code = $('code.latex');

    for (var i = 0; i < code.length; i++) {
      var el = $(code[i])
      var img = $('<img>', { 'src': '//menteslibres.net/api/latex/png?t='+encodeURIComponent(el.html()) });
      img.insertBefore(el);
      el.hide();
    };

    // Starting prettyPrint.
    hljs.initHighlightingOnLoad();

    // Tables without class

    $('table').each(
      function(i, el) {
        if (!$(el).attr('class')) {
          $(el).addClass('table');
        };
      }
    );

    // Navigation
    var links = $('ul.menu li').removeClass('active');

    for (var i = 0; i < links.length; i++) {
      var a = $(links[i]).find('a');
      if (a.attr('href') == document.location.pathname) {
        $(links[i]).addClass('active');
      };
    };

  }
);
*/
