ko.bindingHandlers.placeholder = {
	update: function(elm, valueAccessor) {
		var val = ko.unwrap(valueAccessor());
		$(elm).attr('placeholder', val);
	}
};