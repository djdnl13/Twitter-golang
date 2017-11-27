var ip = "http://192.168.99.100";

function fetchData() {
	$.ajax({
		crossDomain: true,
		url: ip+":6768/get/20",
		success: function(result){
			$("#tweets").empty();
			console.log(result);
			for(var t in result)
			{
				$("#tweets").append("<div class='tweetWrapper'><img class='avatar' src='gfx/user1.png'> <span class='name'> &nbsp;&nbsp;" + result[t].accountId+ "<span class='time'> " + result[t].likesCount+"likes</span> <p>" + result[t].text + "</p></div>");
			}
			
		}
	});

}

$("#addButton").click(function(){
    var text = $("#addText").val();
    var accountId = $("#addAccountId").val();
 
    console.log(text);
    console.log(accountId); 
    $.ajax({crossDomain: true, url: ip+":6768/add?text="+text+"&accountId="+accountId, success: function(result){
	$("#addText").val("");
	$("#addAccountId").val("");
	console.log(result);
    	fetchData();
    }});
});

$( document ).ready(function() {
	fetchData();
});
