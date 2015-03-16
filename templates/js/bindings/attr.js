ko.bindingHandlers.href = {
	update: function(elm, valueAccessor) {
		$(elm).attr('href', ko.unwrap(valueAccessor()));
	}
};

ko.bindingHandlers.hash = {
	update: function(elm, valueAccessor) {
		$(elm).attr('href', '#' + ko.unwrap(valueAccessor()));
	}
};

ko.bindingHandlers.id = {
	update: function(elm, valueAccessor) {
		$(elm).attr('id', ko.unwrap(valueAccessor()));
	}
};

ko.bindingHandlers.placeholder = {
	update: function(elm, valueAccessor) {
		$(elm).attr('placeholder', ko.unwrap(valueAccessor()));
	}
};

ko.bindingHandlers.type = {
	update: function(elm, valueAccessor) {
		$(elm).attr('type', ko.unwrap(valueAccessor()));
	}
};