var transctions=[];
var stockshares = {};
     function loadusers(){
        $.ajax({
            type: "GET",//方法类型
            dataType: "json",//预期服务器返回的数据类型
            url: "/api/users" ,//url,
            success: function (result) {
                    console.log(result);//
                    //console.log(result.length);
                    $("#from").empty();
                    $("#to").empty();
                    for (i = 0; i < result.length; i++) {
                      //console.log(result[i]["fullname"]);
                      var str="<option value=";
                      str +=result[i]["id"];
                      str +=">";
                      str += result[i]["fullname"];
                      str += "</option>";                    
                      $("#from").append(str);
                      $("#to").append(str);
                      if (i == 1){
                        $("#to").val(result[i]["id"]);
                      }
                    }              
            },
            error : function() {
                    console.log("获取数据异常！");
                }
            });        
    }



    function transfer()
    {    
      var data = {
          fromid: $("#from").val(),
          toid: $("#to").val(),
          amount: $("#amount").val()
        };
      //console.log($("#from").val());
      //console.log($("#from").find("option:selected").text());
      if(data["fromid"] == data["toid"]){
         alert("不能自己给自己转账！");
         return;
      }
      val = parseInt(data["amount"])
      if(!val>0){
         alert("Token数量必须是数字！");
         return;
      }
      var name = $("#from").find("option:selected").text();
      if( val > stockshares[name]){
         alert("转出的Token数量不能多于持有数！");
        return;
      }
      
      $.ajax({
          type: "POST",
          url: "/api/transfer",
          //traditional :true,
          //contentType: "application/json",
          data: JSON.stringify(data),        
          //data: data,
          dataType: "text",
          async: true,
          success: function(msg) {
            console.log(msg)
            obj = JSON.parse(msg);
            console.log(obj)
            transctions.push(obj['txid'])
            console.log(transctions)
            var str = ""
            var prex = "https://ropsten.etherscan.io/tx/"
            for(i=0;i < transctions.length;i++){
               var temp="<a href =\"";
               temp += prex;
               temp += transctions[i];
               temp += "\">";
               temp += transctions[i];
               temp += "</a><br/>";
               str +=temp
            }
            $("#trans").html(str);
            //if(msgObject.Status === true) {
              // window.location.href = msgObject.Url;
           // }
          },
          error : function(msg) {
             console.log(msg)
             console.log("异常！")
          }            
        });
    }
   

    function loadstocks(){
      var names = [];
      var values = [];
      $.ajax({
            //几个参数需要注意一下
            type: "GET",//方法类型
            dataType: "json",//预期服务器返回的数据类型
            url: "/api/stockshares" ,//url,
            success: function (result) {
                    //console.log(result);//
                    //console.log(result.length);
                    for (i = 0; i < result.length; i++) {
                      //console.log(result[i]["name"]);
                      names.push(result[i]["name"])
                      //console.log(result[i]["stock"]);                      
                      values.push({name:result[i]["name"],
                        value:parseInt(result[i]["stock"])});
                      stockshares[result[i]["name"]] = parseInt(result[i]["stock"]);
                      //$(".from").append(result[i]["fullname"]);                      
                    }
                    console.log(stockshares)
                    //更新pie图
                    freshpie(names,values)
            },
            error : function() {
                    console.log("获取数据异常！");
                }
            }); 
    }
    
    function freshpie(names,values){
      var dom = document.getElementById("pie_container");
      var myChart = echarts.init(dom);
      var app = {};   

      option = null;      
      option = {
          title: {
              text: '',
              subtext: '',
              //left: 'center'
              x:'left'
          },
          tooltip : {
              trigger: 'item',
              formatter: "{a} <br/>{b} : {c} ({d}%)"
          },
          toolbox: {
            show : true,
            orient: 'vertical',    
            // orient: 'horizontal',      // 布局方式，默认为水平布局，可选为：
            //                            // 'horizontal' ¦ 'vertical'
            x: 'right',                // 水平安放位置，默认为全图右对齐，可选为：
                                       // 'center' ¦ 'left' ¦ 'right'
                                       // ¦ {number}（x坐标，单位px）
            y: 'bottom',                  // 垂直安放位置，默认为全图顶端，可选为：
                                       // 'top' ¦ 'bottom' ¦ 'center'
                                       // ¦ {number}（y坐标，单位px）
            feature : {
                saveAsImage : {
                    show: true
                },               
                myTool:{  
                    show: true, //是否显示    
                    title:'刷新数据', //鼠标移动上去显示的文字    
                    icon:'image://static/img/refresh.png', //图标    
                    option:{},    
                    onclick:function(option1) {  
                        loadusers();
                        loadstocks();
                        myChart.clear();
                        myChart.setOption(option);
                    }
                },
            }           
          },
          legend: {
              //orient: 'vertical',
              //top: '10',
              bottom: 10,
              //left: 'center',
              data: names
          },
          series : [
              {
                  type: 'pie',
                  name: '股权比例',
                  radius : '75%',
                  center: ['50%', '50%'],
                  selectedMode: 'single',
                  data: values,
                  itemStyle: {
                      emphasis: {
                          shadowBlur: 10,
                          shadowOffsetX: 0,
                          shadowColor: 'rgba(0, 0, 0, 0.5)'
                      }
                  }
              }
          ]
      };
      if (option && typeof option === "object") {
          myChart.setOption(option, true);
      }
    }   