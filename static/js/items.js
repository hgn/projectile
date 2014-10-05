
jQuery(document).ready(function($) {

	$(".clickableRow").click(function() {
		alert($(this).attr("href"));
	});

	$.ajax({
		url: "api/users"
	}).then(function(data) {
		$('.greeting-id').append(data.users);
	});

	setTimeout(function(){
		$('#myModal').modal('toggle');},40000);

	$('#example1').datepicker({
		format: "dd/mm/yyyy"
	});

});


$("#myFormSubmit").click(function(e){
	e.preventDefault();
	var xobj = { Command: "add", Data: { Description: "foo" }}
	$.post('/api/items',
			JSON.stringify(xobj),
			function(data, status, xhr) {
				console.log(data.Status);
			});
	// after submit, disable the modal
	$('#myModal').modal('toggle');
});

