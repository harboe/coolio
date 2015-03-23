coolio.templates.register('btn', function(params) {
	this.name = ko.observable(params.name);
	this.hit = function(viewModel, e) {
		coolio.notify(new Date(), 'SupplierName');
	};
});

coolio.templates.register('test', function(params) {
	console.log(params);
	this.palle = ko.observable(params.palle); 
});