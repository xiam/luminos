<!DOCTYPE html>

<html lang="en">

  <head>

    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />

    <link href="//fonts.googleapis.com/css?family=PT+Serif" rel="stylesheet" type="text/css">
    <link href="//fonts.googleapis.com/css?family=PT+Sans" rel="stylesheet" type="text/css">
    <link href="//fonts.googleapis.com/css?family=Source+Code+Pro" rel="stylesheet" type="text/css">

    {{ if .IsHome }}
        <title>{{ setting "page/head/title" }}</title>
    {{ else }}
      {{ if .Title }}
        <title>
          {{ .Title }} {{ if setting "page/head/title" }} // {{ setting "page/head/title" }} {{ end }}</title>
      {{ else }}
        <title>{{ setting "page/head/title" }}</title>
      {{ end }}
    {{ end }}

		<link rel="shortcut icon" href="{{ asset "/favicon.ico" }}" />

		<script src="//ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>

    <link rel="stylesheet" href="//menteslibres.net/static/normalize/normalize.css" />

    <link rel="stylesheet" href="//menteslibres.net/static/bootstrap/css/bootstrap.css" />
    <link rel="stylesheet" href="//menteslibres.net/static/bootstrap/css/bootstrap-responsive.css" />

    <link rel="stylesheet" href="//menteslibres.net/static/highlightjs/styles/solarized_dark.css">
    <script src="//menteslibres.net/static/highlightjs/highlight.pack.js"></script>

    <link rel="stylesheet" href="{{ asset "/css/styles.css" }}" />

    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <script type="text/javascript" src="{{ asset "/js/main.js" }}"></script>

		<style type="text/css">
			body {
				font-family: 'PT Sans';
				font-size: large;
			}
			code {
				font-family: 'Source Code Pro';
			}
		</style>

  </head>

  <body>

    <div class="container" id="container">
    <div class="navbar navbar-fixed-top">
      <div class="navbar-inner">
        <div class="container">

          <a class="brand" href="{{ asset "/" }}">{{ setting "page/brand" }}</a>

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

    <hr />

    <footer>
      Powered by <a href="https://menteslibres.net/luminos" target="_blank">Luminos</a>
    </footer>

    {{ if setting "page/body/scripts/footer" }}
      <script type="text/javascript">
        {{ setting "page/body/scripts/footer" | jstext }}
      </script>
    {{ end }}

  </body>
</html>
