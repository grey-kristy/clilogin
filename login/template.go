package login

func getTemplate() string {
	return `
<!doctype html>
<html>
	
<head>
	  <title>Google SignIn</title>
	  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css">
	  <style>
		body {
		  padding-top: 70px;
		}
	  </style>
</head>
	
<body>
	  <div class="container">
		<div class="jumbotron">
		  <h1 class="text-success  text-center"><span class="fa fa-user"></span>Login is successful</h1>
		  <div class="row">
			<div class="col-sm-6">
			  <div class="well">
				<p>
				  <strong>Id</strong>: {{.UserID}}<br>
				  <strong>Email</strong>: {{.Email}}<br>
				  <strong>Name</strong>: {{.Name}}
				</p>
			  </div>
			</div>
		  </div>
		</div>
	  </div>
</body>
	
</html>	
`
}
