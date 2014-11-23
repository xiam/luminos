<!DOCTYPE html>
<html lang="en">

  <head>
    <link href="http://gmpg.org/xfn/11" rel="profile">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta http-equiv="content-type" content="text/html; charset=utf-8">

    <!-- Enable responsiveness on mobile devices-->
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1">

    <title>
      {{ if .IsHome }}
        {{ setting "page/head/title" }}
      {{ else }}
        {{ if .Title }}
          {{ .Title }} {{ if setting "page/head/title" }} &middot; {{ setting "page/head/title" }} {{ end }}
        {{ else }}
          {{ setting "page/head/title" }}
        {{ end }}
      {{ end }}
    </title>

    <link href="//fonts.googleapis.com/css?family=Source+Code+Pro" rel="stylesheet" type="text/css">

    <!-- CSS -->
    <link rel="stylesheet" href="{{ asset "/css/poole.css" }}">
    <link rel="stylesheet" href="{{ asset "/css/syntax.css" }}">

    <!-- Icons -->
    <link rel="apple-touch-icon-precomposed" sizes="144x144" href="{{ asset "/apple-touch-icon-precomposed.png" }}">
    <link rel="shortcut icon" href="{{ asset "/favicon.ico"}}">

    <style type="text/css">
      code {
        font-family: 'Source Code Pro';
      }
    </style>

  </head>

  <body>

    <div class="container content">
      <header class="masthead">
        <h3 class="masthead-title">
          <a href="{{ asset "/" }}" title="Home">{{ setting "page/body/title" }}</a>
          <small>{{ setting "page/brand" }}</small>
        </h3>
      </header>

      <main>

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

    {{ if .IsHome }}

      <div class="hero-unit">

        <h1>Luminos</h1>
        <p>
					A tiny server for markdown documents
        </p>

        <p class="pull-right">
          <a href="http://luminos.menteslibres.org" target="_blank" class="btn btn-large btn-primary">
            Homepage
          </a>
        </p>

      </div>

      <div class="container-fluid">
        <div class="row">
          <div class="span11">
            {{ .ContentHeader }}

            {{ .Content }}

            {{ .ContentFooter }}
          </div>
        </div>
      </div>

    {{ else }}

      {{ if .BreadCrumb }}
        <ul class="breadcrumb menu">
          {{ range .BreadCrumb }}
            <li><a href="{{ asset .link }}">{{ .text }}</a> <span class="divider">/</span></li>
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
                        <a href="{{ asset .link }}">{{ .text }}</a>
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
                      <a href="{{ asset .link }}">{{ .text }}</a>
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

    {{ end }}


      </main>

      <footer class="footer">
        <small>
          &copy; 2014. {{ setting "page/brand" }}
        </small>
      </footer>
    </div>

    {{ if setting "page/body/scripts/footer" }}
      <script type="text/javascript">
        {{ setting "page/body/scripts/footer" | jstext }}
      </script>
    {{ end }}

  </body>
</html>


