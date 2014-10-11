

function construct_table(data) {
	var oldTable = document.getElementById("maindatatable"),
	    newTable = oldTable.cloneNode();

	var thead = document.createElement('thead');

	var thead_tr = document.createElement('tr');

	var th = document.createElement('th');
	th.appendChild(document.createTextNode("Priority"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("Description"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("Deadline"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("Assigned Task/Deadline"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode("Associated"));
	thead_tr.appendChild(th);

	th = document.createElement('th');
	th.appendChild(document.createTextNode(""));
	thead_tr.appendChild(th);
	
	thead.appendChild(thead_tr);
	newTable.appendChild(thead);

	var tbody = document.createElement('tbody');

	for(var i = 0; i < data.length; i++) {
		var tr = document.createElement('tr');
		tr.setAttribute('tid', data[i]["Id"]);

		var td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["Priority"]));
		tr.appendChild(td);

		var description = data[i]["Description"];
		if (description.length > 120) {
			description = description.substring(0, 120) + " ...";
		}
		td = document.createElement('td');
		td.appendChild(document.createTextNode(description));
		tr.appendChild(td);

		td = document.createElement('td');
		var attr = document.createAttribute('style');
		attr.nodeValue = 'white-space: nowrap';
		td.setAttributeNode(attr);
		td.innerHTML = '<span class="tag tag-red"></span>&nbsp;' + data[i]["Deadline"]
		tr.appendChild(td);

		td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["AssignedTo"]));
		tr.appendChild(td);

		td = document.createElement('td');
		td.appendChild(document.createTextNode(data[i]["AssociatedPerson"]));
		tr.appendChild(td);


		td = document.createElement('td');
		var attr = document.createAttribute('style');
		attr.nodeValue = 'white-space: nowrap';
		td.setAttributeNode(attr);
		td.innerHTML = '<i class="fa fa-pencil-square-o fa-lg"></i>&nbsp;&nbsp;<i class="fa fa-remove fa-lg"></i>'
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

var firstTableConstruction = true;
var delete_happend = false;

function delete_first_searchbox()
{
		if (!firstTableConstruction && !delete_happend) {
			console.log("remove tabl");
			$('#maindatatable_wrapper').find('div').first().remove();
			delete_happend = true;
		}
		firstTableConstruction = false;
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
		// wired workaround. We geneate a new table in construct_table()
		// and later we add active dataTable. But what happends here
		// is that the search bar is added again and again. Dont know
		// how to solve this in a clean manner
		setTimeout(delete_first_searchbox(), 1)
	});

}

var usersData;

function update_user_database() {
  // first we fetch the user database
	$.ajax({
		url: "api/users"
	}).then(function(data) {
		usersData = data
	});

}

jQuery(document).ready(function($) {

	update_user_database();

	update_items_table();

	$('#form-deadline').datepicker({
		format: "dd-mm-yyyy"
	});

	$('.selectpicker').selectpicker({
    size: 20
  });

});


$("#myFormSubmit").click(function(e){
	e.preventDefault();
	var robj =  {};
	robj["Description"] = $('#form-description').val();
	if (robj["Description"] == "") {
		alert("Empty Description String");
		return;
	}
	robj["Deadline"] = $('#form-deadline').val();
	robj["AssignedTo"] = $('#form-assigned').val();
	robj["AssociatedPersons"] = $('#form-associated').val();
	robj["Tags"] = $('#form-tags').val();
	robj["Priority"] = $('#form-priority').val();
	robj["Information"] = $('#form-information').val();

	console.log(robj["Tags"]);
	var xobj = { Command: "add",
								Data: { Description: robj["Description"],
								Deadline: robj["Deadline"],
								Tags: robj["Tags"],
								AssignedTo: robj["AssignedTo"],
								Priority: robj["Priority"],
								Information: robj["Information"],
								AssociatedPersons: robj["AssociatedPersons"]}}
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

