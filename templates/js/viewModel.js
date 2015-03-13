var ParameterViewModel = function(data) {
	this.id = ko.observable(data.id);
	this.name = ko.observable(data.id);
	this.desc = ko.observable(data.desc);
	this.type = ko.observable('text');
	this.help = ko.observable();

	// TODO: sharepoint injection
};

var GroupViewModel = function(data) {
	this.id = ko.observable(data.id || data.name);
	this.name = ko.observable(data.name);
	this.type = ko.observable(data.type || 'group');
	this.paramType = function(viewModel, e) {
		return ko.unwrap(viewModel.type);
	};
	this.groupType = function(viewModel) {
		return ko.unwrap(viewModel.type);
	}
	this.params = ko.observableArray().map(data.params, ParameterViewModel);
	this.groups = ko.observableArray().map(data.groups, GroupViewModel);
};