

function construct_table(data) {
	var oldTable = document.getElementById("maindatatable"),
	    newTable = oldTable.cloneNode();

	var thead = document.createElement('thead');
	var thead_tr = document.createElement('tr');

	var th = document.createElement('th');
	th.appendChild(document.createTextNode("foo"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("foo"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("foo"));
	thead_tr.appendChild(th);
	
	thead.appendChild(thead_tr);
	newTable.appendChild(thead);

	var tbody = document.createElement('tbody');

	for(var i = 0; i < data.length; i++) {
		var tr = document.createElement('tr');

		var td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["Id"]));
		tr.appendChild(td);

		td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["Description"]));
		tr.appendChild(td);

		td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["Description"]));
		tr.appendChild(td);

		tbody.appendChild(tr);
	}
	newTable.appendChild(tbody);

	oldTable.parentNode.replaceChild(newTable, oldTable);

	/*
	var oldTable = document.getElementById('example'),
	    newTable = oldTable.cloneNode();

	for(var i = 0; i < json_example.length; i++){
			var tr = document.createElement('tr');
			for(var j = 0; j < json_example[i].length; j++){
					var td = document.createElement('td');
					td.appendChild(document.createTextNode(json_example[i][j]));
					tr.appendChild(td);
			}
			newTable.appendChild(tr);
	}

	oldTable.parentNode.replaceChild(newTable, oldTable);
	*/

}

function update_items_table() {
	$.ajax({
		url: "api/items"
	}).then(function(data) {
		construct_table(data)
		//$('#tabledisplay').empty();
	  //$('#tabledisplay').append("<table id=\"maindatatable\" class=\"table table-striped table-hover\" cellspacing=\"0\" width=\"100%\"></table>");
		$('#maindatatable').dataTable( {
			"paging":   false,
			"info":     false,
			"order": [[ 0, "desc" ]]
		} );

	});

}


jQuery(document).ready(function($) {

	$(".clickableRow").click(function() {
		alert($(this).attr("href"));
	});

	$.ajax({
		url: "api/users"
	}).then(function(data) {
		$('.greeting-id').append(data.users);
	});

	update_items_table()

	$('#example1').datepicker({
		format: "dd/mm/yyyy"
	});

	$('#example2').datepicker({
		format: "dd/mm/yyyy"
	});

});


$("#myFormSubmit").click(function(e){
	e.preventDefault();
	var robj =  {};
	robj["Description"] = $('#description').val();
	if (robj["Description"] == "") {
		alert("Empty Description String");
		return;
	}
	var data = $('#description').val(); console.log(data);
	var xobj = { Command: "add", Data: { Description: robj["Description"] }}
	$.post('/api/items',
			JSON.stringify(xobj),
			function(data, status, xhr) {
				// we update the table not immediatly because
				// we first remove the modal dialog to get a
				// smoother experience[TM]
				setTimeout(update_items_table(), 1)
			});
	// after submit, disable the modal
	$('#myModal').modal('toggle');
});

