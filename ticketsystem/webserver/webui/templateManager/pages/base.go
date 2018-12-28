package pages

/*
	Base Html template.
*/
var Base = `
	{{ define "Base" }}
	<!DOCTYPE html>
	<html>

	<head>
		<title>
			{{ block "Title" .}} {{ end }}
		</title>
		<link rel="stylesheet" href="/files/style/main">
		{{ block "StylesAndScripts" .}} {{ end }}
	</head>
	
	<body>
		{{ template "Content" .}}
	</body>
	
	</html>
	{{ end }}`
