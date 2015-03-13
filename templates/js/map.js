ko.observableArray.fn.map = function(data, Constructor) {
	var mapped = ko.utils.arrayMap(data, function (item) {
		return new Constructor(item);
	});

	this(mapped);
	return this;
};