var ParameterViewModel = function(data) {

	var spFieldCtx = data.id ? sharepoint.getField(data.id): {}, 
	schema = spFieldCtx.fieldSchema;

	var obj = $.extend({}, {
		name: schema.Name,
		desc: schema.Description,
		type: schema.FieldType.toLowerCase(),
		aux: {
			required: schema.Required,
			maxLength: schema.MaxLength
		}}, data);

	this.id = ko.observable(obj.id);
	this.name = ko.observable(obj.name || '');
	this.desc = ko.observable(obj.desc || this.name());
	this.type = ko.observable(obj.type || 'text');
	this.help = ko.observable(obj.help || '');
	this.value = ko.observable(spFieldCtx.fieldValue);
	this.aux = obj.aux;

	if (!templates[this.type()]) {
		this.type('text');
	}

	this.hasHelp = ko.computed(function() {
		return this.help().length > 0;
	}, this);

	// TODO: sharepoint injection
	if (data.id) {
		sharepoint.registerValue(spFieldCtx, this);
	}

	console.log(ko.toJS(this));
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