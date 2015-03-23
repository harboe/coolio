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

coolio.templates.register('calc', function(params) {
  this.name = ko.observable().subscribeTo('SupplierName');
  this.address = ko.observable().subscribeTo('address');
  
  this.computed = ko.computed(function() {
    return this.name() + " - " + this.address();
  }, this);
});