ko.bindingHandlers.id = {
	update: function(elm, valueAccessor) {
		var val = ko.unwrap(valueAccessor());
		$(elm).attr('id', val);		
	}
};