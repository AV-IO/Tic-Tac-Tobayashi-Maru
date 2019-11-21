<html>
    <head>
    <title></title>
    </head>
    <body>
        <p>{{.Board}}</p>
        <p>"Enter a string to send!"</p>
        <form action="/game" method="post">
            String:<input type="text" name="String">
            <input type="submit" value="String">
        </form>
    </body>
</html>