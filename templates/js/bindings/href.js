ko.bindingHandlers.href = {
	update: function(elm, valueAccessor) {
		var val = ko.unwrap(valueAccessor());
		$(elm).attr('href', val);
	}
};

ko.bindingHandlers.hash = {
	update: function(elm, valueAccessor) {
		var val = ko.unwrap(valueAccessor());
		$(elm).attr('href', '#' + val);
	}
};