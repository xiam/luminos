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
