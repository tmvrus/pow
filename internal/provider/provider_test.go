package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider_HappyPath(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write(response)
	}))

	res, err := NewProvider(server.URL).Get(context.Background())
	require.NoError(t, err)
	require.Equal(t, "I couldn't find any library that can even do that", res)

}

var response = []byte(`<!DOCTYPE HTML>
<html>
<head>
 <title>Excuses For Lazy Coders</title>
 <meta name="keywords" content="Programming, Programmer, Devloping, Developer, Coding, Coders, Excuses, Reasons, Lies, Fibs, Blame, Justifications" />
 <meta name="description" content="Excuses For Lazy Coders" />
 <link rel="canonical" href="http://programmingexcuses.com/" />
 <style type="text/css">* {margin: 0;} html, body {height: 100%;} .wrapper {min-height: 100%; height: auto !important; height: 100%; margin: 0 auto -8em;} .footer, .push {height: 8em;}</style>
</head>
<body>
 <div class="wrapper">
  <center style="color: #333; padding-top: 200px; font-family: Courier; font-size: 24px; font-weight: bold;"><a href="/" rel="nofollow" style="text-decoration: none; color: #333;">I couldn't find any library that can even do that</a></center>
  <div class="push"></div>
 </div>
 <div class="footer">
  <center style="color: #333; font-family: Courier; font-size: 11px;">
   <script type="text/javascript"><!--
    google_ad_client = "ca-pub-4336860580083128";
    google_ad_slot = "1671975908";
    google_ad_width = 728;
    google_ad_height = 90;
    //-->
   </script>
   <script type="text/javascript" src="http://pagead2.googlesyndication.com/pagead/show_ads.js"></script>
   <br /><br />&copy; Copyright 2012 - 2025 programmingexcuses.com - All Rights Reserved
  </center>
 </div>
 <script type="text/javascript">
  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-33167244-1']);
  _gaq.push(['_setDomainName', 'programmingexcuses.com']);
  _gaq.push(['_setAllowLinker', true]);
  _gaq.push(['_trackPageview']);
  (function() {
   var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
   ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
   var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();
 </script>
</body>
</html>`)
