<!DOCTYPE html>
<html lang="en">

  <head>

    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />

    {{ if .PageTitle }}
      <title>
        {{ .PageTitle }} {{ if setting "page/head/short_title" }} // {{ setting "page/head/short_title" }} {{ end }}</title>
    {{ else }}
      <title>{{ setting "page/head/title" }}</title>
    {{ end }}

    <!-- Le HTML5 shim, for IE6-8 support of HTML elements -->
    <!--[if lt IE 9]>
    <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->

    <script type="text/javascript" src="http://get.jsfoo.org/jquery.js"></script>
    <script type="text/javascript" src="http://get.jsfoo.org/jquery.foo.js"></script>

    <link rel="stylesheet" href="http://static.hckr.org/normalize/normalize.css" />

    <link rel="stylesheet" href="http://static.hckr.org/bootstrap/css/bootstrap.css" />
    <link rel="stylesheet" href="http://static.hckr.org/bootstrap/css/bootstrap-responsive.css" />

    <link rel="stylesheet" href="http://static.hckr.org/google-code-prettify/prettify.css" />
    <script type="text/javascript" src="http://static.hckr.org/google-code-prettify/prettify.js"></script>

    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <script type="text/javascript">
      $(document.body).ready(
        function() {
          var links = $('ul.menu li').removeClass('active');

          for (var i = 0; i < links.length; i++) {
            var a = $(links[i]).find('a');
            if (a.attr('href') == document.location.pathname) {
              $(links[i]).addClass('active');
            };
          };
        }
      );
    </script>

    <style type="text/css">
      .navbar .brand {
        margin-left: 0px;
      }
      body {
        padding-top: 50px;
      }
    </style>

  </head>

  <body>

    <script type="text/javascript">
      $(document).ready(prettyPrint);
    </script>

    <div class="container" id="container">
    <div class="navbar navbar-fixed-top">
      <div class="navbar-inner">
        <div class="container">

          <a class="brand" href="{{ url "/" }}">{{ setting "page/brand" }}</a>

          <div class="nav-collapse">
            {{ if settings "page/body/menu" }}
              <ul id="nav" class="nav menu">
                {{ range settings "page/body/menu" }}
                  <li>{{ link .url .text }}</li>
                {{ end }}
              </ul>
            {{ end }}
            {{ if settings "page/body/menu_pull" }}
              <ul id="nav" class="nav pull-right menu">
                {{ range settings "page/body/menu_pull" }}
                  <li>{{ link .url .text }}</li>
                {{ end }}
              </ul>
            {{ end }}
          </div>

        </div>
      </div>
    </div>

    <!--
    <div class="hero-unit">
      <h1>Heading</h1>
      <p>Tagline</p>
      <p>
        <a class="btn btn-primary btn-large">
          Learn more
        </a>
      </p>
    </div>

    <div class="page-header">
      <h1>Example page header</h1>
    </div>
    -->

    {{ if .BreadCrumb }}
      <ul class="breadcrumb menu">
        {{ range .BreadCrumb }}
          <li><a href="{{ url .link }}">{{ .text }}</a> <span class="divider">/</span></li>
        {{ end }}
      </ul>
    {{ end }}

    <div class="container-fluid">

      <div class="row">
        {{ if .SideMenu }}
          {{ if .Content }}
            <div class="span3">
                <ul class="nav nav-list menu">
                  {{ range .SideMenu }}
                    <li>
                      <a href="{{ url .link }}">{{ .text }}</a>
                    </li>
                  {{ end }}
                </ul>
            </div>
            <div class="span8">
              {{ .ContentHeader }}

              {{ .Content }}

              {{ .ContentFooter }}
            </div>
          {{ else }}
            <div class="span11">
              {{ if .CurrentPage }}
                <h1>{{ .CurrentPage.text }}</h1>
              {{ end }}
              <ul class="nav nav-list menu">
                {{ range .SideMenu }}
                  <li>
                    <a href="{{ url .link }}">{{ .text }}</a>
                  </li>
                {{ end }}
              </ul>
            </div>
          {{ end }}
        {{ else }}
          <div class="span11">
            {{ .ContentHeader }}

            {{ .Content }}

            {{ .ContentFooter }}
          </div>
        {{ end }}
      </div>

    </div>

    {{ if setting "page/body/scripts/footer" }}
      <script type="text/javascript">
        {{ setting "page/body/scripts/footer" | jstext }}
      </script>
    {{ end }}

  </body>
</html>
