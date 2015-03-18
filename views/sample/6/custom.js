coolio.templates.register('btn', function(params) {
  this.name = ko.observable(params.name);
  this.hit = function(viewModel, e) {
	alert('test...');
  };
});