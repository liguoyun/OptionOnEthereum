<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">
    <title>区块链股权转让演示</title>
    <!-- Bootstrap core CSS -->
    <link href="/static/css/main.css" rel="stylesheet"/>
    <link href="/static/css/bootstrap.css" rel="stylesheet"/>
  </head>
  <body>
    <nav class="navbar navbar-inverse navbar-fixed-top">
      <div class="container">
        <div class="navbar-header">
          <!--<p>Id has the value of : {{.}} </p>-->
          <h2>区块链股权转让演示</h2>
        </div>
      </div>
    </nav>
    <div class="container">
      <div class="row">
        <div class="col-md-4">
          <!--<h2>Transfer Form</h2>-->
          <form role = "form">
            <div class = "form-group">
                <label for = "from">转出方</label>
                <select id = "from" class = "form-control">
                </select>
            </div>
            <div class = "form-group">
                <label for = "to">转入方</label>
                <select id = "to" class = "form-control">
                </select>
            </div>            
            <div class = "form-group">
                <label for = "amount">Token数量</label>
                <input type="text" class="form-control" id = "amount" 
                        placeholder = "请输入Token数量"></input>
            </div>          
           
            <div class="form-group">
                <button type = "submit" class="btn-info btn-lg"  id="transfertoken">转帐</button> 
            </div>
          </form>
        </div>
        <div class="col-md-8">
          <!--<h2>Pie Chart</h2>-->          
          <div id="pie_container" style="height: 50%"></div>
        </div>   
    </div> <!-- /container -->    
    <div class="panel panel-default">
      <div class="panel-heading">即时转账交易记录</div>      
      <div id="trans">
      </div>
    </div> 

    <script src="/static/js/jquery-3.3.1.min.js"></script>
    <script src="/static/js/bootstrap.min.js"></script>
    <script src="/static/js/echarts.min.js"></script>
    <script src="/static/js/main.js"></script>
    <script>
    
    $("#transfertoken").click(function(event){
        event.preventDefault();   //阻止网页自动提交
        transfer();
    });

    $(document).ready(function(){
        loadusers();
        loadstocks();
    });
    </script>
  </body>
</html>