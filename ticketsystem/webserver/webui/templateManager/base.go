package templateManager

/*
	Base Html template.
*/
var base = `
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
