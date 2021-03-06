package views


var login_template = `<html lang="en">
<head>
    <meta name="google-signin-scope" content="{{range $idx, $scope := .Scopes}}{{$scope}} {{end}}">
    <meta name="google-signin-client_id" content="{{ .ClientID }}">
    <title>{{.ApplicationName}}</title>
    <script src="https://apis.google.com/js/platform.js" async defer></script>
</head>
<body>
<div class="g-signin2" data-onsuccess="onSignIn" data-theme="dark"></div>
<script>
    function onSignIn(googleUser) {
        // Useful data for your client-side scripts:
        var profile = googleUser.getBasicProfile();

        // The ID token you need to pass to your backend:
        var id_token = googleUser.getAuthResponse().id_token;
        //console.log("ID Token: " + id_token);

        // Send it to the backend
        var xhr = new XMLHttpRequest();
        xhr.open('POST', '{{.CallbackURL}}');
        xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
        xhr.onload = function () {
            console.log('Signed in as: ' + xhr.responseText);
        };
        xhr.send('idtoken=' + id_token);

    };
</script>
</body>
</html>`
