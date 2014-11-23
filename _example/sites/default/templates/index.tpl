<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en-us">

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


    <!-- CSS -->
    <link rel="stylesheet" href="{{ asset "/css/poole.css" }}">
    <link rel="stylesheet" href="{{ asset "/css/syntax.css" }}">
    <link rel="stylesheet" href="{{ asset "/css/hyde.css" }}">
    <link rel="stylesheet" href="{{ asset "/css/luminos.css" }}">

    <!-- External fonts -->
    <link href="//fonts.googleapis.com/css?family=Source+Code+Pro" rel="stylesheet" type="text/css">
    <link href="//fonts.googleapis.com/css?family=Open+Sans:300,400italic,400,600,700|Abril+Fatface" rel="stylesheet" type="text/css">

    <!-- Icons -->
    <link rel="apple-touch-icon-precomposed" sizes="144x144" href="{{ asset "/apple-touch-icon-precomposed.png" }}">
    <link rel="shortcut icon" href="{{ asset "/favicon.ico"}}">

  </head>

  <body>

    <!-- sidebar -->
    <div class="sidebar">
      <div class="container sidebar-sticky">
        <div class="sidebar-about">
          <h1>
            <a href="{{ asset "/" }}">
              {{ setting "page/brand" }}
            </a>
          </h1>
          <p class="lead">{{ setting "page/body/title" }}</p>
        </div>

        <nav class="sidebar-nav">
          <a class="sidebar-nav-item active" href="{{ asset "/" }}">Home</a>

          {{ if settings "page/body/menu" }}
            {{ range settings "page/body/menu" }}
              <a class="sidebar-nav-item" href="{{ .url }}">{{ .text }}</a>
            {{ end }}
          {{ end }}

          {{ if settings "page/body/menu_pull" }}
            {{ range settings "page/body/menu_pull" }}
              <a class="sidebar-nav-item" href="{{ .url }}">{{ .text }}</a>
            {{ end }}
          {{ end }}

          {{ if .SideMenu }}
            {{ range .SideMenu }}
              <a class="sidebar-nav-item" href="{{ .url }}">{{ .text }}</a>
            {{ end }}
          {{ end }}

          <a class="sidebar-nav-item" href="https://menteslibres.net/luminos/download">Download</a>
          <a class="sidebar-nav-item" href="https://github.com/xiam/luminos">GitHub project</a>
          <span class="sidebar-nav-item">Currently v0.9</span>
        </nav>

        <p>&copy; 2014. Some rights reserved.</p>
      </div>
    </div>

    <div class="content container">
      {{ if .Content }}
        {{ .ContentHeader }}

        {{ .Content }}

        {{ .ContentFooter }}
      {{ else }}

        {{ if .CurrentPage }}
          <h1>{{ .CurrentPage.text }}</h1>
        {{ end }}

        <ul>
          {{ range .SideMenu }}
            <li>
              <a href="{{ asset .link }}">{{ .text }}</a>
            </li>
          {{ end }}
        </ul>

      {{end}}
    </div>

  {{ if setting "page/body/scripts/footer" }}
    <script type="text/javascript">
      {{ setting "page/body/scripts/footer" | jstext }}
    </script>
  {{ end }}

  </body>
</html>
